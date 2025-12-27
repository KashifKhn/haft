package root

import (
	"os"

	addcmd "github.com/KashifKhn/haft/internal/cli/add"
	completioncmd "github.com/KashifKhn/haft/internal/cli/completion"
	devcmd "github.com/KashifKhn/haft/internal/cli/dev"
	generatecmd "github.com/KashifKhn/haft/internal/cli/generate"
	initcmd "github.com/KashifKhn/haft/internal/cli/init"
	removecmd "github.com/KashifKhn/haft/internal/cli/remove"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
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
  haft add validation

  # Remove dependencies
  haft remove h2
  haft remove lombok validation

  # Development workflow
  haft dev serve          # Start with hot-reload
  haft dev build          # Build project
  haft dev test           # Run tests`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger()
	},
}

func initLogger() {
	logger.SetDefault(logger.New(logger.Options{
		NoColor: noColor,
		Verbose: verbose,
		Output:  os.Stderr,
	}))
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	rootCmd.Version = version
	rootCmd.SetVersionTemplate("haft version {{.Version}}\n")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initcmd.NewCommand())
	rootCmd.AddCommand(generatecmd.NewCommand())
	rootCmd.AddCommand(addcmd.NewCommand())
	rootCmd.AddCommand(removecmd.NewCommand())
	rootCmd.AddCommand(completioncmd.NewCommand())
	rootCmd.AddCommand(devcmd.NewCommand())
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Haft",
	Long:  `Display the current version of Haft CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("haft version", version)
	},
}

func GetVersion() string {
	return version
}

func SetVersion(v string) {
	version = v
	rootCmd.Version = v
}

func IsVerbose() bool {
	return verbose
}

func IsNoColor() bool {
	return noColor
}
