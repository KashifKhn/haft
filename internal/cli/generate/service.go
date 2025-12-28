package generate

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/output"
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

The command auto-detects your project's base package from your build file
and checks for JPA dependency to include repository injection.`,
		Example: `  # Interactive mode
  haft generate service

  # With service name
  haft generate service user
  haft g s product

  # Non-interactive with package override
  haft generate service order --package com.example.app --no-interactive

  # Output as JSON
  haft generate service user --json --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runService,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from build file)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("json", false, "Output as JSON")

	return cmd
}

func runService(cmd *cobra.Command, args []string) error {
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
		cfg, err = RunComponentWizard("Generate Service", cfg, "Service")
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

	tracker := NewGenerateTracker("service", cfg.Name)

	if !jsonOutput {
		log.Info("Generating service", "name", cfg.Name)
	}

	if generated, err := GenerateComponent(cfg, "resource/layered/Service.java.tmpl", "service", "{Name}Service.java"); err != nil {
		tracker.AddError(err.Error())
		if !jsonOutput {
			return err
		}
	} else if generated {
		tracker.AddGenerated(cfg.Name + "Service.java")
	} else {
		tracker.AddSkipped(cfg.Name + "Service.java")
	}

	if generated, err := GenerateComponent(cfg, "resource/layered/ServiceImpl.java.tmpl", "service/impl", "{Name}ServiceImpl.java"); err != nil {
		tracker.AddError(err.Error())
		if !jsonOutput {
			return err
		}
	} else if generated {
		tracker.AddGenerated(cfg.Name + "ServiceImpl.java")
	} else {
		tracker.AddSkipped(cfg.Name + "ServiceImpl.java")
	}

	return OutputGenerateResult(jsonOutput, tracker)
}
