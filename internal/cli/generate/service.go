package generate

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/cobra"
)

func newServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service [name]",
		Aliases: []string{"s"},
		Short:   "Generate a service interface and implementation",
		Long: `Generate a service layer with interface and implementation.

Creates two files:
  - {Name}Service.java - Service interface with CRUD method signatures
  - {Name}ServiceImpl.java - Service implementation

The command auto-detects your project's base package from pom.xml
and checks for JPA dependency to include repository injection.`,
		Example: `  # Interactive mode
  haft generate service

  # With service name
  haft generate service user
  haft g s product

  # Non-interactive with package override
  haft generate service order --package com.example.app --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runService,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from pom.xml)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")

	return cmd
}

func runService(cmd *cobra.Command, args []string) error {
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
		cfg, err = RunComponentWizard("Generate Service", cfg, "Service")
		if err != nil {
			return err
		}
	}

	if err := ValidateComponentConfig(cfg); err != nil {
		return err
	}

	log.Info("Generating service", "name", cfg.Name)

	if err := GenerateComponent(cfg, "resource/Service.java.tmpl", "service", "{Name}Service.java"); err != nil {
		return err
	}

	return GenerateComponent(cfg, "resource/ServiceImpl.java.tmpl", "service/impl", "{Name}ServiceImpl.java")
}
