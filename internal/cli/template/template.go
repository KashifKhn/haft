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
custom annotations, or different patterns.`,
		Example: `  # Initialize custom templates in your project
  haft template init

  # Initialize specific template category
  haft template init --category resource

  # List all available templates and their sources
  haft template list

  # List only custom (overridden) templates
  haft template list --custom`,
	}

	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newListCommand())

	return cmd
}
