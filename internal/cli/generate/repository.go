package generate

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/output"
	"github.com/spf13/cobra"
)

func newRepositoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "repository [name]",
		Aliases: []string{"repo"},
		Short:   "Generate a JPA repository interface",
		Long: `Generate a Spring Data JPA repository interface.

Creates a repository interface extending JpaRepository with
standard CRUD operations. Requires Spring Data JPA dependency.

The command auto-detects your project's base package from your build file.`,
		Example: `  # Interactive mode
  haft generate repository

  # With repository name
  haft generate repository user
  haft g repo product

  # Non-interactive with package override
  haft generate repository order --package com.example.app --no-interactive

  # Output as JSON
  haft generate repository user --json --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runRepository,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from build file)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("json", false, "Output as JSON")

	return cmd
}

func runRepository(cmd *cobra.Command, args []string) error {
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

	if !cfg.HasJpa && !jsonOutput {
		log.Warning("Spring Data JPA not detected in project dependencies")
	}

	if len(args) > 0 {
		cfg.Name = ToPascalCase(args[0])
	}

	if pkg, _ := cmd.Flags().GetString("package"); pkg != "" {
		cfg.BasePackage = pkg
	}

	if !noInteractive {
		cfg, err = RunComponentWizard("Generate Repository", cfg, "Repository")
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

	tracker := NewGenerateTracker("repository", cfg.Name)

	if !jsonOutput {
		log.Info("Generating repository", "name", cfg.Name)
	}

	if generated, err := GenerateComponent(cfg, "resource/layered/Repository.java.tmpl", "repository", "{Name}Repository.java"); err != nil {
		tracker.AddError(err.Error())
		if !jsonOutput {
			return err
		}
	} else if generated {
		tracker.AddGenerated(cfg.Name + "Repository.java")
	} else {
		tracker.AddSkipped(cfg.Name + "Repository.java")
	}

	return OutputGenerateResult(jsonOutput, tracker)
}
