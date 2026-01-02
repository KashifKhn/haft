package dev

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newDepsCommand() *cobra.Command {
	var configuration string
	var verbose bool

	cmd := &cobra.Command{
		Use:     "deps",
		Aliases: []string{"dependencies", "tree"},
		Short:   "Display project dependency tree",
		Long: `Display the project's dependency tree showing all direct and transitive dependencies.

This command helps you understand your project's dependency graph and 
debug dependency conflicts or version issues.

For Maven:  mvn dependency:tree
For Gradle: ./gradlew dependencies`,
		Example: `  # Show full dependency tree
  haft dev deps

  # Show dependencies for specific configuration (Gradle)
  haft dev deps --configuration compileClasspath

  # Show dependencies for specific scope (Maven)
  haft dev deps --configuration compile

  # Show verbose output
  haft dev deps --verbose`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeps(configuration, verbose)
		},
	}

	cmd.Flags().StringVarP(&configuration, "configuration", "c", "", "Configuration/scope to show (e.g., compile, runtime, test)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose dependency information")

	return cmd
}

func runDeps(configuration string, verbose bool) error {
	fs := afero.NewOsFs()
	result, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	logger.Info("Showing dependency tree", "build-tool", result.BuildTool.DisplayName())

	var cmdArgs []string
	var executable string

	switch result.BuildTool {
	case buildtool.Maven:
		executable = getMavenExecutable()
		cmdArgs = []string{"dependency:tree"}
		if configuration != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("-Dscope=%s", configuration))
		}
		if verbose {
			cmdArgs = append(cmdArgs, "-Dverbose")
		}

	case buildtool.Gradle, buildtool.GradleKotln:
		executable = getGradleExecutable()
		cmdArgs = []string{"dependencies"}
		if configuration != "" {
			cmdArgs = append(cmdArgs, "--configuration", configuration)
		}
		if verbose {
			cmdArgs = append(cmdArgs, "--info")
		}
	}

	return executeCommand(executable, cmdArgs)
}
