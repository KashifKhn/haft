package dev

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newCleanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean build artifacts",
		Long: `Clean build artifacts and compiled files.

This command detects your build tool (Maven or Gradle) and runs the 
appropriate command to remove build artifacts.

For Maven:  mvn clean
For Gradle: ./gradlew clean`,
		Example: `  # Clean build artifacts
  haft dev clean`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClean()
		},
	}

	return cmd
}

func runClean() error {
	fs := afero.NewOsFs()
	result, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	logger.Info("Cleaning project", "build-tool", result.BuildTool.DisplayName())

	var cmdArgs []string
	var executable string

	switch result.BuildTool {
	case buildtool.Maven:
		executable = getMavenExecutable()
		cmdArgs = []string{"clean"}

	case buildtool.Gradle, buildtool.GradleKotln:
		executable = getGradleExecutable()
		cmdArgs = []string{"clean"}
	}

	return executeCommand(executable, cmdArgs)
}
