package generate

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/output"
	"github.com/spf13/cobra"
)

func newDtoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dto [name]",
		Short: "Generate DTO classes (Request and Response)",
		Long: `Generate Data Transfer Object classes.

Creates two DTO classes:
  - {Name}Request.java - For incoming request data
  - {Name}Response.java - For outgoing response data

The command auto-detects your project's base package from your build file
and checks for Lombok and Validation dependencies to add annotations.`,
		Example: `  # Interactive mode
  haft generate dto

  # With DTO name
  haft generate dto user
  haft g dto product

  # Non-interactive with package override
  haft generate dto order --package com.example.app --no-interactive

  # Generate only request or response
  haft generate dto user --request-only
  haft generate dto user --response-only

  # Output as JSON
  haft generate dto user --json --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runDto,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from build file)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("request-only", false, "Generate only Request DTO")
	cmd.Flags().Bool("response-only", false, "Generate only Response DTO")
	cmd.Flags().Bool("json", false, "Output as JSON")

	return cmd
}

func runDto(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	requestOnly, _ := cmd.Flags().GetBool("request-only")
	responseOnly, _ := cmd.Flags().GetBool("response-only")
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
		cfg, err = RunComponentWizard("Generate DTO", cfg, "DTO")
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

	tracker := NewGenerateTracker("dto", cfg.Name)

	if !jsonOutput {
		log.Info("Generating DTO", "name", cfg.Name)
	}

	generateBoth := !requestOnly && !responseOnly

	if generateBoth || requestOnly {
		if generated, err := GenerateComponent(cfg, "resource/layered/Request.java.tmpl", "dto", "{Name}Request.java"); err != nil {
			tracker.AddError(err.Error())
			if !jsonOutput {
				return err
			}
		} else if generated {
			tracker.AddGenerated(cfg.Name + "Request.java")
		} else {
			tracker.AddSkipped(cfg.Name + "Request.java")
		}
	}

	if generateBoth || responseOnly {
		if generated, err := GenerateComponent(cfg, "resource/layered/Response.java.tmpl", "dto", "{Name}Response.java"); err != nil {
			tracker.AddError(err.Error())
			if !jsonOutput {
				return err
			}
		} else if generated {
			tracker.AddGenerated(cfg.Name + "Response.java")
		} else {
			tracker.AddSkipped(cfg.Name + "Response.java")
		}
	}

	return OutputGenerateResult(jsonOutput, tracker)
}
