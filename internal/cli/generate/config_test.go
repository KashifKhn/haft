package generate

import (
	"testing"

	"github.com/KashifKhn/haft/internal/detector"
	"github.com/stretchr/testify/assert"
)

func TestConfigCommand(t *testing.T) {
	cmd := newConfigCommand()

	assert.Equal(t, "config", cmd.Use)
	assert.Contains(t, cmd.Aliases, "cfg")
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestConfigCommandFlags(t *testing.T) {
	cmd := newConfigCommand()

	pkgFlag := cmd.Flags().Lookup("package")
	assert.NotNil(t, pkgFlag)
	assert.Equal(t, "p", pkgFlag.Shorthand)

	noInteractiveFlag := cmd.Flags().Lookup("no-interactive")
	assert.NotNil(t, noInteractiveFlag)

	allFlag := cmd.Flags().Lookup("all")
	assert.NotNil(t, allFlag)

	refreshFlag := cmd.Flags().Lookup("refresh")
	assert.NotNil(t, refreshFlag)
}

func TestConfigOptionsDefinition(t *testing.T) {
	assert.Len(t, configOptions, 7)

	expectedConfigs := []string{
		"cors", "openapi", "jackson", "async", "cache", "auditing", "webmvc",
	}

	for i, expected := range expectedConfigs {
		assert.Equal(t, expected, configOptions[i].Key)
		assert.NotEmpty(t, configOptions[i].Name)
		assert.NotEmpty(t, configOptions[i].FileName)
		assert.NotEmpty(t, configOptions[i].Description)
	}
}

func TestBuildConfigSelectedMap(t *testing.T) {
	tests := []struct {
		name     string
		selected []string
		expected map[string]bool
	}{
		{
			name:     "empty input",
			selected: []string{},
			expected: map[string]bool{},
		},
		{
			name:     "single selection",
			selected: []string{"cors"},
			expected: map[string]bool{"cors": true},
		},
		{
			name:     "multiple selections",
			selected: []string{"cors", "jackson", "async"},
			expected: map[string]bool{
				"cors":    true,
				"jackson": true,
				"async":   true,
			},
		},
		{
			name:     "all config options",
			selected: []string{"cors", "openapi", "jackson", "async", "cache", "auditing", "webmvc"},
			expected: map[string]bool{
				"cors":     true,
				"openapi":  true,
				"jackson":  true,
				"async":    true,
				"cache":    true,
				"auditing": true,
				"webmvc":   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildConfigSelectedMap(tt.selected)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetConfigPackageAllArchitectures(t *testing.T) {
	tests := []struct {
		name        string
		arch        detector.ArchitectureType
		basePackage string
		expected    string
	}{
		{
			name:        "layered",
			arch:        detector.ArchLayered,
			basePackage: "com.example.app",
			expected:    "com.example.app.config",
		},
		{
			name:        "feature",
			arch:        detector.ArchFeature,
			basePackage: "com.example.app",
			expected:    "com.example.app.common.config",
		},
		{
			name:        "hexagonal",
			arch:        detector.ArchHexagonal,
			basePackage: "com.example.app",
			expected:    "com.example.app.infrastructure.config",
		},
		{
			name:        "clean",
			arch:        detector.ArchClean,
			basePackage: "com.example.app",
			expected:    "com.example.app.infrastructure.config",
		},
		{
			name:        "flat",
			arch:        detector.ArchFlat,
			basePackage: "com.example.app",
			expected:    "com.example.app.config",
		},
		{
			name:        "modular",
			arch:        detector.ArchModular,
			basePackage: "com.example.app",
			expected:    "com.example.app.config",
		},
		{
			name:        "unknown",
			arch:        detector.ArchitectureType("unknown"),
			basePackage: "com.example.app",
			expected:    "com.example.app.config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &detector.ProjectProfile{
				Architecture: tt.arch,
				BasePackage:  tt.basePackage,
			}
			result := getConfigPackage(profile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildConfigTemplateData(t *testing.T) {
	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.myapp",
		Architecture: detector.ArchLayered,
		Lombok:       detector.LombokProfile{Detected: true},
	}

	data := buildConfigTemplateData(profile, "com.example.myapp.config")

	assert.Equal(t, "com.example.myapp", data["BasePackage"])
	assert.Equal(t, "com.example.myapp.config", data["ConfigPackage"])
	assert.Equal(t, "Myapp", data["AppName"])
	assert.Equal(t, true, data["HasLombok"])
	assert.Equal(t, "layered", data["Architecture"])
}

func TestExtractAppName(t *testing.T) {
	tests := []struct {
		name        string
		basePackage string
		expected    string
	}{
		{
			name:        "standard package",
			basePackage: "com.example.demo",
			expected:    "Demo",
		},
		{
			name:        "deep package",
			basePackage: "com.company.project.myapp",
			expected:    "Myapp",
		},
		{
			name:        "single segment",
			basePackage: "application",
			expected:    "Application",
		},
		{
			name:        "empty package",
			basePackage: "",
			expected:    "Application",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractAppName(tt.basePackage)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfigTemplateFiles(t *testing.T) {
	expectedFiles := map[string]string{
		"cors":     "CorsConfig.java",
		"openapi":  "OpenApiConfig.java",
		"jackson":  "JacksonConfig.java",
		"async":    "AsyncConfig.java",
		"cache":    "CacheConfig.java",
		"auditing": "AuditingConfig.java",
		"webmvc":   "WebMvcConfig.java",
	}

	for _, opt := range configOptions {
		expected, exists := expectedFiles[opt.Key]
		assert.True(t, exists, "unexpected config key: %s", opt.Key)
		assert.Equal(t, expected, opt.FileName)
	}
}
