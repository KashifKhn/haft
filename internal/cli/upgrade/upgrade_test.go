package upgrade

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOS(t *testing.T) {
	os := GetOS()
	validOS := []string{"linux", "darwin", "windows"}

	found := false
	for _, v := range validOS {
		if os == v {
			found = true
			break
		}
	}

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		assert.True(t, found, "GetOS() should return a valid OS")
	}
}

func TestGetArch(t *testing.T) {
	arch := GetArch()
	validArch := []string{"amd64", "arm64", "unknown"}

	found := false
	for _, v := range validArch {
		if arch == v {
			found = true
			break
		}
	}

	assert.True(t, found, "GetArch() should return a valid architecture")
}

func TestGetPlatformInfo(t *testing.T) {
	info, err := GetPlatformInfo()

	if runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64" {
		require.NoError(t, err)
		assert.NotEmpty(t, info.OS)
		assert.NotEmpty(t, info.Arch)
		assert.NotEmpty(t, info.BinaryName)
		assert.NotEmpty(t, info.ArchiveExt)
	}
}

func TestPlatformInfoGetArchiveName(t *testing.T) {
	tests := []struct {
		name     string
		platform PlatformInfo
		version  string
		expected string
	}{
		{
			name:     "linux amd64",
			platform: PlatformInfo{OS: "linux", Arch: "amd64", ArchiveExt: ".tar.gz"},
			version:  "v1.0.0",
			expected: "haft-linux-amd64.tar.gz",
		},
		{
			name:     "darwin arm64",
			platform: PlatformInfo{OS: "darwin", Arch: "arm64", ArchiveExt: ".tar.gz"},
			version:  "v1.0.0",
			expected: "haft-darwin-arm64.tar.gz",
		},
		{
			name:     "windows amd64",
			platform: PlatformInfo{OS: "windows", Arch: "amd64", ArchiveExt: ".zip"},
			version:  "v1.0.0",
			expected: "haft-windows-amd64.zip",
		},
		{
			name:     "darwin amd64",
			platform: PlatformInfo{OS: "darwin", Arch: "amd64", ArchiveExt: ".tar.gz"},
			version:  "v2.5.0",
			expected: "haft-darwin-amd64.tar.gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.platform.GetArchiveName(tt.version)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlatformInfoGetBinaryNameInArchive(t *testing.T) {
	tests := []struct {
		name     string
		platform PlatformInfo
		expected string
	}{
		{
			name:     "linux amd64",
			platform: PlatformInfo{OS: "linux", Arch: "amd64"},
			expected: "haft-linux-amd64",
		},
		{
			name:     "windows amd64",
			platform: PlatformInfo{OS: "windows", Arch: "amd64"},
			expected: "haft-windows-amd64.exe",
		},
		{
			name:     "darwin arm64",
			platform: PlatformInfo{OS: "darwin", Arch: "arm64"},
			expected: "haft-darwin-arm64",
		},
		{
			name:     "linux arm64",
			platform: PlatformInfo{OS: "linux", Arch: "arm64"},
			expected: "haft-linux-arm64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.platform.GetBinaryNameInArchive()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlatformInfoGetDownloadURL(t *testing.T) {
	platform := &PlatformInfo{OS: "linux", Arch: "amd64", ArchiveExt: ".tar.gz"}
	url := platform.GetDownloadURL("v1.0.0")

	assert.Contains(t, url, "github.com")
	assert.Contains(t, url, "KashifKhn/haft")
	assert.Contains(t, url, "v1.0.0")
	assert.Contains(t, url, "haft-linux-amd64.tar.gz")
}

func TestPlatformInfoGetChecksumsURL(t *testing.T) {
	platform := &PlatformInfo{OS: "linux", Arch: "amd64", ArchiveExt: ".tar.gz"}
	url := platform.GetChecksumsURL("v1.0.0")

	assert.Contains(t, url, "github.com")
	assert.Contains(t, url, "KashifKhn/haft")
	assert.Contains(t, url, "v1.0.0")
	assert.Contains(t, url, "checksums.txt")
}

func TestPlatformInfoString(t *testing.T) {
	tests := []struct {
		platform PlatformInfo
		expected string
	}{
		{PlatformInfo{OS: "linux", Arch: "amd64"}, "linux/amd64"},
		{PlatformInfo{OS: "darwin", Arch: "arm64"}, "darwin/arm64"},
		{PlatformInfo{OS: "windows", Arch: "amd64"}, "windows/amd64"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.platform.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMajor int
		wantMinor int
		wantPatch int
		wantPre   string
		wantErr   bool
	}{
		{
			name:      "simple version",
			input:     "1.2.3",
			wantMajor: 1,
			wantMinor: 2,
			wantPatch: 3,
		},
		{
			name:      "version with v prefix",
			input:     "v1.2.3",
			wantMajor: 1,
			wantMinor: 2,
			wantPatch: 3,
		},
		{
			name:      "version with prerelease",
			input:     "v1.2.3-beta.1",
			wantMajor: 1,
			wantMinor: 2,
			wantPatch: 3,
			wantPre:   "beta.1",
		},
		{
			name:      "version with alpha prerelease",
			input:     "v2.0.0-alpha",
			wantMajor: 2,
			wantMinor: 0,
			wantPatch: 0,
			wantPre:   "alpha",
		},
		{
			name:      "version with rc prerelease",
			input:     "v1.0.0-rc.1",
			wantMajor: 1,
			wantMinor: 0,
			wantPatch: 0,
			wantPre:   "rc.1",
		},
		{
			name:      "major only",
			input:     "1",
			wantMajor: 1,
		},
		{
			name:      "major.minor only",
			input:     "1.2",
			wantMajor: 1,
			wantMinor: 2,
		},
		{
			name:      "large version numbers",
			input:     "v100.200.300",
			wantMajor: 100,
			wantMinor: 200,
			wantPatch: 300,
		},
		{
			name:      "version with plus metadata",
			input:     "v1.0.0+build.123",
			wantMajor: 1,
			wantMinor: 0,
			wantPatch: 0,
			wantPre:   "build.123",
		},
		{
			name:    "invalid major version",
			input:   "vx.0.0",
			wantErr: true,
		},
		{
			name:    "invalid minor version",
			input:   "v1.x.0",
			wantErr: true,
		},
		{
			name:    "invalid patch version",
			input:   "v1.0.x",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ParseVersion(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantMajor, v.Major)
			assert.Equal(t, tt.wantMinor, v.Minor)
			assert.Equal(t, tt.wantPatch, v.Patch)
			assert.Equal(t, tt.wantPre, v.Prerelease)
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		latest   string
		expected int
	}{
		{
			name:     "same version",
			current:  "v1.0.0",
			latest:   "v1.0.0",
			expected: 0,
		},
		{
			name:     "current older major",
			current:  "v1.0.0",
			latest:   "v2.0.0",
			expected: -1,
		},
		{
			name:     "current newer major",
			current:  "v2.0.0",
			latest:   "v1.0.0",
			expected: 1,
		},
		{
			name:     "current older minor",
			current:  "v1.0.0",
			latest:   "v1.1.0",
			expected: -1,
		},
		{
			name:     "current newer minor",
			current:  "v1.2.0",
			latest:   "v1.1.0",
			expected: 1,
		},
		{
			name:     "current older patch",
			current:  "v1.0.0",
			latest:   "v1.0.1",
			expected: -1,
		},
		{
			name:     "current newer patch",
			current:  "v1.0.2",
			latest:   "v1.0.1",
			expected: 1,
		},
		{
			name:     "dev version",
			current:  "dev",
			latest:   "v1.0.0",
			expected: -1,
		},
		{
			name:     "empty version",
			current:  "",
			latest:   "v1.0.0",
			expected: -1,
		},
		{
			name:     "prerelease vs stable",
			current:  "v1.0.0-beta.1",
			latest:   "v1.0.0",
			expected: -1,
		},
		{
			name:     "stable vs prerelease",
			current:  "v1.0.0",
			latest:   "v1.0.0-beta.1",
			expected: 1,
		},
		{
			name:     "both prerelease same version",
			current:  "v1.0.0-alpha",
			latest:   "v1.0.0-beta",
			expected: 0,
		},
		{
			name:     "version without v prefix",
			current:  "1.0.0",
			latest:   "1.1.0",
			expected: -1,
		},
		{
			name:     "dirty version",
			current:  "v1.0.0-dirty",
			latest:   "v1.0.0",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareVersions(tt.current, tt.latest)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCompareVersionsError(t *testing.T) {
	_, err := CompareVersions("v1.0.0", "invalid.version.x")
	assert.Error(t, err)
}

func TestIsNewerAvailable(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		latest   string
		expected bool
	}{
		{
			name:     "newer available",
			current:  "v1.0.0",
			latest:   "v1.1.0",
			expected: true,
		},
		{
			name:     "same version",
			current:  "v1.0.0",
			latest:   "v1.0.0",
			expected: false,
		},
		{
			name:     "current newer",
			current:  "v1.1.0",
			latest:   "v1.0.0",
			expected: false,
		},
		{
			name:     "dev version always gets update",
			current:  "dev",
			latest:   "v0.0.1",
			expected: true,
		},
		{
			name:     "prerelease to stable",
			current:  "v1.0.0-beta",
			latest:   "v1.0.0",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := IsNewerAvailable(tt.current, tt.latest)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.0.0", "v1.0.0"},
		{"v1.0.0", "v1.0.0"},
		{"dev", "dev"},
		{"", ""},
		{"0.1.14", "v0.1.14"},
		{"v0.1.14-dirty", "v0.1.14-dirty"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizeVersion(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateChecksum(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	err := os.WriteFile(testFile, []byte("hello world"), 0644)
	require.NoError(t, err)

	checksum, err := calculateChecksum(testFile)
	require.NoError(t, err)
	assert.NotEmpty(t, checksum)
	assert.Len(t, checksum, 64)
	assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", checksum)
}

func TestCalculateChecksumFileNotFound(t *testing.T) {
	_, err := calculateChecksum("/nonexistent/file/path")
	assert.Error(t, err)
}

func TestVerifyChecksum(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	err := os.WriteFile(testFile, []byte("hello world"), 0644)
	require.NoError(t, err)

	checksum, err := calculateChecksum(testFile)
	require.NoError(t, err)

	err = VerifyChecksum(testFile, checksum)
	assert.NoError(t, err)

	err = VerifyChecksum(testFile, "invalid-checksum")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checksum mismatch")

	err = VerifyChecksum(testFile, strings.ToUpper(checksum))
	assert.NoError(t, err, "checksum should be case-insensitive")
}

func TestVerifyChecksumFileNotFound(t *testing.T) {
	err := VerifyChecksum("/nonexistent/file", "somechecksum")
	assert.Error(t, err)
}

func TestCreateAndRestoreBackup(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-binary")
	content := []byte("binary content")

	err := os.WriteFile(testFile, content, 0755)
	require.NoError(t, err)

	backup, err := CreateBackup(testFile)
	require.NoError(t, err)
	assert.NotEmpty(t, backup.BackupPath)
	assert.Equal(t, testFile, backup.OriginalPath)
	assert.False(t, backup.CreatedAt.IsZero())

	backupContent, err := os.ReadFile(backup.BackupPath)
	require.NoError(t, err)
	assert.Equal(t, content, backupContent)

	newContent := []byte("modified content")
	err = os.WriteFile(testFile, newContent, 0755)
	require.NoError(t, err)

	err = RestoreBackup(backup)
	require.NoError(t, err)

	restoredContent, err := os.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, content, restoredContent)

	CleanupBackup(backup)
	_, err = os.Stat(backup.BackupPath)
	assert.True(t, os.IsNotExist(err))
}

func TestCreateBackupNonExistentFile(t *testing.T) {
	_, err := CreateBackup("/nonexistent/file/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "binary not found")
}

func TestRestoreBackupNil(t *testing.T) {
	err := RestoreBackup(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no backup information provided")
}

func TestRestoreBackupMissingBackupFile(t *testing.T) {
	backup := &BackupInfo{
		OriginalPath: "/some/path",
		BackupPath:   "/nonexistent/backup/path",
	}
	err := RestoreBackup(backup)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "backup file not found")
}

func TestCleanupBackupNil(t *testing.T) {
	CleanupBackup(nil)
}

func TestCleanupBackupEmptyPath(t *testing.T) {
	CleanupBackup(&BackupInfo{BackupPath: ""})
}

func TestInstallBinary(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "source")
	destFile := filepath.Join(tmpDir, "dest", "binary")
	content := []byte("binary content")

	err := os.WriteFile(sourceFile, content, 0755)
	require.NoError(t, err)

	err = InstallBinary(sourceFile, destFile)
	require.NoError(t, err)

	installedContent, err := os.ReadFile(destFile)
	require.NoError(t, err)
	assert.Equal(t, content, installedContent)

	info, err := os.Stat(destFile)
	require.NoError(t, err)
	if runtime.GOOS != "windows" {
		assert.True(t, info.Mode()&0100 != 0, "binary should be executable")
	}
}

func TestInstallBinaryOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "source")
	destFile := filepath.Join(tmpDir, "dest")

	oldContent := []byte("old content")
	newContent := []byte("new content")

	err := os.WriteFile(destFile, oldContent, 0755)
	require.NoError(t, err)

	err = os.WriteFile(sourceFile, newContent, 0755)
	require.NoError(t, err)

	err = InstallBinary(sourceFile, destFile)
	require.NoError(t, err)

	installedContent, err := os.ReadFile(destFile)
	require.NoError(t, err)
	assert.Equal(t, newContent, installedContent)
}

func TestInstallBinarySourceNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	err := InstallBinary("/nonexistent/source", filepath.Join(tmpDir, "dest"))
	assert.Error(t, err)
}

func TestInstallToMultipleLocations(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "source")
	content := []byte("binary content")

	err := os.WriteFile(sourceFile, content, 0755)
	require.NoError(t, err)

	dest1 := filepath.Join(tmpDir, "dest1", "haft")
	dest2 := filepath.Join(tmpDir, "dest2", "haft")
	dest3 := "/nonexistent/readonly/path/haft"

	installed, errors := InstallToMultipleLocations(sourceFile, []string{dest1, dest2, dest3})

	assert.Len(t, installed, 2)
	assert.Contains(t, installed, dest1)
	assert.Contains(t, installed, dest2)
	assert.Len(t, errors, 1)
}

func TestSetAndGetCurrentVersion(t *testing.T) {
	SetCurrentVersion("v1.0.0")
	assert.Equal(t, "v1.0.0", GetCurrentVersion())

	SetCurrentVersion("v2.0.0")
	assert.Equal(t, "v2.0.0", GetCurrentVersion())

	SetCurrentVersion("dev")
	assert.Equal(t, "dev", GetCurrentVersion())
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		version  Version
		expected string
	}{
		{
			version:  Version{Major: 1, Minor: 2, Patch: 3},
			expected: "v1.2.3",
		},
		{
			version:  Version{Major: 1, Minor: 2, Patch: 3, Prerelease: "beta.1"},
			expected: "v1.2.3-beta.1",
		},
		{
			version:  Version{Major: 0, Minor: 0, Patch: 1},
			expected: "v0.0.1",
		},
		{
			version:  Version{Major: 10, Minor: 20, Patch: 30, Prerelease: "rc.1"},
			expected: "v10.20.30-rc.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.version.String())
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
		{1024 * 1024 * 1024 * 1024, "1.0 TB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsPath(t *testing.T) {
	paths := []string{"/usr/local/bin", "/home/user/.local/bin", "/opt/bin"}

	assert.True(t, containsPath(paths, "/usr/local/bin"))
	assert.True(t, containsPath(paths, "/home/user/.local/bin"))
	assert.True(t, containsPath(paths, "/opt/bin"))
	assert.False(t, containsPath(paths, "/usr/bin"))
	assert.False(t, containsPath(paths, ""))
	assert.False(t, containsPath([]string{}, "/usr/local/bin"))
}

func TestGetExecutablePath(t *testing.T) {
	path, err := GetExecutablePath()
	require.NoError(t, err)
	assert.NotEmpty(t, path)

	_, err = os.Stat(path)
	assert.NoError(t, err)
}

func TestGetInstallDir(t *testing.T) {
	dir := GetInstallDir()
	assert.NotEmpty(t, dir)

	if runtime.GOOS == "windows" {
		assert.Contains(t, dir, "AppData")
	}
}

func TestFindAllInstallations(t *testing.T) {
	installations := FindAllInstallations()
	assert.IsType(t, []string{}, installations)
}

func TestGetSearchPaths(t *testing.T) {
	paths := getSearchPaths()
	assert.NotEmpty(t, paths)

	if runtime.GOOS != "windows" {
		assert.Contains(t, paths, "/usr/local/bin")
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "source.txt")
	dstFile := filepath.Join(tmpDir, "dest.txt")
	content := []byte("test content for copy")

	err := os.WriteFile(srcFile, content, 0644)
	require.NoError(t, err)

	err = copyFile(srcFile, dstFile)
	require.NoError(t, err)

	copiedContent, err := os.ReadFile(dstFile)
	require.NoError(t, err)
	assert.Equal(t, content, copiedContent)
}

func TestCopyFileSourceNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	err := copyFile("/nonexistent/source", filepath.Join(tmpDir, "dest"))
	assert.Error(t, err)
}

func TestCopyFileDestNotWritable(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "source.txt")
	err := os.WriteFile(srcFile, []byte("content"), 0644)
	require.NoError(t, err)

	err = copyFile(srcFile, "/nonexistent/dir/dest")
	assert.Error(t, err)
}

func TestExtractFromTarGz(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "test.tar.gz")
	binaryName := "haft-linux-amd64"
	binaryContent := []byte("#!/bin/bash\necho 'haft'")

	f, err := os.Create(archivePath)
	require.NoError(t, err)

	gzw := gzip.NewWriter(f)
	tw := tar.NewWriter(gzw)

	hdr := &tar.Header{
		Name: binaryName,
		Mode: 0755,
		Size: int64(len(binaryContent)),
	}
	err = tw.WriteHeader(hdr)
	require.NoError(t, err)
	_, err = tw.Write(binaryContent)
	require.NoError(t, err)

	tw.Close()
	gzw.Close()
	f.Close()

	extractedPath, err := extractFromTarGz(archivePath, tmpDir, binaryName)
	require.NoError(t, err)
	assert.Contains(t, extractedPath, binaryName)

	extractedContent, err := os.ReadFile(extractedPath)
	require.NoError(t, err)
	assert.Equal(t, binaryContent, extractedContent)
}

func TestExtractFromTarGzBinaryNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "test.tar.gz")

	f, err := os.Create(archivePath)
	require.NoError(t, err)

	gzw := gzip.NewWriter(f)
	tw := tar.NewWriter(gzw)

	hdr := &tar.Header{
		Name: "other-file",
		Mode: 0644,
		Size: 5,
	}
	_ = tw.WriteHeader(hdr)
	_, _ = tw.Write([]byte("hello"))

	tw.Close()
	gzw.Close()
	f.Close()

	_, err = extractFromTarGz(archivePath, tmpDir, "haft-linux-amd64")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in archive")
}

func TestExtractFromTarGzInvalidArchive(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "invalid.tar.gz")

	err := os.WriteFile(archivePath, []byte("not a valid archive"), 0644)
	require.NoError(t, err)

	_, err = extractFromTarGz(archivePath, tmpDir, "haft")
	assert.Error(t, err)
}

func TestExtractFromZip(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "test.zip")
	binaryName := "haft-windows-amd64.exe"
	binaryContent := []byte("MZ...fake exe content")

	f, err := os.Create(archivePath)
	require.NoError(t, err)

	zw := zip.NewWriter(f)

	fw, err := zw.Create(binaryName)
	require.NoError(t, err)
	_, err = fw.Write(binaryContent)
	require.NoError(t, err)

	zw.Close()
	f.Close()

	extractedPath, err := extractFromZip(archivePath, tmpDir, binaryName)
	require.NoError(t, err)
	assert.Contains(t, extractedPath, binaryName)

	extractedContent, err := os.ReadFile(extractedPath)
	require.NoError(t, err)
	assert.Equal(t, binaryContent, extractedContent)
}

func TestExtractFromZipBinaryNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "test.zip")

	f, err := os.Create(archivePath)
	require.NoError(t, err)

	zw := zip.NewWriter(f)
	fw, _ := zw.Create("other-file.txt")
	_, _ = fw.Write([]byte("hello"))
	zw.Close()
	f.Close()

	_, err = extractFromZip(archivePath, tmpDir, "haft.exe")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in archive")
}

func TestExtractFromZipInvalidArchive(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "invalid.zip")

	err := os.WriteFile(archivePath, []byte("not a valid zip"), 0644)
	require.NoError(t, err)

	_, err = extractFromZip(archivePath, tmpDir, "haft")
	assert.Error(t, err)
}

