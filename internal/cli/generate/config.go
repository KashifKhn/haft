package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KashifKhn/haft/internal/detector"
	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type configSelection struct {
	BasePackage string
	Selected    []string
}

type configOption struct {
	Name        string
	FileName    string
	Key         string
	Description string
}

var configOptions = []configOption{
	{"CORS", "CorsConfig.java", "cors", "Cross-origin resource sharing"},
	{"OpenAPI", "OpenApiConfig.java", "openapi", "Swagger/OpenAPI documentation"},
	{"Jackson", "JacksonConfig.java", "jackson", "JSON serialization settings"},
	{"Async", "AsyncConfig.java", "async", "Async/thread pool configuration"},
	{"Caching", "CacheConfig.java", "cache", "Spring Cache configuration"},
	{"Auditing", "AuditingConfig.java", "auditing", "JPA auditing (@CreatedDate, etc.)"},
	{"WebMvc", "WebMvcConfig.java", "webmvc", "Web MVC customization"},
}

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Generate configuration classes",
		Long: `Generate common Spring Boot configuration classes.

Available configurations:
  - CorsConfig.java — Cross-origin resource sharing setup
  - OpenApiConfig.java — Swagger/OpenAPI documentation
  - JacksonConfig.java — JSON serialization settings
  - AsyncConfig.java — Async/thread pool configuration
  - CacheConfig.java — Spring Cache configuration
  - AuditingConfig.java — JPA auditing (@CreatedDate, @LastModifiedDate)
  - WebMvcConfig.java — Web MVC customization (resource handlers)

All configurations are optional; select which ones you need via the interactive picker.`,
		Example: `  # Interactive picker to select configurations
  haft generate config
  haft g cfg

  # Generate all configurations
  haft generate config --all

  # Override base package
  haft generate config --package com.example.app`,
		RunE: runConfig,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from project)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("all", false, "Generate all configuration classes")
	cmd.Flags().Bool("refresh", false, "Force re-detection of project profile (ignore cache)")

	return cmd
}

func runConfig(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	includeAll, _ := cmd.Flags().GetBool("all")
	forceRefresh, _ := cmd.Flags().GetBool("refresh")
	log := logger.Default()

	profile, err := DetectProjectProfileWithRefresh(forceRefresh)
	if err != nil {
		if noInteractive {
			return fmt.Errorf("could not detect project profile: %w", err)
		}
		log.Warning("Could not detect project profile, using defaults")
		profile = &detector.ProjectProfile{
			Architecture: detector.ArchLayered,
		}
	}

	enrichProfileFromBuildFile(profile)

	if pkg, _ := cmd.Flags().GetString("package"); pkg != "" {
		profile.BasePackage = pkg
	}

	selection := configSelection{}

	if includeAll {
		for _, cfg := range configOptions {
			selection.Selected = append(selection.Selected, cfg.Key)
		}
	}

	if !noInteractive {
		wizardResult, err := runConfigWizard(profile.BasePackage, includeAll)
		if err != nil {
			return err
		}
		if wizardResult.BasePackage != "" {
			profile.BasePackage = wizardResult.BasePackage
		}
		if !includeAll {
			selection.Selected = wizardResult.Selected
		}
	}

	if profile.BasePackage == "" {
		return fmt.Errorf("base package could not be detected. Use --package flag to specify it (e.g., --package com.example.myapp)")
	}

	if len(selection.Selected) == 0 && !noInteractive {
		log.Info("No configurations selected")
		return nil
	}

	if len(selection.Selected) == 0 && noInteractive && !includeAll {
		return fmt.Errorf("use --all flag or run without --no-interactive to select configurations")
	}

	return generateConfigs(profile, selection)
}

