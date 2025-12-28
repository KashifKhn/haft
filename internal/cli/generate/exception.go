package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/detector"
	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type exceptionConfig struct {
	BasePackage      string
	SelectedOptional []string
}

type optionalException struct {
	Name        string
	FileName    string
	TemplateKey string
	Description string
}

var optionalExceptions = []optionalException{
	{"Conflict (409)", "ConflictException.java", "HasConflict", "Resource already exists"},
	{"MethodNotAllowed (405)", "MethodNotAllowedException.java", "HasMethodNotAllowed", "HTTP method not supported"},
	{"Gone (410)", "GoneException.java", "HasGone", "Resource no longer available"},
	{"UnsupportedMediaType (415)", "UnsupportedMediaTypeException.java", "HasUnsupportedMediaType", "Wrong content type"},
	{"UnprocessableEntity (422)", "UnprocessableEntityException.java", "HasUnprocessableEntity", "Semantic errors in request"},
	{"TooManyRequests (429)", "TooManyRequestsException.java", "HasTooManyRequests", "Rate limiting"},
	{"InternalServerError (500)", "InternalServerErrorException.java", "HasInternalServerError", "Explicit server error handling"},
	{"ServiceUnavailable (503)", "ServiceUnavailableException.java", "HasServiceUnavailable", "Service temporarily down"},
	{"GatewayTimeout (504)", "GatewayTimeoutException.java", "HasGatewayTimeout", "Upstream timeout"},
}

func newExceptionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "exception",
		Aliases: []string{"ex"},
		Short:   "Generate global exception handler",
		Long: `Generate a global exception handler with @ControllerAdvice.

Default exceptions (always generated):
  - GlobalExceptionHandler.java — Central exception handler
  - ErrorResponse.java — Standardized error response DTO
  - ResourceNotFoundException.java — 404 Not Found
  - BadRequestException.java — 400 Bad Request
  - UnauthorizedException.java — 401 Unauthorized
  - ForbiddenException.java — 403 Forbidden

Optional exceptions (select via interactive picker):
  - ConflictException.java — 409 Conflict
  - MethodNotAllowedException.java — 405 Method Not Allowed
  - GoneException.java — 410 Gone
  - UnsupportedMediaTypeException.java — 415 Unsupported Media Type
  - UnprocessableEntityException.java — 422 Unprocessable Entity
  - TooManyRequestsException.java — 429 Too Many Requests
  - InternalServerErrorException.java — 500 Internal Server Error
  - ServiceUnavailableException.java — 503 Service Unavailable
  - GatewayTimeoutException.java — 504 Gateway Timeout

The handler includes built-in support for validation errors when detected.`,
		Example: `  # Generate with interactive picker for optional exceptions
  haft generate exception
  haft g ex

  # Generate with all optional exceptions
  haft generate exception --all

  # Generate only default exceptions (no optional)
  haft generate exception --no-interactive

  # Override base package
  haft generate exception --package com.example.app`,
		RunE: runException,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from project)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard (default exceptions only)")
	cmd.Flags().Bool("all", false, "Include all optional exceptions")
	cmd.Flags().Bool("refresh", false, "Force re-detection of project profile (ignore cache)")

	return cmd
}

func runException(cmd *cobra.Command, args []string) error {
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

	cfg := exceptionConfig{}

	if includeAll {
		for _, ex := range optionalExceptions {
			cfg.SelectedOptional = append(cfg.SelectedOptional, ex.TemplateKey)
		}
	}

	if !noInteractive {
		wizardCfg, err := runExceptionWizard(profile.BasePackage, includeAll)
		if err != nil {
			return err
		}
		if wizardCfg.BasePackage != "" {
			profile.BasePackage = wizardCfg.BasePackage
		}
		if !includeAll {
			cfg.SelectedOptional = wizardCfg.SelectedOptional
		}
	}

	if profile.BasePackage == "" {
		return fmt.Errorf("base package could not be detected. Use --package flag to specify it (e.g., --package com.example.myapp)")
	}

	return generateExceptionHandler(profile, cfg)
}

func enrichProfileFromBuildFile(profile *detector.ProjectProfile) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	fs := afero.NewOsFs()
	result, err := buildtool.Detect(cwd, fs)
	if err != nil {
		return
	}

	project, err := result.Parser.Parse(result.FilePath)
	if err != nil {
		return
	}

	if !profile.HasValidation {
		profile.HasValidation = result.Parser.HasValidation(project)
		if profile.HasValidation && profile.ValidationStyle == detector.ValidationNone {
			profile.ValidationStyle = detector.ValidationJakarta
		}
	}

	if !profile.Lombok.Detected {
		profile.Lombok.Detected = result.Parser.HasLombok(project)
	}

	if profile.BasePackage == "" {
		profile.BasePackage = result.Parser.GetBasePackage(project)
	}
}

func runExceptionWizard(currentPackage string, skipOptionalPicker bool) (exceptionConfig, error) {
	cfg := exceptionConfig{}

	componentCfg := ComponentConfig{
		BasePackage: currentPackage,
		Name:        "Exception",
	}

	result, err := RunComponentWizard("Generate Exception Handler", componentCfg, "Exception")
	if err != nil {
		return cfg, err
	}
	cfg.BasePackage = result.BasePackage

	if skipOptionalPicker {
		return cfg, nil
	}

	cfg.SelectedOptional, err = runOptionalExceptionPicker()
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

type multiSelectWrapper struct {
	model components.MultiSelectModel
}

func (w multiSelectWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w multiSelectWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (w multiSelectWrapper) View() string {
	return w.model.View()
}

func runOptionalExceptionPicker() ([]string, error) {
	items := make([]components.MultiSelectItem, len(optionalExceptions))
	for i, ex := range optionalExceptions {
		items[i] = components.MultiSelectItem{
			Label:    fmt.Sprintf("%s — %s", ex.Name, ex.Description),
			Value:    ex.TemplateKey,
			Selected: false,
		}
	}

	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label:    "Select optional exceptions to include",
		Items:    items,
		Required: false,
	})

	wrapper := multiSelectWrapper{model: model}
	p := tea.NewProgram(wrapper)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(multiSelectWrapper)
	if result.model.GoBack() {
		return nil, fmt.Errorf("wizard cancelled")
	}

	return result.model.Values(), nil
}

