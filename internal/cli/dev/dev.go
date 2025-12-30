package dev

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dev",
		Aliases: []string{"d"},
		Short:   "Development commands for running, building, and testing",
		Long: `Development commands that wrap Maven/Gradle for a unified experience.

These commands automatically detect your build tool (Maven or Gradle) and 
execute the appropriate underlying commands.`,
		Example: `  # Start the application with hot-reload
  haft dev serve

  # Build the project
  haft dev build

  # Run tests
  haft dev test

  # Clean build artifacts
  haft dev clean

  # Trigger restart (for use with haft dev serve)
  haft dev restart`,
	}

	cmd.AddCommand(newServeCommand())
	cmd.AddCommand(newBuildCommand())
	cmd.AddCommand(newTestCommand())
	cmd.AddCommand(newCleanCommand())
	cmd.AddCommand(newRestartCommand())

	return cmd
}
