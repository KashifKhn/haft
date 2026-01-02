package dev

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newPackageCommand() *cobra.Command {
	var skipTests bool
	var clean bool
	var profile string

	cmd := &cobra.Command{
		Use:     "package",
		Aliases: []string{"pkg", "jar"},
		Short:   "Create deployable artifact without tests",
		Long: `Create the deployable artifact (JAR/WAR) without running tests.

This command is faster than 'haft dev build' as it skips tests by default.
It's useful when you need to quickly create an artifact for deployment
or testing in another environment.

For Maven:  mvn package -DskipTests
For Gradle: ./gradlew bootJar -x test`,
		Example: `  # Create artifact without tests
  haft dev package

  # Create artifact with tests
  haft dev package --no-skip-tests

  # Clean before packaging
  haft dev package --clean

  # Package with specific profile
  haft dev package --profile prod

  # Clean build with profile
  haft dev package -c -p prod`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPackage(skipTests, clean, profile)
		},
	}

	cmd.Flags().BoolVarP(&skipTests, "skip-tests", "s", true, "Skip running tests (default: true)")
	cmd.Flags().BoolVarP(&clean, "clean", "c", false, "Clean before packaging")
	cmd.Flags().StringVarP(&profile, "profile", "p", "", "Maven/Gradle profile to activate")

	return cmd
}

func runPackage(skipTests bool, clean bool, profile string) error {
	fs := afero.NewOsFs()
	result, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	logger.Info("Creating deployable artifact", "build-tool", result.BuildTool.DisplayName())

	var cmdArgs []string
	var executable string

	switch result.BuildTool {
	case buildtool.Maven:
		executable = getMavenExecutable()
		if clean {
			cmdArgs = append(cmdArgs, "clean")
		}
		cmdArgs = append(cmdArgs, "package")
		if skipTests {
			cmdArgs = append(cmdArgs, "-DskipTests")
		}
		if profile != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("-P%s", profile))
		}

	case buildtool.Gradle, buildtool.GradleKotln:
		executable = getGradleExecutable()
		if clean {
			cmdArgs = append(cmdArgs, "clean")
		}
		cmdArgs = append(cmdArgs, "bootJar")
		if skipTests {
			cmdArgs = append(cmdArgs, "-x", "test")
		}
		if profile != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("-P%s", profile))
		}
	}

	return executeCommand(executable, cmdArgs)
}
