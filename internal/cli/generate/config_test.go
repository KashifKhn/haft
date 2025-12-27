package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KashifKhn/haft/internal/detector"
	"github.com/KashifKhn/haft/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestBuildConfigTemplateDataNoLombok(t *testing.T) {
	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.demo",
		Architecture: detector.ArchFeature,
		Lombok:       detector.LombokProfile{Detected: false},
	}

	data := buildConfigTemplateData(profile, "com.example.demo.common.config")

	assert.Equal(t, "com.example.demo", data["BasePackage"])
	assert.Equal(t, "com.example.demo.common.config", data["ConfigPackage"])
	assert.Equal(t, "Demo", data["AppName"])
	assert.Equal(t, false, data["HasLombok"])
	assert.Equal(t, "feature", data["Architecture"])
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

func TestConfigMultiSelectWrapperInit(t *testing.T) {
	items := []components.MultiSelectItem{
		{Label: "Test", Value: "test", Selected: false},
	}
	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label: "Test",
		Items: items,
	})
	wrapper := configMultiSelectWrapper{model: model}

	cmd := wrapper.Init()
	assert.Nil(t, cmd)
}

func TestConfigMultiSelectWrapperView(t *testing.T) {
	items := []components.MultiSelectItem{
		{Label: "Test", Value: "test", Selected: false},
	}
	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label: "Test Label",
		Items: items,
	})
	wrapper := configMultiSelectWrapper{model: model}

	view := wrapper.View()
	assert.NotEmpty(t, view)
}

func TestConfigMultiSelectWrapperUpdateCtrlC(t *testing.T) {
	items := []components.MultiSelectItem{
		{Label: "Test", Value: "test", Selected: false},
	}
	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label: "Test",
		Items: items,
	})
	wrapper := configMultiSelectWrapper{model: model}

	newModel, cmd := wrapper.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.NotNil(t, newModel)
	assert.NotNil(t, cmd)
}

func TestConfigMultiSelectWrapperUpdateRegularKey(t *testing.T) {
	items := []components.MultiSelectItem{
		{Label: "Test", Value: "test", Selected: false},
	}
	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label: "Test",
		Items: items,
	})
	wrapper := configMultiSelectWrapper{model: model}

	newModel, _ := wrapper.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.NotNil(t, newModel)
}

func TestConfigSelectionStruct(t *testing.T) {
	selection := configSelection{
		BasePackage: "com.example.app",
		Selected:    []string{"cors", "jackson"},
	}

	assert.Equal(t, "com.example.app", selection.BasePackage)
	assert.Len(t, selection.Selected, 2)
	assert.Contains(t, selection.Selected, "cors")
	assert.Contains(t, selection.Selected, "jackson")
}

func TestConfigOptionStruct(t *testing.T) {
	opt := configOption{
		Name:        "CORS",
		FileName:    "CorsConfig.java",
		Key:         "cors",
		Description: "Cross-origin resource sharing",
	}

	assert.Equal(t, "CORS", opt.Name)
	assert.Equal(t, "CorsConfig.java", opt.FileName)
	assert.Equal(t, "cors", opt.Key)
	assert.Equal(t, "Cross-origin resource sharing", opt.Description)
}

func TestGenerateConfigsNoSourcePath(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-config-nosrc-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.demo",
		Architecture: detector.ArchLayered,
	}

	selection := configSelection{
		Selected: []string{"cors"},
	}

	err = generateConfigs(profile, selection)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not find src/main/java directory")
}

func TestGenerateConfigsIntegration(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-config-integration-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.demo",
		Architecture: detector.ArchLayered,
		Lombok:       detector.LombokProfile{Detected: true},
	}

	selection := configSelection{
		Selected: []string{"cors", "jackson"},
	}

	err = generateConfigs(profile, selection)
	require.NoError(t, err)

	configPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "config")
	assert.DirExists(t, configPath)

	corsFile := filepath.Join(configPath, "CorsConfig.java")
	assert.FileExists(t, corsFile)

	jacksonFile := filepath.Join(configPath, "JacksonConfig.java")
	assert.FileExists(t, jacksonFile)

	openApiFile := filepath.Join(configPath, "OpenApiConfig.java")
	assert.NoFileExists(t, openApiFile)
}

func TestGenerateConfigsSkipsExisting(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-config-skip-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "config")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	existingFile := filepath.Join(srcPath, "CorsConfig.java")
	require.NoError(t, os.WriteFile(existingFile, []byte("existing content"), 0644))

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.demo",
		Architecture: detector.ArchLayered,
	}

	selection := configSelection{
		Selected: []string{"cors"},
	}

	err = generateConfigs(profile, selection)
	require.NoError(t, err)

	content, err := os.ReadFile(existingFile)
	require.NoError(t, err)
	assert.Equal(t, "existing content", string(content))
}

func TestRunConfigNoInteractiveWithoutAll(t *testing.T) {
	cmd := newConfigCommand()
	cmd.SetArgs([]string{"--no-interactive", "--package", "com.example.app"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "use --all flag")
}

func TestRunConfigMissingPackage(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-config-nopkg-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	require.NoError(t, os.Chdir(tmpDir))

	cmd := newConfigCommand()
	cmd.SetArgs([]string{"--no-interactive", "--all"})

	err = cmd.Execute()
	assert.Error(t, err)
}
