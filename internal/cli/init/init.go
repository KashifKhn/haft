package init

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KashifKhn/haft/internal/config"
	"github.com/KashifKhn/haft/internal/initializr"
	"github.com/KashifKhn/haft/internal/tui/components"
	"github.com/KashifKhn/haft/internal/tui/wizard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type ProjectConfig struct {
	Name              string
	GroupId           string
	ArtifactId        string
	Description       string
	PackageName       string
	JavaVersion       string
	SpringBootVersion string
	BuildTool         string
	Packaging         string
	ConfigFormat      string
	Dependencies      []string
	Architecture      string
	InitGit           bool
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new Spring Boot project",
		Long: `Initialize a new Spring Boot project with interactive wizard.

The init command creates a new Spring Boot project with the specified
configuration. If no name is provided, an interactive wizard will guide
you through the setup process.

The wizard presents all dependencies from Spring Initializr organized by
category (Web, SQL, NoSQL, Security, etc.) with descriptions and search.`,
		Example: `  # Interactive mode (recommended)
  haft init

  # With project name
  haft init my-app

  # Non-interactive with all options
  haft init my-app -g com.example -j 21 -s 3.4.0 -b maven --deps web,data-jpa,lombok --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runInit,
	}

	cmd.Flags().StringP("dir", "d", ".", "Directory to create project in")
	cmd.Flags().StringP("group", "g", "", "Group ID (e.g., com.example)")
	cmd.Flags().StringP("artifact", "a", "", "Artifact ID")
	cmd.Flags().String("description", "", "Project description")
	cmd.Flags().String("package", "", "Base package name")
	cmd.Flags().StringP("java", "j", "", "Java version (17, 21, 25)")
	cmd.Flags().StringP("spring", "s", "", "Spring Boot version")
	cmd.Flags().StringP("build", "b", "", "Build tool (maven, gradle, gradle-kotlin)")
	cmd.Flags().String("packaging", "", "Packaging type (jar, war)")
	cmd.Flags().String("config", "", "Config format (properties, yaml)")
	cmd.Flags().StringSlice("deps", nil, "Dependencies (comma-separated IDs)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("json", false, "Output result as JSON")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	var cfg ProjectConfig

	if len(args) > 0 {
		cfg.Name = args[0]
		cfg.ArtifactId = toArtifactId(args[0])
	}

	if v, _ := cmd.Flags().GetString("group"); v != "" {
		cfg.GroupId = v
	}
	if v, _ := cmd.Flags().GetString("artifact"); v != "" {
		cfg.ArtifactId = v
	}
	if v, _ := cmd.Flags().GetString("description"); v != "" {
		cfg.Description = v
	}
	if v, _ := cmd.Flags().GetString("package"); v != "" {
		cfg.PackageName = v
	}
	if v, _ := cmd.Flags().GetString("java"); v != "" {
		cfg.JavaVersion = v
	}
	if v, _ := cmd.Flags().GetString("spring"); v != "" {
		cfg.SpringBootVersion = v
	}
	if v, _ := cmd.Flags().GetString("build"); v != "" {
		cfg.BuildTool = v
	}
	if v, _ := cmd.Flags().GetString("packaging"); v != "" {
		cfg.Packaging = v
	}
	if v, _ := cmd.Flags().GetString("config"); v != "" {
		cfg.ConfigFormat = v
	}
	if deps, _ := cmd.Flags().GetStringSlice("deps"); len(deps) > 0 {
		cfg.Dependencies = deps
	}

	if !noInteractive {
		var err error
		cfg, err = runWizard(cfg)
		if err != nil {
			return err
		}
	}

	if err := applyDefaults(&cfg); err != nil {
		return err
	}

	if err := validateConfig(cfg); err != nil {
		return err
	}

	dir, _ := cmd.Flags().GetString("dir")
	projectDir := filepath.Join(dir, cfg.ArtifactId)

	return generateProject(cfg, projectDir, jsonOutput)
}

func runWizard(cfg ProjectConfig) (ProjectConfig, error) {
	steps, stepKeys := buildWizardSteps(cfg)

	w := wizard.New(wizard.WizardConfig{
		Title:    "Create New Spring Boot Project",
		Steps:    steps,
		StepKeys: stepKeys,
	})

	p := tea.NewProgram(w)
	finalModel, err := p.Run()
	if err != nil {
		return cfg, fmt.Errorf("wizard failed: %w", err)
	}

	wiz, ok := finalModel.(wizard.WizardModel)
	if !ok {
		return cfg, fmt.Errorf("unexpected wizard state")
	}

	if wiz.Cancelled() {
		return cfg, fmt.Errorf("wizard cancelled")
	}

	return extractWizardValues(wiz, cfg), nil
}

func buildWizardSteps(cfg ProjectConfig) ([]wizard.Step, []string) {
	var steps []wizard.Step
	var keys []string

	steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
		Label:       "Project Name",
		Placeholder: "demo",
		Default:     cfg.Name,
		Required:    true,
		Validator:   validateProjectName,
		HelpText:    "The display name for your project",
	}))
	keys = append(keys, "name")

	steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
		Label:       "Group ID",
		Placeholder: "com.example",
		Default:     cfg.GroupId,
		Required:    true,
		Validator:   validateGroupId,
		HelpText:    "The group ID for your project (e.g., com.company)",
	}))
	keys = append(keys, "groupId")

	steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
		Label:       "Artifact ID",
		Placeholder: "demo",
		Default:     cfg.ArtifactId,
		Required:    true,
		Validator:   validateArtifactId,
		HelpText:    "The artifact ID for your project (used in pom.xml)",
	}))
	keys = append(keys, "artifactId")

	steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
		Label:       "Description",
		Placeholder: "Demo project for Spring Boot",
		Default:     cfg.Description,
		HelpText:    "A short description of your project",
	}))
	keys = append(keys, "description")

	steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
		Label:       "Package Name",
		Placeholder: "com.example.demo",
		Default:     cfg.PackageName,
		Validator:   validatePackageName,
		HelpText:    "Base package for your Java classes (auto-generated from Group ID + Artifact ID)",
		DynamicDefault: func(values map[string]any) string {
			groupId, _ := values["groupId"].(string)
			artifactId, _ := values["artifactId"].(string)
			if groupId != "" && artifactId != "" {
				cleanArtifact := strings.ReplaceAll(artifactId, "-", "")
				return groupId + "." + cleanArtifact
			}
			return ""
		},
	}))
	keys = append(keys, "packageName")

	javaItems := buildJavaVersionItems()
	steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
		Label:    "Java Version",
		Items:    javaItems,
		HelpText: "Select the Java version for your project",
	}))
	keys = append(keys, "javaVersion")

	bootItems := buildBootVersionItems()
	steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
		Label:    "Spring Boot Version",
		Items:    bootItems,
		HelpText: "Select the Spring Boot version",
	}))
	keys = append(keys, "springVersion")

	steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
		Label: "Build Tool",
		Items: []components.SelectItem{
			{Label: "Maven", Value: "maven", Description: "Build with Apache Maven"},
			{Label: "Gradle - Groovy", Value: "gradle", Description: "Build with Gradle using Groovy DSL"},
			{Label: "Gradle - Kotlin", Value: "gradle-kotlin", Description: "Build with Gradle using Kotlin DSL"},
		},
		HelpText: "Select the build tool for your project",
	}))
	keys = append(keys, "buildTool")

	steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
		Label: "Packaging",
		Items: []components.SelectItem{
			{Label: "Jar", Value: "jar", Description: "Package as an executable JAR"},
			{Label: "War", Value: "war", Description: "Package as a WAR for deployment to servlet container"},
		},
		HelpText: "Select the packaging type",
	}))
	keys = append(keys, "packaging")

	steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
		Label: "Configuration Format",
		Items: []components.SelectItem{
			{Label: "Properties", Value: "properties", Description: "Use application.properties format"},
			{Label: "YAML", Value: "yaml", Description: "Use application.yml format"},
		},
		Default:  "yaml",
		HelpText: "Select the configuration file format",
	}))
	keys = append(keys, "configFormat")

	depCategories := buildDepCategories()
	steps = append(steps, wizard.NewDepPickerStep(components.DepPickerConfig{
		Label:      "Dependencies",
		Categories: depCategories,
	}))
	keys = append(keys, "dependencies")

	steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
		Label: "Initialize Git Repository?",
		Items: []components.SelectItem{
			{Label: "Yes", Value: "yes", Description: "Initialize git and create initial commit"},
			{Label: "No", Value: "no", Description: "Skip git initialization"},
		},
		HelpText: "Create a git repository in your project",
	}))
	keys = append(keys, "initGit")

	return steps, keys
}

func buildJavaVersionItems() []components.SelectItem {
	versions, err := initializr.GetJavaVersions()
	if err != nil || len(versions) == 0 {
		return []components.SelectItem{
			{Label: "Java 21 (LTS)", Value: "21", Description: "Long Term Support version"},
			{Label: "Java 17 (LTS)", Value: "17", Description: "Long Term Support version"},
		}
	}

	var items []components.SelectItem
	for _, v := range versions {
		label := fmt.Sprintf("Java %s", v.Name)
		desc := ""
		switch v.ID {
		case "21":
			label = fmt.Sprintf("Java %s (LTS)", v.Name)
			desc = "Long Term Support - Recommended"
		case "17":
			label = fmt.Sprintf("Java %s (LTS)", v.Name)
			desc = "Long Term Support"
		case "25":
			desc = "Latest version"
		}
		items = append(items, components.SelectItem{
			Label:       label,
			Value:       v.ID,
			Description: desc,
		})
	}
	return items
}

func buildBootVersionItems() []components.SelectItem {
	versions, err := initializr.GetBootVersions()
	if err != nil || len(versions) == 0 {
		return []components.SelectItem{
			{Label: "3.4.0 (Latest)", Value: "3.4.0"},
			{Label: "3.3.6", Value: "3.3.6"},
		}
	}

	var items []components.SelectItem
	for _, v := range versions {
		if strings.Contains(v.Name, "SNAPSHOT") {
			continue
		}
		items = append(items, components.SelectItem{
			Label: v.Name,
			Value: v.ID,
		})
	}
	return items
}

func buildDepCategories() []components.DepCategory {
	categories, err := initializr.GetDependencyCategories()
	if err != nil {
		return getDefaultDepCategories()
	}

	var result []components.DepCategory
	for _, cat := range categories {
		depCat := components.DepCategory{Name: cat.Name}
		for _, dep := range cat.Values {
			depCat.Dependencies = append(depCat.Dependencies, components.DepItem{
				ID:          dep.ID,
				Name:        dep.Name,
				Description: dep.Description,
			})
		}
		if len(depCat.Dependencies) > 0 {
			result = append(result, depCat)
		}
	}
	return result
}

func getDefaultDepCategories() []components.DepCategory {
	return []components.DepCategory{
		{
			Name: "Web",
			Dependencies: []components.DepItem{
				{ID: "web", Name: "Spring Web", Description: "Build web applications using Spring MVC"},
				{ID: "webflux", Name: "Spring Reactive Web", Description: "Build reactive web applications"},
			},
		},
		{
			Name: "SQL",
			Dependencies: []components.DepItem{
				{ID: "data-jpa", Name: "Spring Data JPA", Description: "Persist data with JPA"},
				{ID: "h2", Name: "H2 Database", Description: "In-memory database for development"},
				{ID: "postgresql", Name: "PostgreSQL Driver", Description: "PostgreSQL JDBC driver"},
				{ID: "mysql", Name: "MySQL Driver", Description: "MySQL JDBC driver"},
			},
		},
		{
			Name: "Developer Tools",
			Dependencies: []components.DepItem{
				{ID: "devtools", Name: "Spring Boot DevTools", Description: "Fast application restarts and LiveReload"},
				{ID: "lombok", Name: "Lombok", Description: "Reduce boilerplate code"},
			},
		},
		{
			Name: "Security",
			Dependencies: []components.DepItem{
				{ID: "security", Name: "Spring Security", Description: "Secure your application"},
			},
		},
	}
}

func extractWizardValues(wiz wizard.WizardModel, cfg ProjectConfig) ProjectConfig {
	if v := wiz.StringValue("name"); v != "" {
		cfg.Name = v
	}
	if v := wiz.StringValue("groupId"); v != "" {
		cfg.GroupId = v
	}
	if v := wiz.StringValue("artifactId"); v != "" {
		cfg.ArtifactId = v
	} else if cfg.ArtifactId == "" && cfg.Name != "" {
		cfg.ArtifactId = toArtifactId(cfg.Name)
	}
	if v := wiz.StringValue("description"); v != "" {
		cfg.Description = v
	}
	if v := wiz.StringValue("packageName"); v != "" {
		cfg.PackageName = v
	}
	if v := wiz.StringValue("javaVersion"); v != "" {
		cfg.JavaVersion = v
	}
	if v := wiz.StringValue("springVersion"); v != "" {
		cfg.SpringBootVersion = v
	}
	if v := wiz.StringValue("buildTool"); v != "" {
		cfg.BuildTool = v
	}
	if v := wiz.StringValue("packaging"); v != "" {
		cfg.Packaging = v
	}
	if v := wiz.StringValue("configFormat"); v != "" {
		cfg.ConfigFormat = v
	}
	if deps := wiz.StringSliceValue("dependencies"); len(deps) > 0 {
		cfg.Dependencies = deps
	}
	if v := wiz.StringValue("initGit"); v == "yes" {
		cfg.InitGit = true
	}

	return cfg
}

func applyDefaults(cfg *ProjectConfig) error {
	homeDir, _ := os.UserHomeDir()
	cm := config.NewConfigManager(afero.NewOsFs(), ".", homeDir)
	globalCfg, _ := cm.LoadGlobalConfig()

	if cfg.ArtifactId == "" && cfg.Name != "" {
		cfg.ArtifactId = toArtifactId(cfg.Name)
	}
	if cfg.PackageName == "" && cfg.GroupId != "" && cfg.ArtifactId != "" {
		cfg.PackageName = cfg.GroupId + "." + strings.ReplaceAll(cfg.ArtifactId, "-", "")
	}
	if cfg.JavaVersion == "" {
		if globalCfg.Defaults.JavaVersion != "" {
			cfg.JavaVersion = globalCfg.Defaults.JavaVersion
		} else {
			cfg.JavaVersion = "21"
		}
	}
	if cfg.SpringBootVersion == "" {
		if globalCfg.Defaults.SpringBoot != "" {
			cfg.SpringBootVersion = globalCfg.Defaults.SpringBoot
		} else {
			cfg.SpringBootVersion = "3.4.0"
		}
	}
	if cfg.BuildTool == "" {
		if globalCfg.Defaults.BuildTool != "" {
			cfg.BuildTool = globalCfg.Defaults.BuildTool
		} else {
			cfg.BuildTool = "maven"
		}
	}
	if cfg.Packaging == "" {
		cfg.Packaging = "jar"
	}
	if cfg.ConfigFormat == "" {
		cfg.ConfigFormat = "yaml"
	}
	if cfg.Architecture == "" {
		cfg.Architecture = globalCfg.Defaults.Architecture
	}
	if cfg.Description == "" {
		cfg.Description = fmt.Sprintf("%s Spring Boot application", cfg.Name)
	}

	return nil
}

func validateConfig(cfg ProjectConfig) error {
	if cfg.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if cfg.GroupId == "" {
		return fmt.Errorf("group ID is required")
	}
	if cfg.ArtifactId == "" {
		return fmt.Errorf("artifact ID is required")
	}
	return nil
}

func validateProjectName(name string) error {
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_]*$`).MatchString(name) {
		return fmt.Errorf("name must start with a letter and contain only letters, numbers, hyphens, and underscores")
	}
	return nil
}

func validateGroupId(groupId string) error {
	if !regexp.MustCompile(`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`).MatchString(groupId) {
		return fmt.Errorf("invalid group ID format (e.g., com.example)")
	}
	return nil
}

func validateArtifactId(artifactId string) error {
	if len(artifactId) < 2 {
		return fmt.Errorf("artifact ID must be at least 2 characters")
	}
	if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(artifactId) {
		return fmt.Errorf("artifact ID must be lowercase and contain only letters, numbers, and hyphens")
	}
	return nil
}

func validatePackageName(pkg string) error {
	if pkg == "" {
		return nil
	}
	if !regexp.MustCompile(`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`).MatchString(pkg) {
		return fmt.Errorf("invalid package name format")
	}
	return nil
}

func toArtifactId(name string) string {
	name = strings.ToLower(name)
	name = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(name, "-")
	name = regexp.MustCompile(`-+`).ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	return name
}
