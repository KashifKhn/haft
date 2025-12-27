package template

import (
	"fmt"
	"os"

	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newInitCommand() *cobra.Command {
	var category string
	var global bool
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize custom templates in your project",
		Long: `Copy built-in templates to your project's .haft/templates directory.

This allows you to customize the generated code templates for:
  - Company coding standards and copyright headers
  - Custom annotations (@Audited, @Cacheable, etc.)
  - Different patterns (records for DTOs, custom base classes)
  - Extra methods in repositories or services

Templates can be initialized at project level (.haft/templates/) or
global level (~/.haft/templates/). Project templates take priority.`,
		Example: `  # Initialize all templates for customization
  haft template init

  # Initialize only resource templates
  haft template init --category resource

  # Initialize only test templates
  haft template init --category test

  # Initialize templates globally (user-level)
  haft template init --global

  # Overwrite existing templates
  haft template init --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(category, global, force)
		},
	}

	cmd.Flags().StringVarP(&category, "category", "c", "", "Template category to initialize (resource, test, project)")
	cmd.Flags().BoolVarP(&global, "global", "g", false, "Initialize templates in global ~/.haft/templates directory")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing templates")

	return cmd
}

func runInit(category string, global bool, force bool) error {
	log := logger.Default()
	fs := afero.NewOsFs()

	var targetDir string
	var targetLabel string

	if global {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not determine home directory: %w", err)
		}
		targetDir = homeDir
		targetLabel = "global"
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not determine current directory: %w", err)
		}
		targetDir = cwd
		targetLabel = "project"
	}

	loader := generator.NewTemplateLoader(fs, targetDir)

	var templates []string
	var err error

	if category != "" {
		templates, err = generator.ListEmbeddedTemplates(category)
		if err != nil {
			return fmt.Errorf("invalid category '%s': %w", category, err)
		}
		if len(templates) == 0 {
			return fmt.Errorf("no templates found for category '%s'", category)
		}
	} else {
		templates, err = generator.ListAllEmbeddedTemplates()
		if err != nil {
			return fmt.Errorf("could not list templates: %w", err)
		}
	}

	templateDir := loader.GetProjectTemplateDir()
	if global {
		templateDir = loader.GetGlobalTemplateDir()
	}

	log.Info(fmt.Sprintf("Initializing %s templates", targetLabel), "dir", templateDir)

	copiedCount := 0
	skippedCount := 0

	for _, tmpl := range templates {
		destPath := fmt.Sprintf("%s/%s", templateDir, tmpl)

		exists, _ := afero.Exists(fs, destPath)
		if exists && !force {
			log.Debug("Skipping existing template", "template", tmpl)
			skippedCount++
			continue
		}

		if err := loader.CopyEmbeddedToProject(tmpl); err != nil {
			log.Warning("Failed to copy template", "template", tmpl, "error", err.Error())
			continue
		}

		log.Debug("Copied template", "template", tmpl)
		copiedCount++
	}

	if copiedCount > 0 {
		log.Success(fmt.Sprintf("Initialized %d templates", copiedCount))
	}

	if skippedCount > 0 {
		log.Info(fmt.Sprintf("Skipped %d existing templates (use --force to overwrite)", skippedCount))
	}

	if copiedCount == 0 && skippedCount == 0 {
		log.Warning("No templates were initialized")
	}

	return nil
}
