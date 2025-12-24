package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/maven"
	"github.com/KashifKhn/haft/internal/tui/components"
	"github.com/KashifKhn/haft/internal/tui/wizard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type ResourceConfig struct {
	Name          string
	BasePackage   string
	HasLombok     bool
	HasJpa        bool
	HasValidation bool
}

func newResourceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource [name]",
		Short: "Generate a complete CRUD resource",
		Long: `Generate a complete CRUD resource with all layers.

Creates the following files:
  - Controller (REST endpoints)
  - Service interface
  - ServiceImpl (implementation)
  - Repository (JPA repository)
  - Entity (JPA entity)
  - Request DTO
  - Response DTO
  - Mapper (entity <-> DTO conversion)
  - ResourceNotFoundException (if not exists)

The command auto-detects your project's base package from pom.xml and
checks for Lombok, JPA, and Validation dependencies to customize the
generated code accordingly.`,
		Example: `  # Interactive mode
  haft generate resource

  # With resource name
  haft generate resource user
  haft g resource product

  # Non-interactive with package override
  haft generate resource user --package com.example.myapp --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runResource,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from pom.xml)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("skip-entity", false, "Skip entity generation")
	cmd.Flags().Bool("skip-repository", false, "Skip repository generation")

	return cmd
}

func runResource(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")

	cfg, err := detectProjectConfig()
	if err != nil && noInteractive {
		return fmt.Errorf("could not detect project configuration: %w", err)
	}

	if len(args) > 0 {
		cfg.Name = toPascalCase(args[0])
	}

	if pkg, _ := cmd.Flags().GetString("package"); pkg != "" {
		cfg.BasePackage = pkg
	}

	if !noInteractive {
		cfg, err = runResourceWizard(cfg)
		if err != nil {
			return err
		}
	}

	if err := validateResourceConfig(cfg); err != nil {
		return err
	}

	skipEntity, _ := cmd.Flags().GetBool("skip-entity")
	skipRepository, _ := cmd.Flags().GetBool("skip-repository")

	return generateResource(cfg, skipEntity, skipRepository)
}

func detectProjectConfig() (ResourceConfig, error) {
	var cfg ResourceConfig

	cwd, err := os.Getwd()
	if err != nil {
		return cfg, err
	}

	parser := maven.NewParser()
	pomPath, err := parser.FindPomXml(cwd)
	if err != nil {
		return cfg, err
	}

	project, err := parser.Parse(pomPath)
	if err != nil {
		return cfg, err
	}

	cfg.BasePackage = parser.GetBasePackage(project)
	cfg.HasLombok = parser.HasLombok(project)
	cfg.HasJpa = parser.HasSpringDataJpa(project)
	cfg.HasValidation = parser.HasValidation(project)

	return cfg, nil
}

func runResourceWizard(cfg ResourceConfig) (ResourceConfig, error) {
	steps, stepKeys := buildResourceWizardSteps(cfg)

	w := wizard.New(wizard.WizardConfig{
		Title:    "Generate Resource",
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

	return extractResourceWizardValues(wiz, cfg), nil
}

func buildResourceWizardSteps(cfg ResourceConfig) ([]wizard.Step, []string) {
	var steps []wizard.Step
	var keys []string

	steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
		Label:       "Resource Name",
		Placeholder: "User",
		Default:     cfg.Name,
		Required:    true,
		Validator:   validateResourceName,
		HelpText:    "Name of the resource (e.g., User, Product, Order)",
	}))
	keys = append(keys, "name")

	steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
		Label:       "Base Package",
		Placeholder: "com.example.demo",
		Default:     cfg.BasePackage,
		Required:    true,
		Validator:   validatePackageName,
		HelpText:    "Base package for generated classes (auto-detected from pom.xml)",
	}))
	keys = append(keys, "basePackage")

	steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
		Label: "Use Lombok annotations?",
		Items: []components.SelectItem{
			{Label: "Yes", Value: "yes", Description: "Generate code with Lombok annotations"},
			{Label: "No", Value: "no", Description: "Generate traditional getters/setters"},
		},
		HelpText: fmt.Sprintf("Lombok %s in your project", detectStatusText(cfg.HasLombok)),
	}))
	keys = append(keys, "hasLombok")

	steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
		Label: "Include JPA/Repository layer?",
		Items: []components.SelectItem{
			{Label: "Yes", Value: "yes", Description: "Generate Entity and Repository"},
			{Label: "No", Value: "no", Description: "Skip database layer"},
		},
		HelpText: fmt.Sprintf("Spring Data JPA %s in your project", detectStatusText(cfg.HasJpa)),
	}))
	keys = append(keys, "hasJpa")

	steps = append(steps, wizard.NewSelectStep(components.SelectConfig{
		Label: "Add validation annotations?",
		Items: []components.SelectItem{
			{Label: "Yes", Value: "yes", Description: "Add @Valid to request DTOs"},
			{Label: "No", Value: "no", Description: "Skip validation"},
		},
		HelpText: fmt.Sprintf("Spring Validation %s in your project", detectStatusText(cfg.HasValidation)),
	}))
	keys = append(keys, "hasValidation")

	return steps, keys
}

func detectStatusText(detected bool) string {
	if detected {
		return "detected"
	}
	return "not detected"
}

func extractResourceWizardValues(wiz wizard.WizardModel, cfg ResourceConfig) ResourceConfig {
	if v := wiz.StringValue("name"); v != "" {
		cfg.Name = toPascalCase(v)
	}
	if v := wiz.StringValue("basePackage"); v != "" {
		cfg.BasePackage = v
	}
	if v := wiz.StringValue("hasLombok"); v == "yes" {
		cfg.HasLombok = true
	} else if v == "no" {
		cfg.HasLombok = false
	}
	if v := wiz.StringValue("hasJpa"); v == "yes" {
		cfg.HasJpa = true
	} else if v == "no" {
		cfg.HasJpa = false
	}
	if v := wiz.StringValue("hasValidation"); v == "yes" {
		cfg.HasValidation = true
	} else if v == "no" {
		cfg.HasValidation = false
	}

	return cfg
}

func validateResourceConfig(cfg ResourceConfig) error {
	if cfg.Name == "" {
		return fmt.Errorf("resource name is required")
	}
	if cfg.BasePackage == "" {
		return fmt.Errorf("base package is required")
	}
	return nil
}

func validateResourceName(name string) error {
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`).MatchString(name) {
		return fmt.Errorf("name must start with a letter and contain only letters and numbers")
	}
	return nil
}

