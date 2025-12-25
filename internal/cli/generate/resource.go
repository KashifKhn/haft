package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/logger"
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
  - ResourceNotFoundException (if JPA and not exists)

The command auto-detects your project's base package from pom.xml and
checks for Lombok, JPA, and Validation dependencies to customize the
generated code accordingly. Dependencies not in your project are 
automatically disabled.`,
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

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from pom.xml)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("skip-entity", false, "Skip entity generation")
	cmd.Flags().Bool("skip-repository", false, "Skip repository generation")

	return cmd
}

func runResource(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	log := logger.Default()

	compCfg, err := DetectProjectConfig()
	if err != nil {
		if noInteractive {
			return fmt.Errorf("could not detect project configuration: %w", err)
		}
		log.Warning("Could not detect project config, using defaults")
	}

	cfg := ResourceConfig{
		Name:          compCfg.Name,
		BasePackage:   compCfg.BasePackage,
		HasLombok:     compCfg.HasLombok,
		HasJpa:        compCfg.HasJpa,
		HasValidation: compCfg.HasValidation,
	}

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
	compCfg := ComponentConfig{
		Name:          cfg.Name,
		BasePackage:   cfg.BasePackage,
		HasLombok:     cfg.HasLombok,
		HasJpa:        cfg.HasJpa,
		HasValidation: cfg.HasValidation,
	}

	result, err := RunComponentWizard("Generate Resource", compCfg, "Resource")
	if err != nil {
		return cfg, err
	}

	return ResourceConfig{
		Name:          result.Name,
		BasePackage:   result.BasePackage,
		HasLombok:     result.HasLombok,
		HasJpa:        result.HasJpa,
		HasValidation: result.HasValidation,
	}, nil
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
