package generate

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/cobra"
)

func newControllerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "controller [name]",
		Aliases: []string{"co"},
		Short:   "Generate a REST controller",
		Long: `Generate a REST controller with CRUD endpoints.

Creates a controller class with:
  - GET /api/{resource}s - List all
  - GET /api/{resource}s/{id} - Get by ID
  - POST /api/{resource}s - Create
  - PUT /api/{resource}s/{id} - Update
  - DELETE /api/{resource}s/{id} - Delete

The command auto-detects your project's base package from pom.xml
and checks for Validation dependency to add @Valid annotations.`,
		Example: `  # Interactive mode
  haft generate controller

  # With controller name
  haft generate controller user
  haft g co product

  # Non-interactive with package override
  haft generate controller order --package com.example.app --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runController,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from pom.xml)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")

	return cmd
}

func runController(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	log := logger.Default()

	cfg, err := DetectProjectConfig()
	if err != nil {
		if noInteractive {
			return fmt.Errorf("could not detect project configuration: %w", err)
		}
		log.Warning("Could not detect project config, using defaults")
	}

	if len(args) > 0 {
		cfg.Name = ToPascalCase(args[0])
	}

	if pkg, _ := cmd.Flags().GetString("package"); pkg != "" {
		cfg.BasePackage = pkg
	}

	if !noInteractive {
		cfg, err = RunComponentWizard("Generate Controller", cfg, "Controller")
		if err != nil {
			return err
		}
	}

	if err := ValidateComponentConfig(cfg); err != nil {
		return err
	}

	log.Info("Generating controller", "name", cfg.Name)
	_, err = GenerateComponent(cfg, "resource/Controller.java.tmpl", "controller", "{Name}Controller.java")
	return err
}
