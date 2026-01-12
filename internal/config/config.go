package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

const (
	ProjectConfigFile = ".haft.json"
	GlobalConfigDir   = ".config/haft"
	GlobalConfigFile  = "config.json"
)

type ProjectConfig struct {
	Version      string            `json:"version"`
	Project      ProjectSettings   `json:"project"`
	Spring       SpringSettings    `json:"spring"`
	Java         JavaSettings      `json:"java"`
	Build        BuildSettings     `json:"build"`
	Architecture ArchSettings      `json:"architecture"`
	Database     DatabaseSettings  `json:"database"`
	Generators   GeneratorSettings `json:"generators"`
}

type ProjectSettings struct {
	Name        string `json:"name"`
	Group       string `json:"group"`
	Artifact    string `json:"artifact"`
	Description string `json:"description"`
	Package     string `json:"package"`
}

type SpringSettings struct {
	Version string `json:"version"`
}

type JavaSettings struct {
	Version string `json:"version"`
}

type BuildSettings struct {
	Tool string `json:"tool"`
}

type ArchSettings struct {
	Style string `json:"style"`
}

type DatabaseSettings struct {
	Type string `json:"type"`
}

type GeneratorSettings struct {
	DTO   DTOSettings  `json:"dto"`
	Tests TestSettings `json:"tests"`
}

type DTOSettings struct {
	Style string `json:"style"`
}

type TestSettings struct {
	Enabled bool `json:"enabled"`
}

type GlobalConfig struct {
	Defaults DefaultSettings `json:"defaults"`
	Output   OutputSettings  `json:"output"`
}

type DefaultSettings struct {
	JavaVersion  string `json:"java_version"`
	BuildTool    string `json:"build_tool"`
	Architecture string `json:"architecture"`
	SpringBoot   string `json:"spring_boot"`
}

type OutputSettings struct {
	Colors  bool `json:"colors"`
	Verbose bool `json:"verbose"`
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
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (cm *ConfigManager) SaveProjectConfig(config *ProjectConfig) error {
	configPath := filepath.Join(cm.projectDir, ProjectConfigFile)

	data, err := json.MarshalIndent(config, "", "  ")
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
	if err := json.Unmarshal(data, &config); err != nil {
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

	data, err := json.MarshalIndent(config, "", "  ")
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
