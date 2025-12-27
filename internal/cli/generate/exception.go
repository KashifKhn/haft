package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KashifKhn/haft/internal/detector"
	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newExceptionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "exception",
		Aliases: []string{"ex"},
		Short:   "Generate global exception handler",
		Long: `Generate a global exception handler with @ControllerAdvice.

Creates the following files:
  - GlobalExceptionHandler.java — Central exception handler with @ControllerAdvice
  - ErrorResponse.java — Standardized error response DTO
  - ResourceNotFoundException.java — 404 Not Found exception
  - BadRequestException.java — 400 Bad Request exception
  - ConflictException.java — 409 Conflict exception

The handler includes built-in support for:
  - Validation errors (MethodArgumentNotValidException)
  - Resource not found (ResourceNotFoundException)
  - Bad request (BadRequestException)
  - Conflict (ConflictException)
  - Generic exceptions with proper error responses`,
		Example: `  # Generate exception handler (auto-detects architecture)
  haft generate exception
  haft g ex

  # Non-interactive mode
  haft generate exception --no-interactive

  # Override base package
  haft generate exception --package com.example.app`,
		RunE: runException,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from project)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("refresh", false, "Force re-detection of project profile (ignore cache)")

	return cmd
}

func runException(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
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

	if pkg, _ := cmd.Flags().GetString("package"); pkg != "" {
		profile.BasePackage = pkg
	}

	if !noInteractive && profile.BasePackage == "" {
		cfg, err := runExceptionWizard(profile.BasePackage)
		if err != nil {
			return err
		}
		profile.BasePackage = cfg
	}

	if profile.BasePackage == "" {
		return fmt.Errorf("base package is required")
	}

	return generateExceptionHandler(profile)
}

func runExceptionWizard(currentPackage string) (string, error) {
	cfg := ComponentConfig{
		BasePackage: currentPackage,
		Name:        "Exception",
	}

	result, err := RunComponentWizard("Generate Exception Handler", cfg, "Exception")
	if err != nil {
		return "", err
	}

	return result.BasePackage, nil
}

func generateExceptionHandler(profile *detector.ProjectProfile) error {
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

	data := buildExceptionTemplateData(profile, exceptionPackage)
	templateDir := getExceptionTemplateDir(profile)

	log.Info("Generating exception handler", "package", exceptionPackage)

	templates := []struct {
		template string
		fileName string
	}{
		{templateDir + "/GlobalExceptionHandler.java.tmpl", "GlobalExceptionHandler.java"},
		{templateDir + "/ErrorResponse.java.tmpl", "ErrorResponse.java"},
		{templateDir + "/ResourceNotFoundException.java.tmpl", "ResourceNotFoundException.java"},
		{templateDir + "/BadRequestException.java.tmpl", "BadRequestException.java"},
		{templateDir + "/ConflictException.java.tmpl", "ConflictException.java"},
	}

	generatedCount := 0
	skippedCount := 0

	for _, t := range templates {
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

func buildExceptionTemplateData(profile *detector.ProjectProfile, exceptionPackage string) map[string]any {
	validationImport := "jakarta.validation"
	if profile.ValidationStyle == detector.ValidationJavax {
		validationImport = "javax.validation"
	}

	return map[string]any{
		"BasePackage":      profile.BasePackage,
		"ExceptionPackage": exceptionPackage,
		"HasLombok":        profile.Lombok.Detected,
		"HasValidation":    profile.HasValidation,
		"ValidationImport": validationImport,
		"Architecture":     string(profile.Architecture),
	}
}