func TestExtractBinary(t *testing.T) {
	tmpDir := t.TempDir()
	binaryContent := []byte("binary content")

	t.Run("tar.gz", func(t *testing.T) {
		archivePath := filepath.Join(tmpDir, "test.tar.gz")
		binaryName := "haft-linux-amd64"

		f, err := os.Create(archivePath)
		require.NoError(t, err)

		gzw := gzip.NewWriter(f)
		tw := tar.NewWriter(gzw)

		hdr := &tar.Header{
			Name: binaryName,
			Mode: 0755,
			Size: int64(len(binaryContent)),
		}
		_ = tw.WriteHeader(hdr)
		_, _ = tw.Write(binaryContent)
		tw.Close()
		gzw.Close()
		f.Close()

		platform := &PlatformInfo{OS: "linux", Arch: "amd64", ArchiveExt: ".tar.gz"}
		extracted, err := ExtractBinary(archivePath, platform)
		require.NoError(t, err)
		assert.NotEmpty(t, extracted)
	})

	t.Run("zip", func(t *testing.T) {
		archivePath := filepath.Join(tmpDir, "test.zip")
		binaryName := "haft-windows-amd64.exe"

		f, err := os.Create(archivePath)
		require.NoError(t, err)

		zw := zip.NewWriter(f)
		fw, _ := zw.Create(binaryName)
		_, _ = fw.Write(binaryContent)
		zw.Close()
		f.Close()

		platform := &PlatformInfo{OS: "windows", Arch: "amd64", ArchiveExt: ".zip"}
		extracted, err := ExtractBinary(archivePath, platform)
		require.NoError(t, err)
		assert.NotEmpty(t, extracted)
	})
}

