package dev

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newVerifyCommand() *cobra.Command {
	var skipTests bool
	var skipIntegration bool
	var profile string

	cmd := &cobra.Command{
		Use:     "verify",
		Aliases: []string{"vfy"},
		Short:   "Run integration tests and quality checks",
		Long: `Run integration tests and quality checks using the detected build tool.

This command runs the full verification lifecycle including:
- Compile the project
- Run unit tests
- Run integration tests
- Execute quality checks (Checkstyle, SpotBugs, etc.)

For Maven:  mvn verify
For Gradle: ./gradlew check

This is more comprehensive than 'haft dev test' which only runs unit tests.`,
		Example: `  # Run full verification
  haft dev verify

  # Skip all tests (only run quality checks)
  haft dev verify --skip-tests

  # Skip integration tests only
  haft dev verify --skip-integration

  # Run with specific profile
  haft dev verify --profile ci`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVerify(skipTests, skipIntegration, profile)
		},
	}

	cmd.Flags().BoolVarP(&skipTests, "skip-tests", "s", false, "Skip all tests during verification")
	cmd.Flags().BoolVarP(&skipIntegration, "skip-integration", "i", false, "Skip integration tests only")
	cmd.Flags().StringVarP(&profile, "profile", "p", "", "Maven/Gradle profile to activate")

	return cmd
}

func runVerify(skipTests bool, skipIntegration bool, profile string) error {
	fs := afero.NewOsFs()
	result, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	logger.Info("Running verification", "build-tool", result.BuildTool.DisplayName())

	var cmdArgs []string
	var executable string

	switch result.BuildTool {
	case buildtool.Maven:
		executable = getMavenExecutable()
		cmdArgs = []string{"verify"}
		if skipTests {
			cmdArgs = append(cmdArgs, "-DskipTests")
		} else if skipIntegration {
			cmdArgs = append(cmdArgs, "-DskipITs")
		}
		if profile != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("-P%s", profile))
		}

	case buildtool.Gradle, buildtool.GradleKotln:
		executable = getGradleExecutable()
		cmdArgs = []string{"check"}
		if skipTests {
			cmdArgs = append(cmdArgs, "-x", "test")
		} else if skipIntegration {
			cmdArgs = append(cmdArgs, "-x", "integrationTest")
		}
		if profile != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("-P%s", profile))
		}
	}

	return executeCommand(executable, cmdArgs)
}
