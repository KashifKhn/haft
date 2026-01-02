package dev

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newOutdatedCommand() *cobra.Command {
	var includePlugins bool
	var allowSnapshots bool

	cmd := &cobra.Command{
		Use:     "outdated",
		Aliases: []string{"updates", "out"},
		Short:   "Check for dependency updates",
		Long: `Check for newer versions of your project dependencies.

This command scans your dependencies and reports which ones have newer 
versions available. It helps keep your project up-to-date and secure.

For Maven:  mvn versions:display-dependency-updates
For Gradle: ./gradlew dependencyUpdates (requires plugin)

Note: For Gradle projects, this requires the 'com.github.ben-manes.versions' 
plugin to be configured in your build file.`,
		Example: `  # Check for outdated dependencies
  haft dev outdated

  # Include plugin updates (Maven only)
  haft dev outdated --plugins

  # Allow snapshot versions in results
  haft dev outdated --snapshots`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOutdated(includePlugins, allowSnapshots)
		},
	}

	cmd.Flags().BoolVarP(&includePlugins, "plugins", "p", false, "Include plugin updates (Maven only)")
	cmd.Flags().BoolVarP(&allowSnapshots, "snapshots", "s", false, "Include snapshot versions in results")

	return cmd
}

func runOutdated(includePlugins bool, allowSnapshots bool) error {
	fs := afero.NewOsFs()
	result, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	logger.Info("Checking for dependency updates", "build-tool", result.BuildTool.DisplayName())

	var cmdArgs []string
	var executable string

	switch result.BuildTool {
	case buildtool.Maven:
		executable = getMavenExecutable()
		if includePlugins {
			cmdArgs = []string{"versions:display-dependency-updates", "versions:display-plugin-updates"}
		} else {
			cmdArgs = []string{"versions:display-dependency-updates"}
		}
		if !allowSnapshots {
			cmdArgs = append(cmdArgs, "-DallowSnapshots=false")
		}

	case buildtool.Gradle, buildtool.GradleKotln:
		executable = getGradleExecutable()
		cmdArgs = []string{"dependencyUpdates"}
		if !allowSnapshots {
			cmdArgs = append(cmdArgs, "-Drevision=release")
		}
		logger.Warning("Gradle requires the 'com.github.ben-manes.versions' plugin for dependency updates")
	}

	return executeCommand(executable, cmdArgs)
}
