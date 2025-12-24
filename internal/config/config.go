package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const (
	ProjectConfigFile = ".haft.yaml"
	GlobalConfigDir   = ".config/haft"
	GlobalConfigFile  = "config.yaml"
)

type ProjectConfig struct {
	Version      string            `yaml:"version"`
	Project      ProjectSettings   `yaml:"project"`
	Spring       SpringSettings    `yaml:"spring"`
	Java         JavaSettings      `yaml:"java"`
	Build        BuildSettings     `yaml:"build"`
	Architecture ArchSettings      `yaml:"architecture"`
	Database     DatabaseSettings  `yaml:"database"`
	Generators   GeneratorSettings `yaml:"generators"`
}

type ProjectSettings struct {
	Name        string `yaml:"name"`
	Group       string `yaml:"group"`
	Artifact    string `yaml:"artifact"`
	Description string `yaml:"description"`
	Package     string `yaml:"package"`
}

type SpringSettings struct {
	Version string `yaml:"version"`
}

type JavaSettings struct {
	Version string `yaml:"version"`
}

type BuildSettings struct {
	Tool string `yaml:"tool"`
}

type ArchSettings struct {
	Style string `yaml:"style"`
}

type DatabaseSettings struct {
	Type string `yaml:"type"`
}

type GeneratorSettings struct {
	DTO   DTOSettings  `yaml:"dto"`
	Tests TestSettings `yaml:"tests"`
}

type DTOSettings struct {
	Style string `yaml:"style"`
}

type TestSettings struct {
	Enabled bool `yaml:"enabled"`
}

type GlobalConfig struct {
	Defaults DefaultSettings `yaml:"defaults"`
	Output   OutputSettings  `yaml:"output"`
}

type DefaultSettings struct {
	JavaVersion  string `yaml:"java_version"`
	BuildTool    string `yaml:"build_tool"`
	Architecture string `yaml:"architecture"`
	SpringBoot   string `yaml:"spring_boot"`
}

type OutputSettings struct {
	Colors  bool `yaml:"colors"`
	Verbose bool `yaml:"verbose"`
}

type ConfigManager struct {
	fs         afero.Fs
	projectDir string
	homeDir    string
}

func NewConfigManager(fs afero.Fs, projectDir, homeDir string) *ConfigManager {
	return &ConfigManager{
		fs:         fs,
		projectDir: projectDir,
		homeDir:    homeDir,
	}
}

func NewDefaultConfigManager() (*ConfigManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &ConfigManager{
		fs:         afero.NewOsFs(),
		projectDir: workDir,
		homeDir:    homeDir,
	}, nil
}

func (cm *ConfigManager) LoadProjectConfig() (*ProjectConfig, error) {
	configPath := filepath.Join(cm.projectDir, ProjectConfigFile)

	data, err := afero.ReadFile(cm.fs, configPath)
	if err != nil {
		return nil, err
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (cm *ConfigManager) SaveProjectConfig(config *ProjectConfig) error {
	configPath := filepath.Join(cm.projectDir, ProjectConfigFile)

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return afero.WriteFile(cm.fs, configPath, data, 0644)
}

func (cm *ConfigManager) ProjectConfigExists() bool {
	configPath := filepath.Join(cm.projectDir, ProjectConfigFile)
	exists, _ := afero.Exists(cm.fs, configPath)
	return exists
}

func (cm *ConfigManager) LoadGlobalConfig() (*GlobalConfig, error) {
	configPath := filepath.Join(cm.homeDir, GlobalConfigDir, GlobalConfigFile)

	data, err := afero.ReadFile(cm.fs, configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultGlobalConfig(), nil
		}
		return nil, err
	}

	var config GlobalConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (cm *ConfigManager) SaveGlobalConfig(config *GlobalConfig) error {
	configDir := filepath.Join(cm.homeDir, GlobalConfigDir)
	configPath := filepath.Join(configDir, GlobalConfigFile)

	if err := cm.fs.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return afero.WriteFile(cm.fs, configPath, data, 0644)
}

func (cm *ConfigManager) GlobalConfigExists() bool {
	configPath := filepath.Join(cm.homeDir, GlobalConfigDir, GlobalConfigFile)
	exists, _ := afero.Exists(cm.fs, configPath)
	return exists
}

func DefaultProjectConfig() *ProjectConfig {
	return &ProjectConfig{
		Version: "1",
		Project: ProjectSettings{},
		Spring: SpringSettings{
			Version: "3.4.0",
		},
		Java: JavaSettings{
			Version: "21",
		},
		Build: BuildSettings{
			Tool: "maven",
		},
		Architecture: ArchSettings{
			Style: "layered",
		},
		Database: DatabaseSettings{
			Type: "postgresql",
		},
		Generators: GeneratorSettings{
			DTO: DTOSettings{
				Style: "record",
			},
			Tests: TestSettings{
				Enabled: true,
			},
		},
	}
}

func DefaultGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		Defaults: DefaultSettings{
			JavaVersion:  "21",
			BuildTool:    "maven",
			Architecture: "layered",
			SpringBoot:   "3.4.0",
		},
		Output: OutputSettings{
			Colors:  true,
			Verbose: false,
		},
	}
}

func (cm *ConfigManager) SetProjectDir(dir string) {
	cm.projectDir = dir
}

func (cm *ConfigManager) GetProjectDir() string {
	return cm.projectDir
}

func (cm *ConfigManager) GetHomeDir() string {
	return cm.homeDir
}