func runConfigWizard(currentPackage string, skipPicker bool) (configSelection, error) {
	cfg := configSelection{}

	componentCfg := ComponentConfig{
		BasePackage: currentPackage,
		Name:        "Config",
	}

	result, err := RunComponentWizard("Generate Configuration Classes", componentCfg, "Config")
	if err != nil {
		return cfg, err
	}
	cfg.BasePackage = result.BasePackage

	if skipPicker {
		return cfg, nil
	}

	cfg.Selected, err = runConfigPicker()
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func runConfigPicker() ([]string, error) {
	items := make([]components.MultiSelectItem, len(configOptions))
	for i, opt := range configOptions {
		items[i] = components.MultiSelectItem{
			Label:    fmt.Sprintf("%s — %s", opt.Name, opt.Description),
			Value:    opt.Key,
			Selected: false,
		}
	}

	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label:    "Select configurations to generate",
		Items:    items,
		Required: false,
	})

	wrapper := configMultiSelectWrapper{model: model}
	p := tea.NewProgram(wrapper)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(configMultiSelectWrapper)
	if result.model.GoBack() {
		return nil, fmt.Errorf("wizard cancelled")
	}

	return result.model.Values(), nil
}

type configMultiSelectWrapper struct {
	model components.MultiSelectModel
}

func (w configMultiSelectWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w configMultiSelectWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return w, tea.Quit
		}
	}

	updated, cmd := w.model.Update(msg)
	w.model = updated

	if w.model.Submitted() || w.model.GoBack() {
		return w, tea.Quit
	}

	return w, cmd
}

func (w configMultiSelectWrapper) View() string {
	return w.model.View()
}

func generateConfigs(profile *detector.ProjectProfile, selection configSelection) error {
	log := logger.Default()
	fs := afero.NewOsFs()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	engine := generator.NewEngineWithLoader(fs, cwd)

	srcPath := FindSourcePath(cwd)
	if srcPath == "" {
		return fmt.Errorf("could not find src/main/java directory")
	}

	configPackage := getConfigPackage(profile)
	packagePath := strings.ReplaceAll(configPackage, ".", string(os.PathSeparator))
	basePath := filepath.Join(srcPath, packagePath)

	selectedMap := buildConfigSelectedMap(selection.Selected)
	data := buildConfigTemplateData(profile, configPackage)

	log.Info("Generating configuration classes", "package", configPackage)

	generatedCount := 0
	skippedCount := 0

	for _, opt := range configOptions {
		if !selectedMap[opt.Key] {
			continue
		}

		templatePath := "config/" + opt.FileName + ".tmpl"
		outputPath := filepath.Join(basePath, opt.FileName)

		if engine.FileExists(outputPath) {
			log.Warning("File exists, skipping", "file", FormatRelativePath(cwd, outputPath))
			skippedCount++
			continue
		}

		if err := engine.RenderAndWrite(templatePath, outputPath, data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", opt.FileName, err)
		}

		log.Info("Created", "file", FormatRelativePath(cwd, outputPath))
		generatedCount++
	}

	if generatedCount > 0 {
		log.Success(fmt.Sprintf("Generated %d configuration files", generatedCount))
	}
	if skippedCount > 0 {
		log.Info(fmt.Sprintf("Skipped %d existing files", skippedCount))
	}

	return nil
}

func buildConfigSelectedMap(selected []string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range selected {
		m[s] = true
	}
	return m
}

func getConfigPackage(profile *detector.ProjectProfile) string {
	switch profile.Architecture {
	case detector.ArchFeature:
		return profile.BasePackage + ".common.config"
	case detector.ArchHexagonal:
		return profile.BasePackage + ".infrastructure.config"
	case detector.ArchClean:
		return profile.BasePackage + ".infrastructure.config"
	default:
		return profile.BasePackage + ".config"
	}
}

func buildConfigTemplateData(profile *detector.ProjectProfile, configPackage string) map[string]any {
	appName := extractAppName(profile.BasePackage)

	return map[string]any{
		"BasePackage":   profile.BasePackage,
		"ConfigPackage": configPackage,
		"AppName":       appName,
		"HasLombok":     profile.Lombok.Detected,
		"Architecture":  string(profile.Architecture),
	}
}

func extractAppName(basePackage string) string {
	parts := strings.Split(basePackage, ".")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		if len(name) > 0 {
			return strings.ToUpper(string(name[0])) + name[1:]
		}
	}
	return "Application"
}
