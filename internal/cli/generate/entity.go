package generate

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/output"
	"github.com/spf13/cobra"
)

func newEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "entity [name]",
		Aliases: []string{"e"},
		Short:   "Generate a JPA entity class",
		Long: `Generate a JPA entity class with standard annotations.

Creates an entity class with:
  - @Entity and @Table annotations
  - Auto-generated ID field with @Id and @GeneratedValue
  - Lombok annotations (if detected in project)

Requires Spring Data JPA dependency.

The command auto-detects your project's base package from your build file
and checks for Lombok dependency to add annotations.`,
		Example: `  # Interactive mode
  haft generate entity

  # With entity name
  haft generate entity user
  haft g e product

  # Non-interactive with package override
  haft generate entity order --package com.example.app --no-interactive

  # Output as JSON
  haft generate entity user --json --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runEntity,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from build file)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("json", false, "Output as JSON")

	return cmd
}

func runEntity(cmd *cobra.Command, args []string) error {
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
		cfg, err = RunComponentWizard("Generate Entity", cfg, "Entity")
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

	tracker := NewGenerateTracker("entity", cfg.Name)

	if !jsonOutput {
		log.Info("Generating entity", "name", cfg.Name)
	}

	if generated, err := GenerateComponent(cfg, "resource/layered/Entity.java.tmpl", "entity", "{Name}.java"); err != nil {
		tracker.AddError(err.Error())
		if !jsonOutput {
			return err
		}
	} else if generated {
		tracker.AddGenerated(cfg.Name + ".java")
	} else {
		tracker.AddSkipped(cfg.Name + ".java")
	}

	return OutputGenerateResult(jsonOutput, tracker)
}
