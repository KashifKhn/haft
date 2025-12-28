package template

import (
	"fmt"
	"os"
	"strings"

	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/output"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	sourceStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	projectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	globalStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	pathStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func newListCommand() *cobra.Command {
	var customOnly bool
	var category string
	var showPaths bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available templates",
		Long: `List all available code generation templates and their sources.

Templates can come from three sources:
  - project: .haft/templates/ (highest priority)
  - global: ~/.haft/templates/
  - embedded: built-in templates (fallback)

When a template exists in multiple locations, only the highest
priority source is used during code generation.`,
		Example: `  # List all templates
  haft template list

  # List only custom (overridden) templates
  haft template list --custom

  # List templates in a specific category
  haft template list --category resource

  # Show full paths
  haft template list --paths

  # Output as JSON
  haft template list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(customOnly, category, showPaths, jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&customOnly, "custom", false, "Show only custom (non-embedded) templates")
	cmd.Flags().StringVarP(&category, "category", "c", "", "Filter by category (resource, test, project)")
	cmd.Flags().BoolVar(&showPaths, "paths", false, "Show full template paths")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output result as JSON")

	return cmd
}

func runList(customOnly bool, category string, showPaths bool, jsonOutput bool) error {
	fs := afero.NewOsFs()

	cwd, err := os.Getwd()
	if err != nil {
		if jsonOutput {
			return output.Error("CWD_ERROR", "could not determine current directory", err.Error())
		}
		return fmt.Errorf("could not determine current directory: %w", err)
	}

	loader := generator.NewTemplateLoader(fs, cwd)

	templates, err := loader.ListAllTemplatesWithSource()
	if err != nil {
		if jsonOutput {
			return output.Error("LIST_ERROR", "could not list templates", err.Error())
		}
		return fmt.Errorf("could not list templates: %w", err)
	}

	if category != "" {
		templates = filterByCategory(templates, category)
	}

	if customOnly {
		templates = filterCustomOnly(templates)
	}

	if len(templates) == 0 {
		if jsonOutput {
			return output.Success(output.TemplateListOutput{
				Templates:     []output.TemplateInfo{},
				Total:         0,
				ProjectCount:  0,
				GlobalCount:   0,
				EmbeddedCount: 0,
			})
		}
		fmt.Println("No templates found matching criteria")
		return nil
	}

	projectCount := countBySource(templates, generator.SourceProject)
	globalCount := countBySource(templates, generator.SourceGlobal)
	embeddedCount := countBySource(templates, generator.SourceEmbedded)

	if jsonOutput {
		var outputTemplates []output.TemplateInfo
		for _, t := range templates {
			sourceStr := "embedded"
			switch t.Source {
			case generator.SourceProject:
				sourceStr = "project"
			case generator.SourceGlobal:
				sourceStr = "global"
			}

			parts := strings.Split(t.Name, "/")
			cat := "other"
			if len(parts) > 0 {
				cat = parts[0]
			}

			outputTemplates = append(outputTemplates, output.TemplateInfo{
				Name:     t.Name,
				Category: cat,
				Source:   sourceStr,
				Path:     t.Path,
			})
		}

		return output.Success(output.TemplateListOutput{
			Templates:     outputTemplates,
			Total:         len(templates),
			ProjectCount:  projectCount,
			GlobalCount:   globalCount,
			EmbeddedCount: embeddedCount,
		})
	}

	fmt.Println()
	fmt.Println(headerStyle.Render("  Available Templates"))
	fmt.Println(strings.Repeat("─", 60))

	grouped := groupByCategory(templates)

	titleCaser := cases.Title(language.English)
	for _, cat := range []string{"resource", "test", "project"} {
		if items, ok := grouped[cat]; ok {
			fmt.Println()
			fmt.Printf("  %s\n", headerStyle.Render(titleCaser.String(cat)))

			for _, info := range items {
				printTemplateInfo(info, showPaths)
			}
		}
	}

	for cat, items := range grouped {
		if cat != "resource" && cat != "test" && cat != "project" {
			fmt.Println()
			fmt.Printf("  %s\n", headerStyle.Render(titleCaser.String(cat)))

			for _, info := range items {
				printTemplateInfo(info, showPaths)
			}
		}
	}

	fmt.Println()

	fmt.Printf("  %s %d total", sourceStyle.Render("Templates:"), len(templates))
	if projectCount > 0 {
		fmt.Printf(" (%s %d)", projectStyle.Render("project:"), projectCount)
	}
	if globalCount > 0 {
		fmt.Printf(" (%s %d)", globalStyle.Render("global:"), globalCount)
	}
	if embeddedCount > 0 {
		fmt.Printf(" (%s %d)", sourceStyle.Render("embedded:"), embeddedCount)
	}
	fmt.Println()
	fmt.Println()

	return nil
}

func printTemplateInfo(info generator.TemplateInfo, showPaths bool) {
	name := info.Name

	var sourceLabel string
	switch info.Source {
	case generator.SourceProject:
		sourceLabel = projectStyle.Render("[project]")
	case generator.SourceGlobal:
		sourceLabel = globalStyle.Render("[global]")
	default:
		sourceLabel = sourceStyle.Render("[embedded]")
	}

	if showPaths {
		fmt.Printf("    %s %s\n", name, sourceLabel)
		fmt.Printf("      %s\n", pathStyle.Render(info.Path))
	} else {
		fmt.Printf("    %s %s\n", name, sourceLabel)
	}
}

func filterByCategory(templates []generator.TemplateInfo, category string) []generator.TemplateInfo {
	var filtered []generator.TemplateInfo
	for _, t := range templates {
		if strings.HasPrefix(t.Name, category+"/") {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func filterCustomOnly(templates []generator.TemplateInfo) []generator.TemplateInfo {
	var filtered []generator.TemplateInfo
	for _, t := range templates {
		if t.Source != generator.SourceEmbedded {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func groupByCategory(templates []generator.TemplateInfo) map[string][]generator.TemplateInfo {
	grouped := make(map[string][]generator.TemplateInfo)
	for _, t := range templates {
		parts := strings.Split(t.Name, "/")
		category := "other"
		if len(parts) > 0 {
			category = parts[0]
		}
		grouped[category] = append(grouped[category], t)
	}
	return grouped
}

func countBySource(templates []generator.TemplateInfo, source generator.TemplateSource) int {
	count := 0
	for _, t := range templates {
		if t.Source == source {
			count++
		}
	}
	return count
}