func generateExceptionHandler(profile *detector.ProjectProfile, cfg exceptionConfig) error {
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

	exceptionPackage := getExceptionPackage(profile)
	packagePath := strings.ReplaceAll(exceptionPackage, ".", string(os.PathSeparator))
	basePath := filepath.Join(srcPath, packagePath)

	selectedMap := buildSelectedMap(cfg.SelectedOptional)
	data := buildExceptionTemplateData(profile, exceptionPackage, selectedMap)
	templateDir := getExceptionTemplateDir(profile)

	log.Info("Generating exception handler", "package", exceptionPackage)

	templates := []struct {
		template    string
		fileName    string
		conditional string
	}{
		{templateDir + "/GlobalExceptionHandler.java.tmpl", "GlobalExceptionHandler.java", ""},
		{templateDir + "/ErrorResponse.java.tmpl", "ErrorResponse.java", ""},
		{templateDir + "/ResourceNotFoundException.java.tmpl", "ResourceNotFoundException.java", ""},
		{templateDir + "/BadRequestException.java.tmpl", "BadRequestException.java", ""},
		{templateDir + "/UnauthorizedException.java.tmpl", "UnauthorizedException.java", ""},
		{templateDir + "/ForbiddenException.java.tmpl", "ForbiddenException.java", ""},
		{templateDir + "/ConflictException.java.tmpl", "ConflictException.java", "HasConflict"},
		{templateDir + "/MethodNotAllowedException.java.tmpl", "MethodNotAllowedException.java", "HasMethodNotAllowed"},
		{templateDir + "/GoneException.java.tmpl", "GoneException.java", "HasGone"},
		{templateDir + "/UnsupportedMediaTypeException.java.tmpl", "UnsupportedMediaTypeException.java", "HasUnsupportedMediaType"},
		{templateDir + "/UnprocessableEntityException.java.tmpl", "UnprocessableEntityException.java", "HasUnprocessableEntity"},
		{templateDir + "/TooManyRequestsException.java.tmpl", "TooManyRequestsException.java", "HasTooManyRequests"},
		{templateDir + "/InternalServerErrorException.java.tmpl", "InternalServerErrorException.java", "HasInternalServerError"},
		{templateDir + "/ServiceUnavailableException.java.tmpl", "ServiceUnavailableException.java", "HasServiceUnavailable"},
		{templateDir + "/GatewayTimeoutException.java.tmpl", "GatewayTimeoutException.java", "HasGatewayTimeout"},
	}

	generatedCount := 0
	skippedCount := 0

	for _, t := range templates {
		if t.conditional != "" && !selectedMap[t.conditional] {
			continue
		}

		outputPath := filepath.Join(basePath, t.fileName)

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
		log.Success(fmt.Sprintf("Generated %d exception handler files", generatedCount))
	}
	if skippedCount > 0 {
		log.Info(fmt.Sprintf("Skipped %d existing files", skippedCount))
	}

	return nil
}

func buildSelectedMap(selected []string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range selected {
		m[s] = true
	}
	return m
}

func getExceptionPackage(profile *detector.ProjectProfile) string {
	switch profile.Architecture {
	case detector.ArchFeature:
		return profile.BasePackage + ".common.exception"
	case detector.ArchHexagonal:
		return profile.BasePackage + ".infrastructure.exception"
	case detector.ArchClean:
		return profile.BasePackage + ".infrastructure.exception"
	default:
		return profile.BasePackage + ".exception"
	}
}

func getExceptionTemplateDir(_ *detector.ProjectProfile) string {
	return "exception"
}

func buildExceptionTemplateData(profile *detector.ProjectProfile, exceptionPackage string, selectedMap map[string]bool) map[string]any {
	validationImport := "jakarta.validation"
	if profile.ValidationStyle == detector.ValidationJavax {
		validationImport = "javax.validation"
	}

	return map[string]any{
		"BasePackage":             profile.BasePackage,
		"ExceptionPackage":        exceptionPackage,
		"HasLombok":               profile.Lombok.Detected,
		"HasValidation":           profile.HasValidation,
		"ValidationImport":        validationImport,
		"Architecture":            string(profile.Architecture),
		"HasConflict":             selectedMap["HasConflict"],
		"HasMethodNotAllowed":     selectedMap["HasMethodNotAllowed"],
		"HasGone":                 selectedMap["HasGone"],
		"HasUnsupportedMediaType": selectedMap["HasUnsupportedMediaType"],
		"HasUnprocessableEntity":  selectedMap["HasUnprocessableEntity"],
		"HasTooManyRequests":      selectedMap["HasTooManyRequests"],
		"HasInternalServerError":  selectedMap["HasInternalServerError"],
		"HasServiceUnavailable":   selectedMap["HasServiceUnavailable"],
		"HasGatewayTimeout":       selectedMap["HasGatewayTimeout"],
	}
}