func TestCleanupDownload(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "haft-test-*")
	require.NoError(t, err)

	testFile := filepath.Join(tmpDir, "test.tar.gz")
	err = os.WriteFile(testFile, []byte("content"), 0644)
	require.NoError(t, err)

	result := &DownloadResult{FilePath: testFile}
	CleanupDownload(result)

	_, err = os.Stat(tmpDir)
	assert.True(t, os.IsNotExist(err))
}

func TestCleanupDownloadNil(t *testing.T) {
	CleanupDownload(nil)
}

func TestCleanupDownloadEmptyPath(t *testing.T) {
	CleanupDownload(&DownloadResult{FilePath: ""})
}

func TestVerifyInstallation(t *testing.T) {
	t.Run("nonexistent binary", func(t *testing.T) {
		err := VerifyInstallation("/nonexistent/binary")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "binary not found")
	})

	t.Run("existing binary without version command", func(t *testing.T) {
		tmpDir := t.TempDir()
		fakeBinary := filepath.Join(tmpDir, "fake")
		err := os.WriteFile(fakeBinary, []byte("not executable"), 0644)
		require.NoError(t, err)

		err = VerifyInstallation(fakeBinary)
		assert.Error(t, err)
	})
}

func TestUpgradeResultStruct(t *testing.T) {
	result := &UpgradeResult{
		CurrentVersion:  "v1.0.0",
		LatestVersion:   "v1.1.0",
		UpdateAvailable: true,
		Upgraded:        true,
		InstalledPaths:  []string{"/usr/local/bin/haft"},
	}
	result.Platform.OS = "linux"
	result.Platform.Arch = "amd64"

	assert.Equal(t, "v1.0.0", result.CurrentVersion)
	assert.Equal(t, "v1.1.0", result.LatestVersion)
	assert.True(t, result.UpdateAvailable)
	assert.True(t, result.Upgraded)
	assert.Equal(t, "linux", result.Platform.OS)
	assert.Equal(t, "amd64", result.Platform.Arch)
}

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	assert.Equal(t, "upgrade", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)

	checkFlag := cmd.Flags().Lookup("check")
	assert.NotNil(t, checkFlag)
	assert.Equal(t, "c", checkFlag.Shorthand)

	forceFlag := cmd.Flags().Lookup("force")
	assert.NotNil(t, forceFlag)
	assert.Equal(t, "f", forceFlag.Shorthand)

	versionFlag := cmd.Flags().Lookup("version")
	assert.NotNil(t, versionFlag)
	assert.Equal(t, "v", versionFlag.Shorthand)

	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag)
}

