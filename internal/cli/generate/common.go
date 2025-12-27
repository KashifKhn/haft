package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/detector"
	"github.com/KashifKhn/haft/internal/generator"
	_ "github.com/KashifKhn/haft/internal/gradle"
	"github.com/KashifKhn/haft/internal/logger"
	_ "github.com/KashifKhn/haft/internal/maven"
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

	fs := afero.NewOsFs()
	result, err := buildtool.Detect(cwd, fs)
	if err != nil {
		return cfg, err
	}

	project, err := result.Parser.Parse(result.FilePath)
	if err != nil {
		return cfg, err
	}

	cfg.BasePackage = result.Parser.GetBasePackage(project)
	cfg.HasLombok = result.Parser.HasLombok(project)
	cfg.HasJpa = result.Parser.HasSpringDataJpa(project)
	cfg.HasValidation = result.Parser.HasValidation(project)

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
		HelpText:    "Base package for generated classes (auto-detected from build file)",
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
	if err := ValidateComponentName(cfg.Name); err != nil {
		return err
	}
	if cfg.BasePackage == "" {
		return fmt.Errorf("base package is required")
	}
	if err := ValidatePackageName(cfg.BasePackage); err != nil {
		return err
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

func GenerateComponent(cfg ComponentConfig, templateName, subPackage, fileNamePattern string) (bool, error) {
	log := logger.Default()
	fs := afero.NewOsFs()

	cwd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	engine := generator.NewEngineWithLoader(fs, cwd)

	srcPath := FindSourcePath(cwd)
	if srcPath == "" {
		return false, fmt.Errorf("could not find src/main/java directory")
	}

	basePath := filepath.Join(srcPath, strings.ReplaceAll(cfg.BasePackage, ".", string(os.PathSeparator)))
	data := BuildTemplateData(cfg)

	fileName := strings.ReplaceAll(fileNamePattern, "{Name}", cfg.Name)
	outputPath := filepath.Join(basePath, subPackage, fileName)

	if engine.FileExists(outputPath) {
		log.Warning("Skipped (already exists)", "file", FormatRelativePath(cwd, outputPath))
		return false, nil
	}

	if err := engine.RenderAndWrite(templateName, outputPath, data); err != nil {
		return false, fmt.Errorf("failed to generate %s: %w", fileName, err)
	}

	log.Success("Created", "file", FormatRelativePath(cwd, outputPath))
	return true, nil
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

func DetectProjectProfile() (*detector.ProjectProfile, error) {
	return DetectProjectProfileWithRefresh(false)
}

func DetectProjectProfileWithRefresh(forceRefresh bool) (*detector.ProjectProfile, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	log := logger.Default()
	cache := detector.NewProfileCache(cwd)

	if !forceRefresh && cache.IsValid() {
		profile, err := cache.Load()
		if err == nil && profile != nil {
			log.Debug("Using cached profile from .haft/profile.yaml")
			return profile, nil
		}
	}

	log.Debug("Scanning project to detect architecture...")
	d := detector.NewDetector(cwd)
	profile, err := d.Detect()
	if err != nil {
		return nil, err
	}

	if err := cache.Save(profile); err != nil {
		log.Debug("Failed to cache profile", "error", err.Error())
	} else {
		log.Debug("Profile cached to .haft/profile.yaml")
	}

	return profile, nil
}

type TemplateContext struct {
	Name      string
	NameLower string
	NameCamel string

	BasePackage    string
	FeaturePackage string
	TestPackage    string

	Architecture string

	HasLombok     bool
	HasJpa        bool
	HasValidation bool
	HasSwagger    bool
	HasMapStruct  bool

	IDType      string
	IDImport    string
	TestIdValue string

	ControllerSuffix string
	RequestSuffix    string
	ResponseSuffix   string

	HasBaseEntity    bool
	BaseEntityName   string
	BaseEntityImport string

	HasResponseWrapper    bool
	ResponseWrapperName   string
	ResponseWrapperImport string

	HasGlobalException bool
	ExceptionPackage   string

	ValidationImport string

	Lombok detector.LombokProfile
}

func BuildTemplateContextFromProfile(name string, profile *detector.ProjectProfile) TemplateContext {
	nameLower := strings.ToLower(name)

	featurePackage := profile.BasePackage
	testPackage := profile.BasePackage
	if profile.Architecture == detector.ArchFeature {
		featurePackage = profile.BasePackage + "." + nameLower
		testPackage = profile.BasePackage + "." + nameLower
	}

	testIdValue := "1L"
	if profile.IDType == "UUID" {
		testIdValue = "UUID.randomUUID()"
	}

	ctx := TemplateContext{
		Name:      name,
		NameLower: nameLower,
		NameCamel: ToCamelCase(name),

		BasePackage:    profile.BasePackage,
		FeaturePackage: featurePackage,
		TestPackage:    testPackage,

		Architecture: string(profile.Architecture),

		HasLombok:     profile.Lombok.Detected,
		HasJpa:        profile.Database == detector.DatabaseJPA || profile.Database == detector.DatabaseMulti,
		HasValidation: profile.HasValidation,
		HasSwagger:    profile.HasSwagger,
		HasMapStruct:  profile.Mapper == detector.MapperMapStruct,

		IDType:      profile.IDType,
		IDImport:    profile.GetIDImport(),
		TestIdValue: testIdValue,

		ControllerSuffix: profile.ControllerSuffix,
		RequestSuffix:    name + profile.GetDTORequestSuffix(),
		ResponseSuffix:   name + profile.GetDTOResponseSuffix(),

		Lombok: profile.Lombok,
	}

	if profile.BaseEntity != nil {
		ctx.HasBaseEntity = true
		ctx.BaseEntityName = profile.BaseEntity.Name
		ctx.BaseEntityImport = profile.GetBaseEntityImport()
	}

	if profile.ResponseWrapper != nil {
		ctx.HasResponseWrapper = true
		ctx.ResponseWrapperName = profile.ResponseWrapper.Name
		ctx.ResponseWrapperImport = profile.GetResponseWrapperImport()
	}

	if profile.Exceptions.HasGlobalHandler {
		ctx.HasGlobalException = true
		ctx.ExceptionPackage = profile.BasePackage + ".exception"
		if profile.Architecture == detector.ArchFeature {
			ctx.ExceptionPackage = profile.BasePackage + ".common.exception"
		}
	}

	if profile.ValidationStyle == detector.ValidationJakarta {
		ctx.ValidationImport = "jakarta.validation"
	} else if profile.ValidationStyle == detector.ValidationJavax {
		ctx.ValidationImport = "javax.validation"
	}

	return ctx
}

func (ctx TemplateContext) ToMap() map[string]any {
	return map[string]any{
		"Name":                  ctx.Name,
		"NameLower":             ctx.NameLower,
		"NameCamel":             ctx.NameCamel,
		"BasePackage":           ctx.BasePackage,
		"FeaturePackage":        ctx.FeaturePackage,
		"TestPackage":           ctx.TestPackage,
		"Architecture":          ctx.Architecture,
		"HasLombok":             ctx.HasLombok,
		"HasJpa":                ctx.HasJpa,
		"HasValidation":         ctx.HasValidation,
		"HasSwagger":            ctx.HasSwagger,
		"HasMapStruct":          ctx.HasMapStruct,
		"IDType":                ctx.IDType,
		"IDImport":              ctx.IDImport,
		"TestIdValue":           ctx.TestIdValue,
		"ControllerSuffix":      ctx.ControllerSuffix,
		"RequestSuffix":         ctx.RequestSuffix,
		"ResponseSuffix":        ctx.ResponseSuffix,
		"HasBaseEntity":         ctx.HasBaseEntity,
		"BaseEntityName":        ctx.BaseEntityName,
		"BaseEntityImport":      ctx.BaseEntityImport,
		"HasResponseWrapper":    ctx.HasResponseWrapper,
		"ResponseWrapperName":   ctx.ResponseWrapperName,
		"ResponseWrapperImport": ctx.ResponseWrapperImport,
		"HasGlobalException":    ctx.HasGlobalException,
		"ExceptionPackage":      ctx.ExceptionPackage,
		"ValidationImport":      ctx.ValidationImport,
		"Lombok":                ctx.Lombok,
	}
}

func GetTemplateDir(profile *detector.ProjectProfile) string {
	switch profile.Architecture {
	case detector.ArchFeature:
		return "resource/feature"
	case detector.ArchHexagonal:
		return "resource/feature"
	case detector.ArchClean:
		return "resource/feature"
	default:
		return "resource/layered"
	}
}

func GetTestTemplateDir(profile *detector.ProjectProfile) string {
	switch profile.Architecture {
	case detector.ArchFeature:
		return "test/feature"
	case detector.ArchHexagonal:
		return "test/feature"
	case detector.ArchClean:
		return "test/feature"
	default:
		return "test/layered"
	}
}

func FindTestPath(startDir string) string {
	candidates := []string{
		filepath.Join(startDir, "src", "test", "java"),
		filepath.Join(startDir, "app", "src", "test", "java"),
	}

	for _, path := range candidates {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path
		}
	}

	return ""
}
