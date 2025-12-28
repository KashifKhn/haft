package generate

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/output"
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

The command auto-detects your project's base package from your build file
and checks for Validation dependency to add @Valid annotations.`,
		Example: `  # Interactive mode
  haft generate controller

  # With controller name
  haft generate controller user
  haft g co product

  # Non-interactive with package override
  haft generate controller order --package com.example.app --no-interactive

  # Output as JSON
  haft generate controller user --json --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runController,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from build file)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("json", false, "Output as JSON")

	return cmd
}

func runController(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	jsonOutput, _ := cmd.Flags().GetBool("json")
	log := logger.Default()

	cfg, err := DetectProjectConfig()
	if err != nil {
		if noInteractive {
			if jsonOutput {
				return output.Error("DETECTION_ERROR", "Could not detect project configuration", err.Error())
			}
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
			if jsonOutput {
				return output.Error("WIZARD_ERROR", err.Error())
			}
			return err
		}
	}

	if err := ValidateComponentConfig(cfg); err != nil {
		if jsonOutput {
			return output.Error("VALIDATION_ERROR", err.Error())
		}
		return err
	}

	tracker := NewGenerateTracker("controller", cfg.Name)

	if !jsonOutput {
		log.Info("Generating controller", "name", cfg.Name)
	}

	generated, err := GenerateComponent(cfg, "resource/layered/Controller.java.tmpl", "controller", "{Name}Controller.java")
	if err != nil {
		if jsonOutput {
			tracker.AddError(err.Error())
			return OutputGenerateResult(true, tracker)
		}
		return err
	}

	if generated {
		tracker.AddGenerated(cfg.Name + "Controller.java")
	} else {
		tracker.AddSkipped(cfg.Name + "Controller.java")
	}

	return OutputGenerateResult(jsonOutput, tracker)
}