func validatePackageName(pkg string) error {
	if pkg == "" {
		return nil
	}
	if !regexp.MustCompile(`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`).MatchString(pkg) {
		return fmt.Errorf("invalid package name format (e.g., com.example.demo)")
	}
	return nil
}

func generateResource(cfg ResourceConfig, skipEntity, skipRepository bool) error {
	log := logger.Default()
	fs := afero.NewOsFs()
	engine := generator.NewEngine(fs)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	srcPath := findSourcePath(cwd)
	if srcPath == "" {
		return fmt.Errorf("could not find src/main/java directory")
	}

	basePath := filepath.Join(srcPath, strings.ReplaceAll(cfg.BasePackage, ".", string(os.PathSeparator)))

	data := buildTemplateData(cfg)

	templates := []struct {
		template   string
		subPackage string
		fileName   string
		skip       bool
	}{
		{"resource/Controller.java.tmpl", "controller", cfg.Name + "Controller.java", false},
		{"resource/Service.java.tmpl", "service", cfg.Name + "Service.java", false},
		{"resource/ServiceImpl.java.tmpl", "service/impl", cfg.Name + "ServiceImpl.java", false},
		{"resource/Repository.java.tmpl", "repository", cfg.Name + "Repository.java", skipRepository || !cfg.HasJpa},
		{"resource/Entity.java.tmpl", "entity", cfg.Name + ".java", skipEntity || !cfg.HasJpa},
		{"resource/Request.java.tmpl", "dto", cfg.Name + "Request.java", false},
		{"resource/Response.java.tmpl", "dto", cfg.Name + "Response.java", false},
		{"resource/Mapper.java.tmpl", "mapper", cfg.Name + "Mapper.java", false},
	}

	if cfg.HasJpa && !skipEntity {
		exceptionPath := filepath.Join(basePath, "exception", "ResourceNotFoundException.java")
		if !engine.FileExists(exceptionPath) {
			templates = append(templates, struct {
				template   string
				subPackage string
				fileName   string
				skip       bool
			}{"resource/ResourceNotFoundException.java.tmpl", "exception", "ResourceNotFoundException.java", false})
		}
	}

	log.Info("Generating resource", "name", cfg.Name)

	for _, t := range templates {
		if t.skip {
			continue
		}

		outputPath := filepath.Join(basePath, t.subPackage, t.fileName)

		if engine.FileExists(outputPath) {
			log.Warning("File already exists, skipping", "path", outputPath)
			continue
		}

		if err := engine.RenderAndWrite(t.template, outputPath, data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", t.fileName, err)
		}

		log.Info("Created", "file", formatRelativePath(cwd, outputPath))
	}

	log.Info("Resource generated successfully", "name", cfg.Name)
	return nil
}

func buildTemplateData(cfg ResourceConfig) map[string]any {
	return map[string]any{
		"Name":          cfg.Name,
		"NameLower":     strings.ToLower(cfg.Name),
		"NameCamel":     toCamelCase(cfg.Name),
		"BasePackage":   cfg.BasePackage,
		"HasLombok":     cfg.HasLombok,
		"HasJpa":        cfg.HasJpa,
		"HasValidation": cfg.HasValidation,
	}
}

func findSourcePath(startDir string) string {
	candidates := []string{
		filepath.Join(startDir, "src", "main", "java"),
		filepath.Join(startDir, "app", "src", "main", "java"),
	}

	for _, path := range candidates {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path
		}
	}

	return ""
}

func formatRelativePath(base, path string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	return rel
}

func toPascalCase(s string) string {
	words := splitWords(s)
	var result string
	for _, word := range words {
		result += capitalize(strings.ToLower(word))
	}
	return result
}

func toCamelCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return s
	}
	result := strings.ToLower(words[0])
	for _, word := range words[1:] {
		result += capitalize(strings.ToLower(word))
	}
	return result
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func splitWords(s string) []string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	var words []string
	var currentWord strings.Builder

	for i, r := range s {
		if r == ' ' {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			continue
		}

		if i > 0 && isUpper(r) && currentWord.Len() > 0 {
			lastChar := []rune(currentWord.String())[currentWord.Len()-1]
			if !isUpper(lastChar) {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		}

		currentWord.WriteRune(r)
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}