func TestGitHubReleaseStruct(t *testing.T) {
	release := GitHubRelease{
		TagName:     "v1.0.0",
		Name:        "Release 1.0.0",
		Prerelease:  false,
		Draft:       false,
		PublishedAt: "2024-01-01T00:00:00Z",
		HTMLURL:     "https://github.com/owner/repo/releases/tag/v1.0.0",
	}

	assert.Equal(t, "v1.0.0", release.TagName)
	assert.Equal(t, "Release 1.0.0", release.Name)
	assert.False(t, release.Prerelease)
	assert.False(t, release.Draft)
}

func TestDownloadResultStruct(t *testing.T) {
	result := &DownloadResult{
		FilePath:   "/tmp/haft.tar.gz",
		Size:       1024000,
		Checksum:   "abc123",
		BinaryPath: "/tmp/haft",
	}

	assert.Equal(t, "/tmp/haft.tar.gz", result.FilePath)
	assert.Equal(t, int64(1024000), result.Size)
	assert.Equal(t, "abc123", result.Checksum)
}

func TestBackupInfoStruct(t *testing.T) {
	backup := &BackupInfo{
		OriginalPath: "/usr/local/bin/haft",
		BackupPath:   "/tmp/haft.backup",
	}

	assert.Equal(t, "/usr/local/bin/haft", backup.OriginalPath)
	assert.Equal(t, "/tmp/haft.backup", backup.BackupPath)
}

