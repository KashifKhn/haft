package detector

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

const (
	CacheDir      = ".haft"
	ProfileFile   = "profile.json"
	ChecksumFile  = "checksum"
	DefaultMaxAge = 24 * time.Hour
)

type ProfileCache struct {
	fs         afero.Fs
	projectDir string
	maxAge     time.Duration
}

func NewProfileCache(projectDir string) *ProfileCache {
	return &ProfileCache{
		fs:         afero.NewOsFs(),
		projectDir: projectDir,
		maxAge:     DefaultMaxAge,
	}
}

func NewProfileCacheWithFs(fs afero.Fs, projectDir string) *ProfileCache {
	return &ProfileCache{
		fs:         fs,
		projectDir: projectDir,
		maxAge:     DefaultMaxAge,
	}
}

func (c *ProfileCache) SetMaxAge(maxAge time.Duration) {
	c.maxAge = maxAge
}

func (c *ProfileCache) getCacheDir() string {
	return filepath.Join(c.projectDir, CacheDir)
}

func (c *ProfileCache) getProfilePath() string {
	return filepath.Join(c.getCacheDir(), ProfileFile)
}

func (c *ProfileCache) getChecksumPath() string {
	return filepath.Join(c.getCacheDir(), ChecksumFile)
}

func (c *ProfileCache) Save(profile *ProjectProfile) error {
	cacheDir := c.getCacheDir()
	if err := c.fs.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	profile.DetectedAt = time.Now()

	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	if err := afero.WriteFile(c.fs, c.getProfilePath(), data, 0644); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	checksum, err := c.computeSourceChecksum()
	if err != nil {
		return nil
	}

	if err := afero.WriteFile(c.fs, c.getChecksumPath(), []byte(checksum), 0644); err != nil {
		return nil
	}

	return nil
}

func (c *ProfileCache) Load() (*ProjectProfile, error) {
	profilePath := c.getProfilePath()

	exists, err := afero.Exists(c.fs, profilePath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	data, err := afero.ReadFile(c.fs, profilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile: %w", err)
	}

	var profile ProjectProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	return &profile, nil
}

func (c *ProfileCache) IsValid() bool {
	profile, err := c.Load()
	if err != nil || profile == nil {
		return false
	}

	if profile.IsStale(c.maxAge) {
		return false
	}

	if c.hasSourceChanged() {
		return false
	}

	return true
}

func (c *ProfileCache) hasSourceChanged() bool {
	checksumPath := c.getChecksumPath()

	exists, err := afero.Exists(c.fs, checksumPath)
	if err != nil || !exists {
		return true
	}

	savedChecksum, err := afero.ReadFile(c.fs, checksumPath)
	if err != nil {
		return true
	}

	currentChecksum, err := c.computeSourceChecksum()
	if err != nil {
		return true
	}

	return string(savedChecksum) != currentChecksum
}

func (c *ProfileCache) computeSourceChecksum() (string, error) {
	srcDir := filepath.Join(c.projectDir, "src", "main", "java")

	exists, err := afero.DirExists(c.fs, srcDir)
	if err != nil || !exists {
		return "", fmt.Errorf("source directory not found")
	}

	var files []string
	err = afero.Walk(c.fs, srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".java" {
			relPath, _ := filepath.Rel(c.projectDir, path)
			files = append(files, fmt.Sprintf("%s:%d:%d", relPath, info.Size(), info.ModTime().Unix()))
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	hash := md5.New()
	for _, f := range files {
		hash.Write([]byte(f))
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (c *ProfileCache) Clear() error {
	cacheDir := c.getCacheDir()

	exists, err := afero.DirExists(c.fs, cacheDir)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	return c.fs.RemoveAll(cacheDir)
}

func (c *ProfileCache) Exists() bool {
	exists, _ := afero.Exists(c.fs, c.getProfilePath())
	return exists
}
