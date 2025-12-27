package dev

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newServeCommand() *cobra.Command {
	var profile string
	var debug bool
	var port int

	cmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"run", "start"},
		Short:   "Start the application with hot-reload",
		Long: `Start the Spring Boot application with DevTools for hot-reload.

This command detects your build tool (Maven or Gradle) and runs the 
appropriate command to start your application with Spring DevTools enabled.

For Maven:  mvn spring-boot:run
For Gradle: ./gradlew bootRun`,
		Example: `  # Start with default settings
  haft dev serve

  # Start with specific profile
  haft dev serve --profile dev

  # Start with debug mode
  haft dev serve --debug

  # Start on specific port
  haft dev serve --port 8081`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(profile, debug, port)
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", "", "Spring profile to activate (e.g., dev, prod)")
	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable remote debugging on port 5005")
	cmd.Flags().IntVar(&port, "port", 0, "Server port (overrides application config)")

	return cmd
}

func runServe(profile string, debug bool, port int) error {
	fs := afero.NewOsFs()
	result, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	logger.Info("Starting application", "build-tool", result.BuildTool.DisplayName())

	var cmdArgs []string
	var executable string

	switch result.BuildTool {
	case buildtool.Maven:
		executable = getMavenExecutable()
		cmdArgs = []string{"spring-boot:run"}
		if profile != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("-Dspring-boot.run.profiles=%s", profile))
		}
		if debug {
			cmdArgs = append(cmdArgs, "-Dspring-boot.run.jvmArguments=-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005")
		}
		if port > 0 {
			cmdArgs = append(cmdArgs, fmt.Sprintf("-Dspring-boot.run.arguments=--server.port=%d", port))
		}

	case buildtool.Gradle, buildtool.GradleKotln:
		executable = getGradleExecutable()
		cmdArgs = []string{"bootRun"}
		if profile != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("--args=--spring.profiles.active=%s", profile))
		}
		if debug {
			cmdArgs = append(cmdArgs, "--debug-jvm")
		}
		if port > 0 {
			if profile != "" {
				cmdArgs[len(cmdArgs)-1] = fmt.Sprintf("--args=--spring.profiles.active=%s --server.port=%d", profile, port)
			} else {
				cmdArgs = append(cmdArgs, fmt.Sprintf("--args=--server.port=%d", port))
			}
		}
	}

	return executeCommand(executable, cmdArgs)
}

func getMavenExecutable() string {
	if runtime.GOOS == "windows" {
		if _, err := os.Stat("mvnw.cmd"); err == nil {
			return "mvnw.cmd"
		}
	} else {
		if _, err := os.Stat("mvnw"); err == nil {
			return "./mvnw"
		}
	}
	return "mvn"
}

func getGradleExecutable() string {
	if runtime.GOOS == "windows" {
		if _, err := os.Stat("gradlew.bat"); err == nil {
			return "gradlew.bat"
		}
	} else {
		if _, err := os.Stat("gradlew"); err == nil {
			return "./gradlew"
		}
	}
	return "gradle"
}

func executeCommand(executable string, args []string) error {
	cmd := exec.Command(executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
