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
)

type ComponentConfig struct {
	Name          string
	BasePackage   string
	HasLombok     bool
	HasJpa        bool
	HasValidation bool
}

func DetectProjectConfig() (ComponentConfig, error) {
	var cfg ComponentConfig

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

func RunComponentWizard(title string, cfg ComponentConfig, componentType string) (ComponentConfig, error) {
	steps, stepKeys := buildComponentWizardSteps(cfg, componentType)

	w := wizard.New(wizard.WizardConfig{
		Title:    title,
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

	return extractComponentWizardValues(wiz, cfg), nil
}

func buildComponentWizardSteps(cfg ComponentConfig, componentType string) ([]wizard.Step, []string) {
	var steps []wizard.Step
	var keys []string

	steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
		Label:       componentType + " Name",
		Placeholder: "User",
		Default:     cfg.Name,
		Required:    true,
		Validator:   ValidateComponentName,
		HelpText:    fmt.Sprintf("Name of the %s (e.g., User, Product, Order)", strings.ToLower(componentType)),
	}))
	keys = append(keys, "name")

	steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
		Label:       "Base Package",
		Placeholder: "com.example.demo",
		Default:     cfg.BasePackage,
		Required:    true,
		Validator:   ValidatePackageName,
		HelpText:    "Base package for generated classes (auto-detected from pom.xml)",
	}))
	keys = append(keys, "basePackage")

	return steps, keys
}

func extractComponentWizardValues(wiz wizard.WizardModel, cfg ComponentConfig) ComponentConfig {
	if v := wiz.StringValue("name"); v != "" {
		cfg.Name = ToPascalCase(v)
	}
	if v := wiz.StringValue("basePackage"); v != "" {
		cfg.BasePackage = v
	}

	return cfg
}

func ValidateComponentConfig(cfg ComponentConfig) error {
	if cfg.Name == "" {
		return fmt.Errorf("name is required")
	}
	if cfg.BasePackage == "" {
		return fmt.Errorf("base package is required")
	}
	return nil
}

func ValidateComponentName(name string) error {
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`).MatchString(name) {
		return fmt.Errorf("name must start with a letter and contain only letters and numbers")
	}
	return nil
}

func ValidatePackageName(pkg string) error {
	if pkg == "" {
		return nil
	}
	if !regexp.MustCompile(`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`).MatchString(pkg) {
		return fmt.Errorf("invalid package name format (e.g., com.example.demo)")
	}
	return nil
}

func GenerateComponent(cfg ComponentConfig, templateName, subPackage, fileNamePattern string) error {
	log := logger.Default()
	fs := afero.NewOsFs()
	engine := generator.NewEngine(fs)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	srcPath := FindSourcePath(cwd)
	if srcPath == "" {
		return fmt.Errorf("could not find src/main/java directory")
	}

	basePath := filepath.Join(srcPath, strings.ReplaceAll(cfg.BasePackage, ".", string(os.PathSeparator)))
	data := BuildTemplateData(cfg)

	fileName := strings.ReplaceAll(fileNamePattern, "{Name}", cfg.Name)
	outputPath := filepath.Join(basePath, subPackage, fileName)

	if engine.FileExists(outputPath) {
		return fmt.Errorf("file already exists: %s", FormatRelativePath(cwd, outputPath))
	}

	if err := engine.RenderAndWrite(templateName, outputPath, data); err != nil {
		return fmt.Errorf("failed to generate %s: %w", fileName, err)
	}

	log.Success("Created", "file", FormatRelativePath(cwd, outputPath))
	return nil
}

func BuildTemplateData(cfg ComponentConfig) map[string]any {
	return map[string]any{
		"Name":          cfg.Name,
		"NameLower":     strings.ToLower(cfg.Name),
		"NameCamel":     ToCamelCase(cfg.Name),
		"BasePackage":   cfg.BasePackage,
		"HasLombok":     cfg.HasLombok,
		"HasJpa":        cfg.HasJpa,
		"HasValidation": cfg.HasValidation,
	}
}

func FindSourcePath(startDir string) string {
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

func FormatRelativePath(base, path string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	return rel
}

func ToPascalCase(s string) string {
	words := SplitWords(s)
	var result string
	for _, word := range words {
		result += Capitalize(strings.ToLower(word))
	}
	return result
}

func ToCamelCase(s string) string {
	words := SplitWords(s)
	if len(words) == 0 {
		return s
	}
	result := strings.ToLower(words[0])
	for _, word := range words[1:] {
		result += Capitalize(strings.ToLower(word))
	}
	return result
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func SplitWords(s string) []string {
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

		if i > 0 && IsUpper(r) && currentWord.Len() > 0 {
			lastChar := []rune(currentWord.String())[currentWord.Len()-1]
			if !IsUpper(lastChar) {
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

func IsUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}
