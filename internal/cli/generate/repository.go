package generate

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/logger"
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
  haft generate repository order --package com.example.app --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runRepository,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from build file)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")

	return cmd
}

func runRepository(cmd *cobra.Command, args []string) error {
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
		cfg, err = RunComponentWizard("Generate Repository", cfg, "Repository")
		if err != nil {
			return err
		}
	}

	if err := ValidateComponentConfig(cfg); err != nil {
		return err
	}

	log.Info("Generating repository", "name", cfg.Name)
	_, err = GenerateComponent(cfg, "resource/Repository.java.tmpl", "repository", "{Name}Repository.java")
	return err
}
