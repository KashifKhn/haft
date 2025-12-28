package generate

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/logger"
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
  haft generate dto user --response-only`,
		Args: cobra.MaximumNArgs(1),
		RunE: runDto,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from build file)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("request-only", false, "Generate only Request DTO")
	cmd.Flags().Bool("response-only", false, "Generate only Response DTO")

	return cmd
}

func runDto(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	requestOnly, _ := cmd.Flags().GetBool("request-only")
	responseOnly, _ := cmd.Flags().GetBool("response-only")
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
		cfg, err = RunComponentWizard("Generate DTO", cfg, "DTO")
		if err != nil {
			return err
		}
	}

	if err := ValidateComponentConfig(cfg); err != nil {
		return err
	}

	log.Info("Generating DTO", "name", cfg.Name)

	generateBoth := !requestOnly && !responseOnly

	if generateBoth || requestOnly {
		if _, err := GenerateComponent(cfg, "resource/layered/Request.java.tmpl", "dto", "{Name}Request.java"); err != nil {
			return err
		}
	}

	if generateBoth || responseOnly {
		if _, err := GenerateComponent(cfg, "resource/layered/Response.java.tmpl", "dto", "{Name}Response.java"); err != nil {
			return err
		}
	}

	return nil
}
