package root

import (
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0-dev"
	verbose bool
	noColor bool
)

var rootCmd = &cobra.Command{
	Use:   "haft",
	Short: "Spring Boot productivity CLI - the developer's grip on Spring Boot projects",
	Long: `Haft is a high-performance CLI tool for the Java Spring Boot ecosystem.

It automates project scaffolding, boilerplate code generation, and 
architectural consistency. Unlike Spring Initializr which only bootstraps 
a project, Haft serves as a lifecycle companion.

Features:
  - Interactive project initialization with TUI wizard
  - Resource generation (Controller, Service, Repository, Entity, DTO)
  - Smart dependency management for Maven/Gradle
  - Architectural enforcement (Layered, Hexagonal)`,
	Example: `  # Initialize a new Spring Boot project
  haft init

  # Generate a complete CRUD resource
  haft generate resource user

  # Generate individual components
  haft generate controller product
  haft generate service order

  # Add dependencies
  haft add lombok
  haft add validation`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Haft",
	Long:  `Display the current version of Haft CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("haft version", version)
	},
}
