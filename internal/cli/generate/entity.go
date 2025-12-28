package generate

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/logger"
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
  haft generate entity order --package com.example.app --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runEntity,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from build file)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")

	return cmd
}

func runEntity(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	log := logger.Default()

	cfg, err := DetectProjectConfig()
	if err != nil {
		if noInteractive {
			return fmt.Errorf("could not detect project configuration: %w", err)
		}
		log.Warning("Could not detect project config, using defaults")
	}

	if !cfg.HasJpa {
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
			return err
		}
	}

	if err := ValidateComponentConfig(cfg); err != nil {
		return err
	}

	log.Info("Generating entity", "name", cfg.Name)
	_, err = GenerateComponent(cfg, "resource/layered/Entity.java.tmpl", "entity", "{Name}.java")
	return err
}
