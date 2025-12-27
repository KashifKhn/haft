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
		Use:     "resource [name]",
		Aliases: []string{"r"},
		Short:   "Generate a complete CRUD resource",
		Long: `Generate a complete CRUD resource with all layers.

Creates the following files:
  - Controller (REST endpoints)
  - Service interface
  - ServiceImpl (implementation)
  - Repository (JPA repository) - if JPA detected
  - Entity (JPA entity) - if JPA detected
  - Request DTO
  - Response DTO
  - Mapper (entity <-> DTO conversion)

The command intelligently detects your project's architecture pattern and
generates code that matches your existing conventions:
  - Base package and feature modules
  - Lombok annotations (@Data, @Builder, etc.)
  - DTO naming style (Request/Response vs DTO)
  - ID type (UUID vs Long)
  - Mapper type (MapStruct vs manual)
  - Base entity inheritance
  - Response wrapper patterns`,
		Example: `  # Interactive mode
  haft generate resource

  # With resource name
  haft generate resource user
  haft g r product

  # Non-interactive with package override
  haft generate resource user --package com.example.myapp --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runResource,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from project)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("skip-entity", false, "Skip entity generation")
	cmd.Flags().Bool("skip-repository", false, "Skip repository generation")
	cmd.Flags().Bool("skip-tests", false, "Skip test generation")
	cmd.Flags().Bool("legacy", false, "Use legacy layered generation (ignores architecture detection)")

	return cmd
}

func runResource(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	useLegacy, _ := cmd.Flags().GetBool("legacy")
	log := logger.Default()

	if useLegacy {
		return runLegacyResource(cmd, args)
	}

	profile, err := DetectProjectProfile()
	if err != nil {
		log.Warning("Could not detect project profile, falling back to legacy mode")
		return runLegacyResource(cmd, args)
	}

	var resourceName string
	if len(args) > 0 {
		resourceName = ToPascalCase(args[0])
	}

	if pkg, _ := cmd.Flags().GetString("package"); pkg != "" {
		profile.BasePackage = pkg
	}

	if !noInteractive {
		var wizErr error
		resourceName, wizErr = runResourceNameWizard(resourceName)
		if wizErr != nil {
			return wizErr
		}
	}

	if resourceName == "" {
		return fmt.Errorf("resource name is required")
	}

	if profile.BasePackage == "" {
		return fmt.Errorf("base package is required")
	}

	skipEntity, _ := cmd.Flags().GetBool("skip-entity")
	skipRepository, _ := cmd.Flags().GetBool("skip-repository")
	skipTests, _ := cmd.Flags().GetBool("skip-tests")

	return generateResourceWithProfile(resourceName, profile, skipEntity, skipRepository, skipTests)
}

func runResourceNameWizard(currentName string) (string, error) {
	steps, keys := []wizard.Step{}, []string{}

	steps = append(steps, wizard.NewTextInputStep(components.TextInputConfig{
		Label:       "Resource Name",
		Placeholder: "User",
		Default:     currentName,
		Required:    true,
		Validator:   ValidateComponentName,
		HelpText:    "Name of the resource (e.g., User, Product, Order)",
	}))
	keys = append(keys, "name")

	w := wizard.New(wizard.WizardConfig{
		Title:    "Generate Resource",
		Steps:    steps,
		StepKeys: keys,
	})

	p := tea.NewProgram(w)
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("wizard failed: %w", err)
	}

	wiz, ok := finalModel.(wizard.WizardModel)
	if !ok {
		return "", fmt.Errorf("unexpected wizard state")
	}

	if wiz.Cancelled() {
		return "", fmt.Errorf("wizard cancelled")
	}

	return ToPascalCase(wiz.StringValue("name")), nil
}

func generateResourceWithProfile(name string, profile *detector.ProjectProfile, skipEntity, skipRepository, skipTests bool) error {
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

	ctx := BuildTemplateContextFromProfile(name, profile)
	templateDir := GetTemplateDir(profile)
	data := ctx.ToMap()

	log.Info("Generating resource", "name", name, "architecture", profile.Architecture)
	log.Debug("Template directory", "dir", templateDir)

	if profile.Lombok.Detected {
		log.Debug("Using Lombok annotations")
	}
	if ctx.HasJpa {
		log.Debug("Generating JPA Entity and Repository")
	}
	if profile.HasSwagger {
		log.Debug("Adding Swagger/OpenAPI annotations")
	}
	if ctx.HasMapStruct {
		log.Debug("Using MapStruct for mapping")
	}
	if ctx.HasBaseEntity {
		log.Debug("Extending base entity", "base", ctx.BaseEntityName)
	}

	templates := buildTemplateList(name, profile, templateDir, ctx, skipEntity, skipRepository)

	generatedCount := 0
	skippedCount := 0

	for _, t := range templates {
		if t.skip {
			continue
		}

		outputPath := computeOutputPath(srcPath, profile, name, t.subPackage, t.fileName)

		if engine.FileExists(outputPath) {
			log.Warning("File exists, skipping", "file", FormatRelativePath(cwd, outputPath))
			skippedCount++
			continue
		}

		if err := engine.RenderAndWrite(t.template, outputPath, data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", t.fileName, err)
		}

		log.Info("Created", "file", FormatRelativePath(cwd, outputPath))
		generatedCount++
	}

	if !skipTests {
		testCount, testSkipped, err := generateTestsWithProfile(name, profile, ctx, skipEntity, skipRepository)
		if err != nil {
			log.Warning("Failed to generate tests", "error", err.Error())
		} else {
			generatedCount += testCount
			skippedCount += testSkipped
		}
	}

	if generatedCount > 0 {
		log.Success(fmt.Sprintf("Generated %d files for %s resource", generatedCount, name))
	}
	if skippedCount > 0 {
		log.Info(fmt.Sprintf("Skipped %d existing files", skippedCount))
	}

	return nil
}

func generateTestsWithProfile(name string, profile *detector.ProjectProfile, ctx TemplateContext, skipEntity, skipRepository bool) (int, int, error) {
	log := logger.Default()
	fs := afero.NewOsFs()
	engine := generator.NewEngine(fs)

	cwd, err := os.Getwd()
	if err != nil {
		return 0, 0, err
	}

	testPath := FindTestPath(cwd)
	if testPath == "" {
		return 0, 0, fmt.Errorf("could not find src/test/java directory")
	}

	testTemplateDir := GetTestTemplateDir(profile)
	data := ctx.ToMap()

	log.Debug("Generating tests", "template_dir", testTemplateDir)

	testTemplates := buildTestTemplateList(name, profile, testTemplateDir, ctx, skipEntity, skipRepository)

	generatedCount := 0
	skippedCount := 0

	for _, t := range testTemplates {
		if t.skip {
			continue
		}

		outputPath := computeTestOutputPath(testPath, profile, name, t.subPackage, t.fileName)

		if engine.FileExists(outputPath) {
			log.Warning("Test file exists, skipping", "file", FormatRelativePath(cwd, outputPath))
			skippedCount++
			continue
		}

		if err := engine.RenderAndWrite(t.template, outputPath, data); err != nil {
			return generatedCount, skippedCount, fmt.Errorf("failed to generate %s: %w", t.fileName, err)
		}

		log.Info("Created test", "file", FormatRelativePath(cwd, outputPath))
		generatedCount++
	}

	return generatedCount, skippedCount, nil
}

type templateSpec struct {
	template   string
	subPackage string
	fileName   string
	skip       bool
}

func buildTemplateList(name string, profile *detector.ProjectProfile, templateDir string, ctx TemplateContext, skipEntity, skipRepository bool) []templateSpec {
	controllerSuffix := profile.ControllerSuffix
	requestSuffix := profile.GetDTORequestSuffix()
	responseSuffix := profile.GetDTOResponseSuffix()

	hasJpa := ctx.HasJpa

	templates := []templateSpec{
		{templateDir + "/Controller.java.tmpl", "controller", name + controllerSuffix + ".java", false},
		{templateDir + "/Service.java.tmpl", "service", name + "Service.java", false},
		{templateDir + "/ServiceImpl.java.tmpl", "service/impl", name + "ServiceImpl.java", false},
		{templateDir + "/Repository.java.tmpl", "repository", name + "Repository.java", skipRepository || !hasJpa},
		{templateDir + "/Entity.java.tmpl", "entity", name + ".java", skipEntity || !hasJpa},
		{templateDir + "/Request.java.tmpl", "dto", name + requestSuffix + ".java", false},
		{templateDir + "/Response.java.tmpl", "dto", name + responseSuffix + ".java", false},
		{templateDir + "/Mapper.java.tmpl", "mapper", name + "Mapper.java", false},
	}

	return templates
}

func computeOutputPath(srcPath string, profile *detector.ProjectProfile, resourceName, subPackage, fileName string) string {
	resourceLower := strings.ToLower(resourceName)

	var packagePath string
	switch profile.Architecture {
	case detector.ArchFeature:
		if profile.FeatureStyle == detector.FeatureStyleFlat {
			if subPackage == "dto" {
				packagePath = filepath.Join(
					strings.ReplaceAll(profile.BasePackage, ".", string(os.PathSeparator)),
					resourceLower,
					"dto",
				)
			} else {
				packagePath = filepath.Join(
					strings.ReplaceAll(profile.BasePackage, ".", string(os.PathSeparator)),
					resourceLower,
				)
			}
		} else {
			packagePath = filepath.Join(
				strings.ReplaceAll(profile.BasePackage, ".", string(os.PathSeparator)),
				resourceLower,
				subPackage,
			)
		}
	case detector.ArchHexagonal, detector.ArchClean:
		packagePath = filepath.Join(
			strings.ReplaceAll(profile.BasePackage, ".", string(os.PathSeparator)),
			resourceLower,
			subPackage,
		)
	default:
		packagePath = filepath.Join(
			strings.ReplaceAll(profile.BasePackage, ".", string(os.PathSeparator)),
			subPackage,
		)
	}

	return filepath.Join(srcPath, packagePath, fileName)
}

func buildTestTemplateList(name string, profile *detector.ProjectProfile, testTemplateDir string, ctx TemplateContext, skipEntity, skipRepository bool) []templateSpec {
	hasJpa := ctx.HasJpa

	templates := []templateSpec{
		{testTemplateDir + "/ServiceTest.java.tmpl", "service", name + "ServiceTest.java", false},
		{testTemplateDir + "/ControllerTest.java.tmpl", "controller", name + "ControllerTest.java", false},
		{testTemplateDir + "/RepositoryTest.java.tmpl", "repository", name + "RepositoryTest.java", skipRepository || !hasJpa},
		{testTemplateDir + "/EntityTest.java.tmpl", "entity", name + "Test.java", skipEntity || !hasJpa},
	}

	return templates
}

func computeTestOutputPath(testPath string, profile *detector.ProjectProfile, resourceName, subPackage, fileName string) string {
	resourceLower := strings.ToLower(resourceName)

	var packagePath string
	switch profile.Architecture {
	case detector.ArchFeature:
		if profile.FeatureStyle == detector.FeatureStyleFlat {
			if subPackage == "dto" {
				packagePath = filepath.Join(
					strings.ReplaceAll(profile.BasePackage, ".", string(os.PathSeparator)),
					resourceLower,
					"dto",
				)
			} else {
				packagePath = filepath.Join(
					strings.ReplaceAll(profile.BasePackage, ".", string(os.PathSeparator)),
					resourceLower,
				)
			}
		} else {
			packagePath = filepath.Join(
				strings.ReplaceAll(profile.BasePackage, ".", string(os.PathSeparator)),
				resourceLower,
				subPackage,
			)
		}
	case detector.ArchHexagonal, detector.ArchClean:
		packagePath = filepath.Join(
			strings.ReplaceAll(profile.BasePackage, ".", string(os.PathSeparator)),
			resourceLower,
			subPackage,
		)
	default:
		packagePath = filepath.Join(
			strings.ReplaceAll(profile.BasePackage, ".", string(os.PathSeparator)),
			subPackage,
		)
	}

	return filepath.Join(testPath, packagePath, fileName)
}

func runLegacyResource(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	log := logger.Default()

	compCfg, err := DetectProjectConfig()
	if err != nil {
		if noInteractive {
			return fmt.Errorf("could not detect project configuration: %w", err)
		}
		log.Warning("Could not detect project config, using defaults")
	}

	cfg := ResourceConfig(compCfg)

	if len(args) > 0 {
		cfg.Name = ToPascalCase(args[0])
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

func runResourceWizard(cfg ResourceConfig) (ResourceConfig, error) {
	compCfg := ComponentConfig(cfg)

	result, err := RunComponentWizard("Generate Resource", compCfg, "Resource")
	if err != nil {
		return cfg, err
	}

	return ResourceConfig(result), nil
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

func generateResource(cfg ResourceConfig, skipEntity, skipRepository bool) error {
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

	data := buildResourceTemplateData(cfg)

	templates := []struct {
		template   string
		subPackage string
		fileName   string
		skip       bool
	}{
		{"resource/layered/Controller.java.tmpl", "controller", cfg.Name + "Controller.java", false},
		{"resource/layered/Service.java.tmpl", "service", cfg.Name + "Service.java", false},
		{"resource/layered/ServiceImpl.java.tmpl", "service/impl", cfg.Name + "ServiceImpl.java", false},
		{"resource/layered/Repository.java.tmpl", "repository", cfg.Name + "Repository.java", skipRepository || !cfg.HasJpa},
		{"resource/layered/Entity.java.tmpl", "entity", cfg.Name + ".java", skipEntity || !cfg.HasJpa},
		{"resource/layered/Request.java.tmpl", "dto", cfg.Name + "Request.java", false},
		{"resource/layered/Response.java.tmpl", "dto", cfg.Name + "Response.java", false},
		{"resource/layered/Mapper.java.tmpl", "mapper", cfg.Name + "Mapper.java", false},
	}

	if cfg.HasJpa && !skipEntity {
		exceptionPath := filepath.Join(basePath, "exception", "ResourceNotFoundException.java")
		if !engine.FileExists(exceptionPath) {
			templates = append(templates, struct {
				template   string
				subPackage string
				fileName   string
				skip       bool
			}{"resource/layered/ResourceNotFoundException.java.tmpl", "exception", "ResourceNotFoundException.java", false})
		}
	}

	log.Info("Generating resource (legacy mode)", "name", cfg.Name)

	if cfg.HasLombok {
		log.Debug("Using Lombok annotations")
	}
	if cfg.HasJpa {
		log.Debug("Generating JPA Entity and Repository")
	}
	if cfg.HasValidation {
		log.Debug("Adding @Valid annotations")
	}

	generatedCount := 0
	skippedCount := 0

	for _, t := range templates {
		if t.skip {
			continue
		}

		outputPath := filepath.Join(basePath, t.subPackage, t.fileName)

		if engine.FileExists(outputPath) {
			log.Warning("File exists, skipping", "file", FormatRelativePath(cwd, outputPath))
			skippedCount++
			continue
		}

		if err := engine.RenderAndWrite(t.template, outputPath, data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", t.fileName, err)
		}

		log.Info("Created", "file", FormatRelativePath(cwd, outputPath))
		generatedCount++
	}

	if generatedCount > 0 {
		log.Success(fmt.Sprintf("Generated %d files for %s resource", generatedCount, cfg.Name))
	}
	if skippedCount > 0 {
		log.Info(fmt.Sprintf("Skipped %d existing files", skippedCount))
	}

	return nil
}

func buildResourceTemplateData(cfg ResourceConfig) map[string]any {
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
