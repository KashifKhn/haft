package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	lineStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
)

func newValidateCommand() *cobra.Command {
	var showVars bool
	var showConditions bool

	cmd := &cobra.Command{
		Use:   "validate [template-path]",
		Short: "Validate custom templates for syntax errors",
		Long: `Validate custom templates before using them for code generation.

This command checks for:
  - Unclosed placeholders (missing '}')
  - Unmatched @if/@endif directives
  - Unknown variables (warnings)
  - Go template syntax errors

You can validate a single template file or all templates in a directory.`,
		Example: `  # Validate all project templates
  haft template validate

  # Validate a specific template
  haft template validate .haft/templates/resource/layered/Controller.java.tmpl

  # Validate all templates in a directory
  haft template validate .haft/templates/resource/

  # Show available variables
  haft template validate --vars

  # Show available conditions
  haft template validate --conditions`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVars {
				printAvailableVariables()
				return nil
			}
			if showConditions {
				printAvailableConditions()
				return nil
			}
			return runValidate(args)
		},
	}

	cmd.Flags().BoolVar(&showVars, "vars", false, "Show available template variables")
	cmd.Flags().BoolVar(&showConditions, "conditions", false, "Show available conditions for @if directives")

	return cmd
}

func runValidate(args []string) error {
	log := logger.Default()
	fs := afero.NewOsFs()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not determine current directory: %w", err)
	}

	var templatePaths []string

	if len(args) == 0 {
		loader := generator.NewTemplateLoader(fs, cwd)
		if !loader.ProjectTemplatesExist() {
			log.Info("No custom templates found in .haft/templates/")
			log.Info("Use 'haft template init' to create custom templates")
			return nil
		}

		projectTemplates, err := loader.ListProjectTemplates()
		if err != nil {
			return fmt.Errorf("could not list templates: %w", err)
		}

		templateDir := loader.GetProjectTemplateDir()
		for _, tmpl := range projectTemplates {
			templatePaths = append(templatePaths, filepath.Join(templateDir, tmpl))
		}
	} else {
		for _, arg := range args {
			path := arg
			if !filepath.IsAbs(path) {
				path = filepath.Join(cwd, arg)
			}

			info, err := fs.Stat(path)
			if err != nil {
				return fmt.Errorf("path not found: %s", arg)
			}

			if info.IsDir() {
				err := afero.Walk(fs, path, func(p string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() && strings.HasSuffix(p, ".tmpl") {
						templatePaths = append(templatePaths, p)
					}
					return nil
				})
				if err != nil {
					return fmt.Errorf("could not walk directory: %w", err)
				}
			} else {
				templatePaths = append(templatePaths, path)
			}
		}
	}

	if len(templatePaths) == 0 {
		log.Info("No templates found to validate")
		return nil
	}

	fmt.Println()
	fmt.Printf("  Validating %d template(s)...\n\n", len(templatePaths))

	totalErrors := 0
	totalWarnings := 0
	validCount := 0

	for _, templatePath := range templatePaths {
		content, err := afero.ReadFile(fs, templatePath)
		if err != nil {
			log.Warning("Could not read template", "path", templatePath, "error", err.Error())
			continue
		}

		relPath, _ := filepath.Rel(cwd, templatePath)
		if relPath == "" {
			relPath = templatePath
		}

		result := generator.ValidateTemplate(string(content), filepath.Base(templatePath))

		if result.Valid && len(result.Warnings) == 0 {
			fmt.Printf("  %s %s\n", successStyle.Render("✓"), relPath)
			validCount++
		} else if result.Valid && len(result.Warnings) > 0 {
			fmt.Printf("  %s %s\n", warningStyle.Render("⚠"), relPath)
			for _, warn := range result.Warnings {
				fmt.Printf("      %s line %d: %s\n",
					warningStyle.Render("warning"),
					warn.Line,
					warn.Message)
			}
			totalWarnings += len(result.Warnings)
			validCount++
		} else {
			fmt.Printf("  %s %s\n", errorStyle.Render("✗"), relPath)
			for _, err := range result.Errors {
				fmt.Printf("      %s line %d: %s\n",
					errorStyle.Render("error"),
					err.Line,
					err.Message)
			}
			for _, warn := range result.Warnings {
				fmt.Printf("      %s line %d: %s\n",
					warningStyle.Render("warning"),
					warn.Line,
					warn.Message)
			}
			totalErrors += len(result.Errors)
			totalWarnings += len(result.Warnings)
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("─", 50))

	if totalErrors == 0 {
		fmt.Printf("  %s All %d template(s) are valid\n",
			successStyle.Render("✓"),
			validCount)
	} else {
		fmt.Printf("  %s %d error(s) in %d template(s)\n",
			errorStyle.Render("✗"),
			totalErrors,
			len(templatePaths)-validCount)
	}

	if totalWarnings > 0 {
		fmt.Printf("  %s %d warning(s)\n",
			warningStyle.Render("⚠"),
			totalWarnings)
	}

	fmt.Println()

	if totalErrors > 0 {
		return fmt.Errorf("validation failed with %d error(s)", totalErrors)
	}

	return nil
}

func printAvailableVariables() {
	fmt.Println()
	fmt.Println(infoStyle.Render("  Available Template Variables"))
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()

	vars := generator.GetAvailableVariables()
	for _, v := range vars {
		fmt.Printf("  %-20s %s\n", successStyle.Render(v.Name), v.Description)
		fmt.Printf("  %-20s Example: %s\n", "", lineStyle.Render(v.Example))
		fmt.Println()
	}

	fmt.Println(lineStyle.Render("  Usage: ${VariableName} in your template"))
	fmt.Println()
}

func printAvailableConditions() {
	fmt.Println()
	fmt.Println(infoStyle.Render("  Available Conditions for @if Directives"))
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()

	conditions := generator.GetAvailableConditions()
	for _, c := range conditions {
		fmt.Printf("  %-20s %s\n", successStyle.Render(c.Name), c.Description)
	}

	fmt.Println()
	fmt.Println(lineStyle.Render("  Usage:"))
	fmt.Println(lineStyle.Render("    // @if HasLombok"))
	fmt.Println(lineStyle.Render("    @Data"))
	fmt.Println(lineStyle.Render("    // @endif"))
	fmt.Println()
}
