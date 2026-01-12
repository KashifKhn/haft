package config

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigManager(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	assert.NotNil(t, cm)
	assert.Equal(t, "/project", cm.GetProjectDir())
	assert.Equal(t, "/home/user", cm.GetHomeDir())
}

func TestSaveAndLoadProjectConfig(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	config := &ProjectConfig{
		Version: "1",
		Project: ProjectSettings{
			Name:     "test-app",
			Group:    "com.example",
			Artifact: "test-app",
		},
		Spring: SpringSettings{
			Version: "3.4.0",
		},
		Java: JavaSettings{
			Version: "21",
		},
		Build: BuildSettings{
			Tool: "maven",
		},
	}

	err := cm.SaveProjectConfig(config)
	require.NoError(t, err)

	loaded, err := cm.LoadProjectConfig()
	require.NoError(t, err)

	assert.Equal(t, config.Version, loaded.Version)
	assert.Equal(t, config.Project.Name, loaded.Project.Name)
	assert.Equal(t, config.Project.Group, loaded.Project.Group)
	assert.Equal(t, config.Spring.Version, loaded.Spring.Version)
	assert.Equal(t, config.Java.Version, loaded.Java.Version)
	assert.Equal(t, config.Build.Tool, loaded.Build.Tool)
}

func TestLoadProjectConfigNotExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	_, err := cm.LoadProjectConfig()
	assert.Error(t, err)
}

func TestProjectConfigExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	assert.False(t, cm.ProjectConfigExists())

	err := cm.SaveProjectConfig(DefaultProjectConfig())
	require.NoError(t, err)

	assert.True(t, cm.ProjectConfigExists())
}

func TestSaveAndLoadGlobalConfig(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	config := &GlobalConfig{
		Defaults: DefaultSettings{
			JavaVersion:  "17",
			BuildTool:    "gradle",
			Architecture: "hexagonal",
			SpringBoot:   "3.3.0",
		},
		Output: OutputSettings{
			Colors:  false,
			Verbose: true,
		},
	}

	err := cm.SaveGlobalConfig(config)
	require.NoError(t, err)

	loaded, err := cm.LoadGlobalConfig()
	require.NoError(t, err)

	assert.Equal(t, config.Defaults.JavaVersion, loaded.Defaults.JavaVersion)
	assert.Equal(t, config.Defaults.BuildTool, loaded.Defaults.BuildTool)
	assert.Equal(t, config.Defaults.Architecture, loaded.Defaults.Architecture)
	assert.Equal(t, config.Output.Colors, loaded.Output.Colors)
	assert.Equal(t, config.Output.Verbose, loaded.Output.Verbose)
}

func TestLoadGlobalConfigReturnsDefaultWhenNotExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	loaded, err := cm.LoadGlobalConfig()
	require.NoError(t, err)

	defaultConfig := DefaultGlobalConfig()
	assert.Equal(t, defaultConfig.Defaults.JavaVersion, loaded.Defaults.JavaVersion)
	assert.Equal(t, defaultConfig.Defaults.BuildTool, loaded.Defaults.BuildTool)
}

func TestGlobalConfigExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	assert.False(t, cm.GlobalConfigExists())

	err := cm.SaveGlobalConfig(DefaultGlobalConfig())
	require.NoError(t, err)

	assert.True(t, cm.GlobalConfigExists())
}

func TestGlobalConfigCreatesDirectory(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	err := cm.SaveGlobalConfig(DefaultGlobalConfig())
	require.NoError(t, err)

	configDir := filepath.Join("/home/user", GlobalConfigDir)
	exists, err := afero.DirExists(fs, configDir)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestDefaultProjectConfig(t *testing.T) {
	config := DefaultProjectConfig()

	assert.Equal(t, "1", config.Version)
	assert.Equal(t, "3.4.0", config.Spring.Version)
	assert.Equal(t, "21", config.Java.Version)
	assert.Equal(t, "maven", config.Build.Tool)
	assert.Equal(t, "layered", config.Architecture.Style)
	assert.Equal(t, "postgresql", config.Database.Type)
	assert.Equal(t, "record", config.Generators.DTO.Style)
	assert.True(t, config.Generators.Tests.Enabled)
}

func TestDefaultGlobalConfig(t *testing.T) {
	config := DefaultGlobalConfig()

	assert.Equal(t, "21", config.Defaults.JavaVersion)
	assert.Equal(t, "maven", config.Defaults.BuildTool)
	assert.Equal(t, "layered", config.Defaults.Architecture)
	assert.Equal(t, "3.4.0", config.Defaults.SpringBoot)
	assert.True(t, config.Output.Colors)
	assert.False(t, config.Output.Verbose)
}

func TestSetProjectDir(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	cm.SetProjectDir("/new/project")
	assert.Equal(t, "/new/project", cm.GetProjectDir())
}

func TestLoadProjectConfigInvalidYAML(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	configPath := filepath.Join("/project", ProjectConfigFile)
	err := afero.WriteFile(fs, configPath, []byte("{invalid json content"), 0644)
	require.NoError(t, err)

	_, err = cm.LoadProjectConfig()
	assert.Error(t, err)
}

func TestLoadGlobalConfigInvalidYAML(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	configDir := filepath.Join("/home/user", GlobalConfigDir)
	err := fs.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, GlobalConfigFile)
	err = afero.WriteFile(fs, configPath, []byte("{invalid json content"), 0644)
	require.NoError(t, err)

	_, err = cm.LoadGlobalConfig()
	assert.Error(t, err)
}

func TestNewDefaultConfigManager(t *testing.T) {
	cm, err := NewDefaultConfigManager()

	assert.NoError(t, err)
	assert.NotNil(t, cm)
	assert.NotEmpty(t, cm.GetProjectDir())
	assert.NotEmpty(t, cm.GetHomeDir())
}

func TestSaveProjectConfigMarshalError(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	config := DefaultProjectConfig()
	err := cm.SaveProjectConfig(config)

	assert.NoError(t, err)
}

func TestSaveGlobalConfigMarshalError(t *testing.T) {
	fs := afero.NewMemMapFs()
	cm := NewConfigManager(fs, "/project", "/home/user")

	config := DefaultGlobalConfig()
	err := cm.SaveGlobalConfig(config)

	assert.NoError(t, err)
}

func TestSaveProjectConfigWriteError(t *testing.T) {
	fs := afero.NewReadOnlyFs(afero.NewMemMapFs())
	cm := NewConfigManager(fs, "/project", "/home/user")

	config := DefaultProjectConfig()
	err := cm.SaveProjectConfig(config)

	assert.Error(t, err)
}

func TestSaveGlobalConfigDirError(t *testing.T) {
	fs := afero.NewReadOnlyFs(afero.NewMemMapFs())
	cm := NewConfigManager(fs, "/project", "/home/user")

	config := DefaultGlobalConfig()
	err := cm.SaveGlobalConfig(config)

	assert.Error(t, err)
}
