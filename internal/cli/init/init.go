package init

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KashifKhn/haft/internal/config"
	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/tui/components"
	"github.com/KashifKhn/haft/internal/tui/styles"
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
	JavaVersion       string
	SpringBootVersion string
	BuildTool         string
	Dependencies      []string
	Architecture      string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new Spring Boot project",
		Long: `Initialize a new Spring Boot project with interactive wizard.

The init command creates a new Spring Boot project with the specified
configuration. If no name is provided, an interactive wizard will guide
you through the setup process.`,
		Example: `  # Interactive mode
  haft init

  # With project name
  haft init my-app

  # With project name in specific directory
  haft init my-app --dir ./projects

  # Non-interactive with all options
  haft init my-app -g com.example -j 21 -s 3.4.0 -b maven --deps web,jpa,lombok --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runInit,
	}

	cmd.Flags().StringP("dir", "d", ".", "Directory to create project in")
	cmd.Flags().StringP("group", "g", "", "Group ID (e.g., com.example)")
	cmd.Flags().StringP("java", "j", "", "Java version (11, 17, 21)")
	cmd.Flags().StringP("spring", "s", "", "Spring Boot version")
	cmd.Flags().StringP("build", "b", "", "Build tool (maven, gradle)")
	cmd.Flags().StringSlice("deps", nil, "Dependencies (web,jpa,security,validation,lombok,h2,postgresql,mysql)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")

	var projectCfg ProjectConfig

	if len(args) > 0 {
		projectCfg.Name = args[0]
		projectCfg.ArtifactId = toArtifactId(args[0])
	}

	if group, _ := cmd.Flags().GetString("group"); group != "" {
		projectCfg.GroupId = group
	}
	if java, _ := cmd.Flags().GetString("java"); java != "" {
		projectCfg.JavaVersion = java
	}
	if spring, _ := cmd.Flags().GetString("spring"); spring != "" {
		projectCfg.SpringBootVersion = spring
	}
	if build, _ := cmd.Flags().GetString("build"); build != "" {
		projectCfg.BuildTool = build
	}
	if deps, _ := cmd.Flags().GetStringSlice("deps"); len(deps) > 0 {
		projectCfg.Dependencies = deps
	}

	if !noInteractive && needsWizard(projectCfg) {
		var err error
		projectCfg, err = runWizard(projectCfg)
		if err != nil {
			return err
		}
	}

	if err := applyDefaults(&projectCfg); err != nil {
		return err
	}

	if err := validateConfig(projectCfg); err != nil {
		return err
	}

	dir, _ := cmd.Flags().GetString("dir")
	projectDir := filepath.Join(dir, projectCfg.ArtifactId)

	return generateProject(projectCfg, projectDir)
}

func needsWizard(cfg ProjectConfig) bool {
	return cfg.Name == "" || cfg.GroupId == "" || cfg.JavaVersion == "" ||
		cfg.SpringBootVersion == "" || cfg.BuildTool == ""
}

func runWizard(cfg ProjectConfig) (ProjectConfig, error) {
	steps := buildWizardSteps(cfg)
	stepKeys := []string{"name", "groupId", "javaVersion", "springVersion", "buildTool", "dependencies"}

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

	if name := wiz.StringValue("name"); name != "" {
		cfg.Name = name
		cfg.ArtifactId = toArtifactId(name)
	}
	if groupId := wiz.StringValue("groupId"); groupId != "" {
		cfg.GroupId = groupId
	}
	if javaVersion := wiz.StringValue("javaVersion"); javaVersion != "" {
		cfg.JavaVersion = javaVersion
	}
	if springVersion := wiz.StringValue("springVersion"); springVersion != "" {
		cfg.SpringBootVersion = springVersion
	}
	if buildTool := wiz.StringValue("buildTool"); buildTool != "" {
		cfg.BuildTool = buildTool
	}
	if deps := wiz.StringSliceValue("dependencies"); len(deps) > 0 {
		cfg.Dependencies = deps
	}

	return cfg, nil
}

func buildWizardSteps(cfg ProjectConfig) []wizard.Step {
	var steps []wizard.Step

	if cfg.Name == "" {
		steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
			Label:       "Project Name",
			Placeholder: "my-spring-app",
			Required:    true,
			Validator:   validateProjectName,
		}))
	}

	if cfg.GroupId == "" {
		steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
			Label:       "Group ID",
			Placeholder: "com.example",
			Required:    true,
			Validator:   validateGroupId,
		}))
	}

	if cfg.JavaVersion == "" {
		steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
			Label: "Java Version",
			Items: []components.SelectItem{
				{Label: "Java 21 (LTS) - Recommended", Value: "21"},
				{Label: "Java 17 (LTS)", Value: "17"},
				{Label: "Java 11 (LTS)", Value: "11"},
			},
		}))
	}

	if cfg.SpringBootVersion == "" {
		steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
			Label: "Spring Boot Version",
			Items: []components.SelectItem{
				{Label: "3.4.0 (Latest)", Value: "3.4.0"},
				{Label: "3.3.6", Value: "3.3.6"},
				{Label: "3.2.12", Value: "3.2.12"},
			},
		}))
	}

	if cfg.BuildTool == "" {
		steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
			Label: "Build Tool",
			Items: []components.SelectItem{
				{Label: "Maven", Value: "maven"},
				{Label: "Gradle (Groovy)", Value: "gradle"},
				{Label: "Gradle (Kotlin)", Value: "gradle-kotlin"},
			},
		}))
	}

	steps = append(steps, wizard.NewMultiSelectStep(components.MultiSelectConfig{
		Label: "Dependencies",
		Items: []components.MultiSelectItem{
			{Label: "Spring Web", Value: "web"},
			{Label: "Spring Data JPA", Value: "jpa"},
			{Label: "Spring Security", Value: "security"},
			{Label: "Spring Validation", Value: "validation"},
			{Label: "Lombok", Value: "lombok"},
			{Label: "H2 Database", Value: "h2"},
			{Label: "PostgreSQL Driver", Value: "postgresql"},
			{Label: "MySQL Driver", Value: "mysql"},
		},
	}))

	return steps
}

func applyDefaults(cfg *ProjectConfig) error {
	homeDir, _ := os.UserHomeDir()
	cm := config.NewConfigManager(afero.NewOsFs(), ".", homeDir)
	globalCfg, _ := cm.LoadGlobalConfig()

	if cfg.JavaVersion == "" {
		cfg.JavaVersion = globalCfg.Defaults.JavaVersion
	}
	if cfg.SpringBootVersion == "" {
		cfg.SpringBootVersion = globalCfg.Defaults.SpringBoot
	}
	if cfg.BuildTool == "" {
		cfg.BuildTool = globalCfg.Defaults.BuildTool
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

func toArtifactId(name string) string {
	name = strings.ToLower(name)
	name = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(name, "-")
	name = regexp.MustCompile(`-+`).ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	return name
}

func generateProject(cfg ProjectConfig, projectDir string) error {
	fs := afero.NewOsFs()

	if exists, _ := afero.DirExists(fs, projectDir); exists {
		return fmt.Errorf("directory %s already exists", projectDir)
	}

	fmt.Println()
	fmt.Printf("  Creating project %s...\n", styles.Focused.Render(cfg.Name))
	fmt.Println()

	engine := generator.NewEngine(fs)

	fmt.Printf("  %s Creating directory structure\n", styles.CheckMark)
	if err := createDirectoryStructure(engine, projectDir, cfg); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	fmt.Printf("  %s Generating project files\n", styles.CheckMark)
	if err := generateProjectFiles(engine, projectDir, cfg); err != nil {
		return fmt.Errorf("failed to generate project files: %w", err)
	}

	if cfg.BuildTool == "maven" {
		fmt.Printf("  %s Adding Maven wrapper\n", styles.CheckMark)
		if err := copyMavenWrapper(engine, projectDir); err != nil {
			return fmt.Errorf("failed to copy Maven wrapper: %w", err)
		}
	}

	fmt.Printf("  %s Writing configuration\n", styles.CheckMark)
	if err := writeHaftConfig(fs, projectDir, cfg); err != nil {
		return fmt.Errorf("failed to write .haft.yaml: %w", err)
	}

	fmt.Println()
	fmt.Printf("  %s Project created successfully!\n", styles.SuccessText.Render("✓"))
	fmt.Println()
	fmt.Println(styles.HelpText.Render("  Next steps:"))
	fmt.Println()
	fmt.Printf("    %s cd %s\n", styles.Arrow, styles.Focused.Render(projectDir))
	if cfg.BuildTool == "maven" {
		fmt.Printf("    %s ./mvnw spring-boot:run\n", styles.Arrow)
	} else {
		fmt.Printf("    %s ./gradlew bootRun\n", styles.Arrow)
	}
	fmt.Println()

	return nil
}

func createDirectoryStructure(engine *generator.Engine, projectDir string, cfg ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(cfg.GroupId, ".", string(os.PathSeparator))
	artifactPath := strings.ReplaceAll(cfg.ArtifactId, "-", "")

	dirs := []string{
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, artifactPath),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, artifactPath, "controller"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, artifactPath, "service"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, artifactPath, "service", "impl"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, artifactPath, "repository"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, artifactPath, "entity"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, artifactPath, "dto"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, artifactPath, "mapper"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, artifactPath, "exception"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, artifactPath, "config"),
		filepath.Join(projectDir, "src", "main", "resources"),
		filepath.Join(projectDir, "src", "test", "java", basePackagePath, artifactPath),
	}

	for _, dir := range dirs {
		if err := engine.GetFS().MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func generateProjectFiles(engine *generator.Engine, projectDir string, cfg ProjectConfig) error {
	basePackage := cfg.GroupId + "." + strings.ReplaceAll(cfg.ArtifactId, "-", "")
	basePackagePath := strings.ReplaceAll(basePackage, ".", string(os.PathSeparator))
	applicationName := toPascalCase(cfg.ArtifactId)

	data := map[string]any{
		"Name":              cfg.Name,
		"GroupId":           cfg.GroupId,
		"ArtifactId":        cfg.ArtifactId,
		"Version":           "0.0.1-SNAPSHOT",
		"Description":       cfg.Description,
		"JavaVersion":       cfg.JavaVersion,
		"SpringBootVersion": cfg.SpringBootVersion,
		"BasePackage":       basePackage,
		"ApplicationName":   applicationName,
		"Dependencies":      buildDependencies(cfg.Dependencies),
		"HasLombok":         contains(cfg.Dependencies, "lombok"),
		"HasJpa":            contains(cfg.Dependencies, "jpa"),
		"HasWeb":            contains(cfg.Dependencies, "web"),
		"HasSecurity":       contains(cfg.Dependencies, "security"),
		"HasValidation":     contains(cfg.Dependencies, "validation"),
	}

	if cfg.BuildTool == "maven" {
		if err := engine.RenderAndWrite(
			"project/pom.xml.tmpl",
			filepath.Join(projectDir, "pom.xml"),
			data,
		); err != nil {
			return err
		}
	}

	if err := engine.RenderAndWrite(
		"project/Application.java.tmpl",
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, applicationName+"Application.java"),
		data,
	); err != nil {
		return err
	}

	if err := engine.RenderAndWrite(
		"project/application.properties.tmpl",
		filepath.Join(projectDir, "src", "main", "resources", "application.properties"),
		data,
	); err != nil {
		return err
	}

	if err := engine.RenderAndWrite(
		"project/ApplicationTests.java.tmpl",
		filepath.Join(projectDir, "src", "test", "java", basePackagePath, applicationName+"ApplicationTests.java"),
		data,
	); err != nil {
		return err
	}

	return nil
}

func writeHaftConfig(fs afero.Fs, projectDir string, cfg ProjectConfig) error {
	basePackage := cfg.GroupId + "." + strings.ReplaceAll(cfg.ArtifactId, "-", "")

	projectCfg := config.ProjectConfig{
		Version: "1",
		Project: config.ProjectSettings{
			Name:        cfg.Name,
			Group:       cfg.GroupId,
			Artifact:    cfg.ArtifactId,
			Description: cfg.Description,
			Package:     basePackage,
		},
		Spring: config.SpringSettings{
			Version: cfg.SpringBootVersion,
		},
		Java: config.JavaSettings{
			Version: cfg.JavaVersion,
		},
		Build: config.BuildSettings{
			Tool: cfg.BuildTool,
		},
		Architecture: config.ArchSettings{
			Style: cfg.Architecture,
		},
		Database: config.DatabaseSettings{
			Type: "h2",
		},
		Generators: config.GeneratorSettings{
			DTO: config.DTOSettings{
				Style: "record",
			},
			Tests: config.TestSettings{
				Enabled: true,
			},
		},
	}

	homeDir, _ := os.UserHomeDir()
	cm := config.NewConfigManager(fs, projectDir, homeDir)
	return cm.SaveProjectConfig(&projectCfg)
}

type Dependency struct {
	GroupId    string
	ArtifactId string
	Version    string
	Scope      string
}

func buildDependencies(deps []string) []Dependency {
	depMap := map[string]Dependency{
		"web": {
			GroupId:    "org.springframework.boot",
			ArtifactId: "spring-boot-starter-web",
		},
		"jpa": {
			GroupId:    "org.springframework.boot",
			ArtifactId: "spring-boot-starter-data-jpa",
		},
		"security": {
			GroupId:    "org.springframework.boot",
			ArtifactId: "spring-boot-starter-security",
		},
		"validation": {
			GroupId:    "org.springframework.boot",
			ArtifactId: "spring-boot-starter-validation",
		},
		"lombok": {
			GroupId:    "org.projectlombok",
			ArtifactId: "lombok",
			Scope:      "provided",
		},
		"h2": {
			GroupId:    "com.h2database",
			ArtifactId: "h2",
			Scope:      "runtime",
		},
		"postgresql": {
			GroupId:    "org.postgresql",
			ArtifactId: "postgresql",
			Scope:      "runtime",
		},
		"mysql": {
			GroupId:    "com.mysql",
			ArtifactId: "mysql-connector-j",
			Scope:      "runtime",
		},
	}

	normalizedDeps := normalizeDependencies(deps)

	var result []Dependency
	for _, dep := range normalizedDeps {
		if d, ok := depMap[dep]; ok {
			result = append(result, d)
		}
	}
	return result
}

func normalizeDependencies(deps []string) []string {
	hasJpa := contains(deps, "jpa")
	hasH2 := contains(deps, "h2")
	hasPostgres := contains(deps, "postgresql")
	hasMysql := contains(deps, "mysql")

	if hasJpa && !hasH2 && !hasPostgres && !hasMysql {
		deps = append(deps, "h2")
	}

	return deps
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func toPascalCase(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}

func copyMavenWrapper(engine *generator.Engine, projectDir string) error {
	wrapperDir := filepath.Join(projectDir, ".mvn", "wrapper")
	if err := engine.GetFS().MkdirAll(wrapperDir, 0755); err != nil {
		return fmt.Errorf("failed to create .mvn/wrapper directory: %w", err)
	}

	propsContent, err := engine.ReadTemplateFile("wrapper/maven-wrapper.properties")
	if err != nil {
		return fmt.Errorf("failed to read maven-wrapper.properties: %w", err)
	}
	propsPath := filepath.Join(wrapperDir, "maven-wrapper.properties")
	if err := engine.WriteFileWithPerm(propsPath, propsContent, 0644); err != nil {
		return fmt.Errorf("failed to write maven-wrapper.properties: %w", err)
	}

	mvnwContent, err := engine.ReadTemplateFile("wrapper/mvnw")
	if err != nil {
		return fmt.Errorf("failed to read mvnw: %w", err)
	}
	mvnwPath := filepath.Join(projectDir, "mvnw")
	if err := engine.WriteFileWithPerm(mvnwPath, mvnwContent, 0755); err != nil {
		return fmt.Errorf("failed to write mvnw: %w", err)
	}

	mvnwCmdContent, err := engine.ReadTemplateFile("wrapper/mvnw.cmd")
	if err != nil {
		return fmt.Errorf("failed to read mvnw.cmd: %w", err)
	}
	mvnwCmdPath := filepath.Join(projectDir, "mvnw.cmd")
	if err := engine.WriteFileWithPerm(mvnwCmdPath, mvnwCmdContent, 0644); err != nil {
		return fmt.Errorf("failed to write mvnw.cmd: %w", err)
	}

	return nil
}
