package generate

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Generate Spring Boot components",
		Long: `Generate Spring Boot components like controllers, services, and resources.

The generate command provides sub-commands to scaffold individual components
or complete CRUD resources with all necessary layers.`,
		Example: `  # Generate a complete CRUD resource (recommended)
  haft generate resource user
  haft g resource product

  # Generate individual components (coming soon)
  haft generate controller order
  haft generate service payment`,
	}

	cmd.AddCommand(newResourceCommand())

	return cmd
}