func TestIsWritable(t *testing.T) {
	tmpDir := t.TempDir()

	assert.True(t, isWritable(tmpDir))
	assert.False(t, isWritable("/nonexistent/path"))

	testFile := filepath.Join(tmpDir, "testfile")
	err := os.WriteFile(testFile, []byte("content"), 0644)
	require.NoError(t, err)
	assert.False(t, isWritable(testFile))
}

func TestDownloadFileWithMockServer(t *testing.T) {
	expectedContent := []byte("downloaded content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(expectedContent)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	destPath := filepath.Join(tmpDir, "downloaded")

	err := downloadFile(server.URL, destPath)
	require.NoError(t, err)

	content, err := os.ReadFile(destPath)
	require.NoError(t, err)
	assert.Equal(t, expectedContent, content)
}

func TestDownloadFileNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	destPath := filepath.Join(tmpDir, "downloaded")

	err := downloadFile(server.URL, destPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDownloadFileServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	destPath := filepath.Join(tmpDir, "downloaded")

	err := downloadFile(server.URL, destPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed")
}

func TestFetchChecksumsWithMockServer(t *testing.T) {
	checksumContent := `abc123def456  haft-linux-amd64.tar.gz
789xyz000111  haft-darwin-arm64.tar.gz
fedcba654321  haft-windows-amd64.zip`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, checksumContent)
	}))
	defer server.Close()

	origURL := fmt.Sprintf("https://github.com/%s/%s/releases/download", RepoOwner, RepoName)
	t.Logf("Original URL pattern: %s", origURL)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFetchChecksumsNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestOutputJSON(t *testing.T) {
	result := &UpgradeResult{
		CurrentVersion:  "v1.0.0",
		LatestVersion:   "v1.1.0",
		UpdateAvailable: true,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := outputJSON(result)
	assert.NoError(t, err)

	w.Close()
	os.Stdout = oldStdout

	var output strings.Builder
	_, _ = io.Copy(&output, r)

	assert.Contains(t, output.String(), "v1.0.0")
	assert.Contains(t, output.String(), "v1.1.0")
	assert.Contains(t, output.String(), "update_available")
}

func TestOutputError(t *testing.T) {
	result := &UpgradeResult{}
	testErr := fmt.Errorf("test error")

	returnedErr := outputError(false, result, testErr)
	assert.Equal(t, testErr, returnedErr)
	assert.Empty(t, result.Error)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	_ = outputError(true, result, testErr)

	w.Close()
	os.Stdout = oldStdout

	var output strings.Builder
	_, _ = io.Copy(&output, r)

	assert.Equal(t, "test error", result.Error)
	assert.Contains(t, output.String(), "test error")
}
