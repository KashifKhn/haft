package dev

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newTestCommand() *cobra.Command {
	var filter string
	var verbose bool
	var failFast bool

	cmd := &cobra.Command{
		Use:     "test",
		Aliases: []string{"t"},
		Short:   "Run tests",
		Long: `Run project tests using the detected build tool.

This command detects your build tool (Maven or Gradle) and runs the 
appropriate command to execute your test suite.

For Maven:  mvn test
For Gradle: ./gradlew test`,
		Example: `  # Run all tests
  haft dev test

  # Run tests matching a pattern
  haft dev test --filter UserService

  # Run with verbose output
  haft dev test --verbose

  # Stop on first failure
  haft dev test --fail-fast`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTest(filter, verbose, failFast)
		},
	}

	cmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter tests by class or method name")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose test output")
	cmd.Flags().BoolVar(&failFast, "fail-fast", false, "Stop on first test failure")

	return cmd
}

func runTest(filter string, verbose bool, failFast bool) error {
	fs := afero.NewOsFs()
	result, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	logger.Info("Running tests", "build-tool", result.BuildTool.DisplayName())

	var cmdArgs []string
	var executable string

	switch result.BuildTool {
	case buildtool.Maven:
		executable = getMavenExecutable()
		cmdArgs = []string{"test"}
		if filter != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("-Dtest=%s", filter))
		}
		if failFast {
			cmdArgs = append(cmdArgs, "-Dsurefire.failIfNoSpecifiedTests=false", "-DfailIfNoTests=false", "--fail-at-end")
		}

	case buildtool.Gradle, buildtool.GradleKotln:
		executable = getGradleExecutable()
		cmdArgs = []string{"test"}
		if filter != "" {
			cmdArgs = append(cmdArgs, "--tests", fmt.Sprintf("*%s*", filter))
		}
		if verbose {
			cmdArgs = append(cmdArgs, "--info")
		}
		if failFast {
			cmdArgs = append(cmdArgs, "--fail-fast")
		}
	}

	return executeCommand(executable, cmdArgs)
}
