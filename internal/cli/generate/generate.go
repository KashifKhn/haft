package generate

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Generate Spring Boot components",
		Long: `Generate Spring Boot components like controllers, services, and resources.

The generate command provides sub-commands to scaffold individual components
or complete CRUD resources with all necessary layers.

All generators auto-detect your project configuration from your build file (pom.xml or build.gradle):
  - Base package
  - Lombok dependency (for annotations)
  - Spring Data JPA (for entity/repository)
  - Validation (for @Valid annotations)`,
		Example: `  # Generate a complete CRUD resource (recommended)
  haft generate resource user
  haft g r product

  # Generate individual components
  haft generate controller order
  haft g co payment

  haft generate service user
  haft g s product

  haft generate repository order
  haft g repo payment

  haft generate entity user
  haft g e product

  haft generate dto order

  # Generate exception handler
  haft generate exception
  haft g ex

  # Generate configuration classes
  haft generate config
  haft g cfg

  # Generate security configuration
  haft generate security
  haft g sec --jwt`,
	}

	cmd.AddCommand(newResourceCommand())
	cmd.AddCommand(newControllerCommand())
	cmd.AddCommand(newServiceCommand())
	cmd.AddCommand(newRepositoryCommand())
	cmd.AddCommand(newEntityCommand())
	cmd.AddCommand(newDtoCommand())
	cmd.AddCommand(newExceptionCommand())
	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newSecurityCommand())

	return cmd
}
