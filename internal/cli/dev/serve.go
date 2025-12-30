package dev

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newServeCommand() *cobra.Command {
	var profile string
	var debug bool
	var port int
	var noInteractive bool

	cmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"run", "start"},
		Short:   "Start the application with hot-reload",
		Long: `Start the Spring Boot application in supervisor mode.

This command detects your build tool (Maven or Gradle) and runs the 
appropriate command to start your application.

In interactive mode (default when running in a terminal), you can:
  - Press 'r' to restart (compiles first, keeps server if compile fails)
  - Press 'q' to quit
  - Press 'c' to clear screen
  - Press 'h' for help

For Maven:  mvn spring-boot:run -DskipTests
For Gradle: ./gradlew bootRun -x test

External plugins (Neovim, VSCode, IntelliJ) can trigger restart by
creating a trigger file at .haft/restart`,
		Example: `  # Start with default settings (interactive mode)
  haft dev serve

  # Start with specific profile
  haft dev serve --profile dev

  # Start with debug mode
  haft dev serve --debug

  # Start on specific port
  haft dev serve --port 8081

  # Start in non-interactive mode (for CI/scripts)
  haft dev serve --no-interactive`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(profile, debug, port, noInteractive)
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", "", "Spring profile to activate (e.g., dev, prod)")
	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable remote debugging on port 5005")
	cmd.Flags().IntVar(&port, "port", 0, "Server port (overrides application config)")
	cmd.Flags().BoolVar(&noInteractive, "no-interactive", false, "Disable interactive mode (no keyboard commands)")

	return cmd
}

func runServe(profile string, debug bool, port int, noInteractive bool) error {
	fs := afero.NewOsFs()
	result, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	keyboard := NewKeyboardListener()
	isInteractive := keyboard.IsInteractive() && !noInteractive

	var stdout, stderr io.Writer = os.Stdout, os.Stderr
	if isInteractive {
		stdout = newCRLFWriter(os.Stdout)
		stderr = newCRLFWriter(os.Stderr)
	}

	pm := NewProcessManager(ProcessConfig{
		BuildTool: result.BuildTool,
		Profile:   profile,
		Debug:     debug,
		Port:      port,
		Stdout:    stdout,
		Stderr:    stderr,
	})

	trigger := NewTriggerWatcher()
	if err := trigger.Setup(); err != nil {
		logger.Warning("Failed to setup trigger watcher", "error", err)
	}
	defer trigger.Cleanup()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	if isInteractive {
		PrintBanner()
		logger.Info("Starting application", "build-tool", result.BuildTool.DisplayName())

		if err := keyboard.Start(); err != nil {
			logger.Warning("Failed to enable interactive mode", "error", err)
			isInteractive = false
		} else {
			defer func() {
				if err := keyboard.Stop(); err != nil {
					logger.Debug("Failed to restore terminal state", "error", err)
				}
			}()
		}
	} else {
		logger.Info("Starting application (non-interactive)", "build-tool", result.BuildTool.DisplayName())
	}

	if err := pm.Start(); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case sig := <-sigChan:
			logger.Info("Received signal, shutting down", "signal", sig)
			if err := pm.Stop(); err != nil {
				logger.Debug("Error during shutdown", "error", err)
			}
			return nil

		case cmd := <-keyboard.Commands():
			switch cmd {
			case KeyRestart:
				if pm.IsBusy() {
					fmt.Print("\r\n\033[33m→ Restart already in progress...\033[0m\r\n")
					continue
				}
				if err := pm.Restart(ctx); err != nil {
					logger.Debug("Restart failed", "error", err)
				}

			case KeyQuit:
				fmt.Print("\r\n\033[33m→ Shutting down...\033[0m\r\n")
				if err := pm.Stop(); err != nil {
					logger.Debug("Error during shutdown", "error", err)
				}
				return nil

			case KeyClear:
				ClearScreen()
				PrintBanner()

			case KeyHelp:
				PrintKeyCommands()
			}

		case <-trigger.Events():
			if pm.IsBusy() {
				continue
			}
			fmt.Print("\r\n\033[35m→ External restart triggered\033[0m\r\n")
			if err := pm.Restart(ctx); err != nil {
				logger.Debug("External restart failed", "error", err)
			}

		default:
			if pm.State() == StateIdle || pm.State() == StateFailed {
				if lastErr := pm.LastError(); lastErr != nil {
					if isInteractive {
						fmt.Printf("\r\n\033[31mProcess stopped: %v\033[0m\r\n", lastErr)
						fmt.Print("\033[33mPress 'r' to restart or 'q' to quit\033[0m\r\n")
						waitForRestartOrQuit(keyboard, sigChan, pm, ctx)
					} else {
						return lastErr
					}
				}
			}

			select {
			case <-ctx.Done():
				return nil
			default:
			}
		}
	}
}

func waitForRestartOrQuit(keyboard *KeyboardListener, sigChan chan os.Signal, pm *ProcessManager, ctx context.Context) {
	for {
		select {
		case sig := <-sigChan:
			logger.Info("Received signal", "signal", sig)
			if err := pm.Stop(); err != nil {
				logger.Debug("Error during shutdown", "error", err)
			}
			return

		case cmd := <-keyboard.Commands():
			switch cmd {
			case KeyRestart:
				if err := pm.Restart(ctx); err != nil {
					logger.Debug("Restart failed", "error", err)
				}
				return
			case KeyQuit:
				fmt.Print("\r\n\033[33m→ Shutting down...\033[0m\r\n")
				if err := pm.Stop(); err != nil {
					logger.Debug("Error during shutdown", "error", err)
				}
				return
			case KeyHelp:
				PrintKeyCommands()
			case KeyClear:
				ClearScreen()
				PrintBanner()
			}
		}
	}
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
