package template

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage code generation templates",
		Long: `Manage code generation templates for customizing generated code.

Haft supports custom templates at two levels:
  1. Project-level: .haft/templates/ (highest priority)
  2. Global user-level: ~/.haft/templates/
  3. Built-in embedded templates (fallback)

When generating code, Haft checks for custom templates in this order.
This allows you to customize generated code for company standards,
custom annotations, or different patterns.

Templates support a simple placeholder syntax:
  ${Name}           - Resource name (PascalCase)
  ${name}           - Resource name (lowercase)
  ${nameCamel}      - Resource name (camelCase)
  ${BasePackage}    - Base package path

And comment-based conditionals:
  // @if HasLombok
  @Data
  // @endif`,
		Example: `  # Initialize custom templates in your project
  haft template init

  # Initialize specific template category
  haft template init --category resource

  # List all available templates and their sources
  haft template list

  # Validate custom templates
  haft template validate

  # Show available template variables
  haft template validate --vars`,
	}

	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newValidateCommand())

	return cmd
}
