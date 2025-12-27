package dev

import (
	"fmt"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newBuildCommand() *cobra.Command {
	var skipTests bool
	var profile string
	var clean bool

	cmd := &cobra.Command{
		Use:     "build",
		Aliases: []string{"b", "compile"},
		Short:   "Build the project",
		Long: `Build the Spring Boot project using the detected build tool.

This command detects your build tool (Maven or Gradle) and runs the 
appropriate command to compile and package your application.

For Maven:  mvn package
For Gradle: ./gradlew build`,
		Example: `  # Build the project
  haft dev build

  # Build without running tests
  haft dev build --skip-tests

  # Clean and build
  haft dev build --clean

  # Build with specific profile
  haft dev build --profile prod`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBuild(skipTests, profile, clean)
		},
	}

	cmd.Flags().BoolVarP(&skipTests, "skip-tests", "s", false, "Skip running tests during build")
	cmd.Flags().StringVarP(&profile, "profile", "p", "", "Maven/Gradle profile to activate")
	cmd.Flags().BoolVarP(&clean, "clean", "c", false, "Clean before building")

	return cmd
}

func runBuild(skipTests bool, profile string, clean bool) error {
	fs := afero.NewOsFs()
	result, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	logger.Info("Building project", "build-tool", result.BuildTool.DisplayName())

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
		cmdArgs = append(cmdArgs, "build")
		if skipTests {
			cmdArgs = append(cmdArgs, "-x", "test")
		}
		if profile != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("-P%s", profile))
		}
	}

	return executeCommand(executable, cmdArgs)
}
