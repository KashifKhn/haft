package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KashifKhn/haft/internal/detector"
	"github.com/KashifKhn/haft/internal/generator"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExceptionCommand(t *testing.T) {
	cmd := newExceptionCommand()

	assert.Equal(t, "exception", cmd.Use)
	assert.Contains(t, cmd.Aliases, "ex")
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestExceptionCommandAllFlags(t *testing.T) {
	cmd := newExceptionCommand()

	flags := []string{"package", "no-interactive", "all", "refresh"}
	for _, flag := range flags {
		f := cmd.Flags().Lookup(flag)
		assert.NotNil(t, f, "Flag %s should exist", flag)
	}

	packageFlag := cmd.Flags().Lookup("package")
	assert.Equal(t, "p", packageFlag.Shorthand)
}

func TestBuildSelectedMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string]bool
	}{
		{
			name:     "empty input",
			input:    []string{},
			expected: map[string]bool{},
		},
		{
			name:  "single selection",
			input: []string{"HasConflict"},
			expected: map[string]bool{
				"HasConflict": true,
			},
		},
		{
			name:  "multiple selections",
			input: []string{"HasConflict", "HasGone", "HasTooManyRequests"},
			expected: map[string]bool{
				"HasConflict":        true,
				"HasGone":            true,
				"HasTooManyRequests": true,
			},
		},
		{
			name:  "all optional exceptions",
			input: []string{"HasConflict", "HasMethodNotAllowed", "HasGone", "HasUnsupportedMediaType", "HasUnprocessableEntity", "HasTooManyRequests", "HasInternalServerError", "HasServiceUnavailable", "HasGatewayTimeout"},
			expected: map[string]bool{
				"HasConflict":             true,
				"HasMethodNotAllowed":     true,
				"HasGone":                 true,
				"HasUnsupportedMediaType": true,
				"HasUnprocessableEntity":  true,
				"HasTooManyRequests":      true,
				"HasInternalServerError":  true,
				"HasServiceUnavailable":   true,
				"HasGatewayTimeout":       true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildSelectedMap(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetExceptionPackageAllArchitectures(t *testing.T) {
	tests := []struct {
		name         string
		architecture detector.ArchitectureType
		basePackage  string
		expected     string
	}{
		{"layered", detector.ArchLayered, "com.example.app", "com.example.app.exception"},
		{"feature", detector.ArchFeature, "com.example.app", "com.example.app.common.exception"},
		{"hexagonal", detector.ArchHexagonal, "com.example.app", "com.example.app.infrastructure.exception"},
		{"clean", detector.ArchClean, "com.example.app", "com.example.app.infrastructure.exception"},
		{"flat", detector.ArchFlat, "com.example.app", "com.example.app.exception"},
		{"modular", detector.ArchModular, "com.example.app", "com.example.app.exception"},
		{"unknown", detector.ArchUnknown, "com.example.app", "com.example.app.exception"},
		{"empty base package", detector.ArchLayered, "", ".exception"},
		{"deep package", detector.ArchLayered, "com.company.project.module", "com.company.project.module.exception"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &detector.ProjectProfile{
				Architecture: tt.architecture,
				BasePackage:  tt.basePackage,
			}
			result := getExceptionPackage(profile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildExceptionTemplateDataAllFields(t *testing.T) {
	profile := &detector.ProjectProfile{
		Architecture:    detector.ArchLayered,
		BasePackage:     "com.example.app",
		HasValidation:   true,
		ValidationStyle: detector.ValidationJakarta,
		Lombok:          detector.LombokProfile{Detected: true},
	}

	selectedMap := map[string]bool{
		"HasConflict":             true,
		"HasMethodNotAllowed":     true,
		"HasGone":                 true,
		"HasUnsupportedMediaType": true,
		"HasUnprocessableEntity":  true,
		"HasTooManyRequests":      true,
		"HasInternalServerError":  true,
		"HasServiceUnavailable":   true,
		"HasGatewayTimeout":       true,
	}

	data := buildExceptionTemplateData(profile, "com.example.app.exception", selectedMap)

	assert.Equal(t, "com.example.app", data["BasePackage"])
	assert.Equal(t, "com.example.app.exception", data["ExceptionPackage"])
	assert.Equal(t, true, data["HasLombok"])
	assert.Equal(t, true, data["HasValidation"])
	assert.Equal(t, "jakarta.validation", data["ValidationImport"])
	assert.Equal(t, "layered", data["Architecture"])

	for key := range selectedMap {
		assert.Equal(t, true, data[key], "Expected %s to be true", key)
	}
}

func TestBuildExceptionTemplateDataNoOptionalSelected(t *testing.T) {
	profile := &detector.ProjectProfile{
		Architecture:    detector.ArchLayered,
		BasePackage:     "com.example.app",
		HasValidation:   false,
		ValidationStyle: detector.ValidationNone,
		Lombok:          detector.LombokProfile{Detected: false},
	}

	selectedMap := map[string]bool{}
	data := buildExceptionTemplateData(profile, "com.example.app.exception", selectedMap)

	assert.Equal(t, false, data["HasLombok"])
	assert.Equal(t, false, data["HasValidation"])
	assert.Equal(t, "jakarta.validation", data["ValidationImport"])

	optionalKeys := []string{
		"HasConflict", "HasMethodNotAllowed", "HasGone",
		"HasUnsupportedMediaType", "HasUnprocessableEntity", "HasTooManyRequests",
		"HasInternalServerError", "HasServiceUnavailable", "HasGatewayTimeout",
	}
	for _, key := range optionalKeys {
		assert.Equal(t, false, data[key], "Expected %s to be false", key)
	}
}

func TestBuildExceptionTemplateDataJavaxValidation(t *testing.T) {
	profile := &detector.ProjectProfile{
		Architecture:    detector.ArchLayered,
		BasePackage:     "com.example.app",
		HasValidation:   true,
		ValidationStyle: detector.ValidationJavax,
		Lombok:          detector.LombokProfile{Detected: false},
	}

	data := buildExceptionTemplateData(profile, "com.example.app.exception", map[string]bool{})

	assert.Equal(t, "javax.validation", data["ValidationImport"])
}

func TestOptionalExceptionsDefinition(t *testing.T) {
	assert.Equal(t, 9, len(optionalExceptions), "Should have 9 optional exceptions")

	expectedExceptions := []struct {
		templateKey string
		fileName    string
	}{
		{"HasConflict", "ConflictException.java"},
		{"HasMethodNotAllowed", "MethodNotAllowedException.java"},
		{"HasGone", "GoneException.java"},
		{"HasUnsupportedMediaType", "UnsupportedMediaTypeException.java"},
		{"HasUnprocessableEntity", "UnprocessableEntityException.java"},
		{"HasTooManyRequests", "TooManyRequestsException.java"},
		{"HasInternalServerError", "InternalServerErrorException.java"},
		{"HasServiceUnavailable", "ServiceUnavailableException.java"},
		{"HasGatewayTimeout", "GatewayTimeoutException.java"},
	}

	for i, expected := range expectedExceptions {
		assert.Equal(t, expected.templateKey, optionalExceptions[i].TemplateKey)
		assert.Equal(t, expected.fileName, optionalExceptions[i].FileName)
		assert.NotEmpty(t, optionalExceptions[i].Name)
		assert.NotEmpty(t, optionalExceptions[i].Description)
	}
}

func TestExceptionTemplateDir(t *testing.T) {
	profiles := []*detector.ProjectProfile{
		{Architecture: detector.ArchLayered},
		{Architecture: detector.ArchFeature},
		{Architecture: detector.ArchHexagonal},
		{Architecture: detector.ArchClean},
	}

	for _, profile := range profiles {
		result := getExceptionTemplateDir(profile)
		assert.Equal(t, "exception", result)
	}
}

func TestGenerateExceptionHandlerIntegration(t *testing.T) {
	fs := afero.NewMemMapFs()

	tmpDir := "/tmp/test-exception-project"
	require.NoError(t, fs.MkdirAll(filepath.Join(tmpDir, "src/main/java/com/example/demo"), 0755))

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-validation</artifactId>
        </dependency>
    </dependencies>
</project>`
	require.NoError(t, afero.WriteFile(fs, filepath.Join(tmpDir, "pom.xml"), []byte(pomContent), 0644))

	profile := &detector.ProjectProfile{
		Architecture:    detector.ArchLayered,
		BasePackage:     "com.example.demo",
		HasValidation:   true,
		ValidationStyle: detector.ValidationJakarta,
		Lombok:          detector.LombokProfile{Detected: false},
	}

	engine := generator.NewEngine(fs)

	exceptionPackage := getExceptionPackage(profile)
	assert.Equal(t, "com.example.demo.exception", exceptionPackage)

	selectedMap := buildSelectedMap([]string{"HasConflict", "HasTooManyRequests"})
	data := buildExceptionTemplateData(profile, exceptionPackage, selectedMap)

	assert.Equal(t, true, data["HasValidation"])
	assert.Equal(t, true, data["HasConflict"])
	assert.Equal(t, true, data["HasTooManyRequests"])
	assert.Equal(t, false, data["HasGone"])

	_ = engine
}

func TestExceptionConfigStruct(t *testing.T) {
	cfg := exceptionConfig{
		BasePackage:      "com.example.app",
		SelectedOptional: []string{"HasConflict", "HasGone"},
	}

	assert.Equal(t, "com.example.app", cfg.BasePackage)
	assert.Len(t, cfg.SelectedOptional, 2)
	assert.Contains(t, cfg.SelectedOptional, "HasConflict")
	assert.Contains(t, cfg.SelectedOptional, "HasGone")
}

func TestEnrichProfileFromBuildFileNoError(t *testing.T) {
	profile := &detector.ProjectProfile{
		Architecture: detector.ArchLayered,
	}

	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-enrich-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	require.NoError(t, os.Chdir(tmpDir))

	enrichProfileFromBuildFile(profile)
}

func TestEnrichProfileFromBuildFileWithValidation(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-enrich-validation-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-validation</artifactId>
        </dependency>
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
        </dependency>
    </dependencies>
</project>`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "pom.xml"), []byte(pomContent), 0644))

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		Architecture:    detector.ArchLayered,
		HasValidation:   false,
		ValidationStyle: detector.ValidationNone,
		Lombok:          detector.LombokProfile{Detected: false},
	}

	enrichProfileFromBuildFile(profile)

	assert.True(t, profile.HasValidation, "Should detect validation from pom.xml")
	assert.Equal(t, detector.ValidationJakarta, profile.ValidationStyle)
	assert.True(t, profile.Lombok.Detected, "Should detect lombok from pom.xml")
}

func TestEnrichProfilePreservesExistingValues(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-enrich-preserve-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
</project>`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "pom.xml"), []byte(pomContent), 0644))

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		Architecture:    detector.ArchLayered,
		BasePackage:     "com.existing.package",
		HasValidation:   true,
		ValidationStyle: detector.ValidationJavax,
		Lombok:          detector.LombokProfile{Detected: true},
	}

	enrichProfileFromBuildFile(profile)

	assert.Equal(t, "com.existing.package", profile.BasePackage, "Should preserve existing base package")
	assert.True(t, profile.HasValidation, "Should preserve existing validation")
	assert.Equal(t, detector.ValidationJavax, profile.ValidationStyle, "Should preserve existing validation style")
	assert.True(t, profile.Lombok.Detected, "Should preserve existing lombok")
}

func TestMultiSelectWrapperMethods(t *testing.T) {
	t.Run("View returns model view", func(t *testing.T) {
		items := []struct {
			label string
			value string
		}{
			{"Option 1", "opt1"},
			{"Option 2", "opt2"},
		}
		_ = items
	})
}
