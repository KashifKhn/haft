package dev

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
)

type ProcessState int

const (
	StateIdle ProcessState = iota
	StateStarting
	StateRunning
	StateStopping
	StateRestarting
	StateCompiling
	StateFailed
)

func (s ProcessState) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateStarting:
		return "starting"
	case StateRunning:
		return "running"
	case StateStopping:
		return "stopping"
	case StateRestarting:
		return "restarting"
	case StateCompiling:
		return "compiling"
	case StateFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type ProcessManager struct {
	mu            sync.Mutex
	cmd           *exec.Cmd
	state         ProcessState
	buildTool     buildtool.Type
	profile       string
	debug         bool
	port          int
	lastError     error
	restartCount  int
	stdout        io.Writer
	stderr        io.Writer
	onStateChange func(ProcessState)
}

type ProcessConfig struct {
	BuildTool     buildtool.Type
	Profile       string
	Debug         bool
	Port          int
	Stdout        io.Writer
	Stderr        io.Writer
	OnStateChange func(ProcessState)
}

func NewProcessManager(cfg ProcessConfig) *ProcessManager {
	stdout := cfg.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := cfg.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}
	return &ProcessManager{
		state:         StateIdle,
		buildTool:     cfg.BuildTool,
		profile:       cfg.Profile,
		debug:         cfg.Debug,
		port:          cfg.Port,
		stdout:        stdout,
		stderr:        stderr,
		onStateChange: cfg.OnStateChange,
	}
}

func (pm *ProcessManager) setState(s ProcessState) {
	pm.state = s
	if pm.onStateChange != nil {
		pm.onStateChange(s)
	}
}

func (pm *ProcessManager) State() ProcessState {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.state
}

func (pm *ProcessManager) LastError() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.lastError
}

func (pm *ProcessManager) ClearLastError() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.lastError = nil
}

func (pm *ProcessManager) RestartCount() int {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.restartCount
}

func (pm *ProcessManager) IsRunning() bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.state == StateRunning && pm.cmd != nil && pm.cmd.Process != nil
}

func (pm *ProcessManager) IsBusy() bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.state == StateStarting || pm.state == StateStopping ||
		pm.state == StateRestarting || pm.state == StateCompiling
}

func (pm *ProcessManager) Start() error {
	pm.mu.Lock()
	if pm.IsBusyLocked() {
		pm.mu.Unlock()
		return fmt.Errorf("cannot start: process manager is busy (state: %s)", pm.state)
	}
	pm.setState(StateStarting)
	pm.mu.Unlock()

	executable, args := pm.buildRunCommand()

	cmd := exec.Command(executable, args...)
	cmd.Stdout = pm.stdout
	cmd.Stderr = pm.stderr
	cmd.Stdin = nil

	setSysProcAttr(cmd)

	if err := cmd.Start(); err != nil {
		pm.mu.Lock()
		pm.lastError = fmt.Errorf("failed to start process: %w", err)
		pm.setState(StateFailed)
		pm.mu.Unlock()
		return pm.lastError
	}

	pm.mu.Lock()
	pm.cmd = cmd
	pm.setState(StateRunning)
	pm.lastError = nil
	pm.mu.Unlock()

	go pm.waitForExit()

	return nil
}

func (pm *ProcessManager) IsBusyLocked() bool {
	return pm.state == StateStarting || pm.state == StateStopping ||
		pm.state == StateRestarting || pm.state == StateCompiling
}

func (pm *ProcessManager) waitForExit() {
	pm.mu.Lock()
	cmd := pm.cmd
	pm.mu.Unlock()

	if cmd == nil {
		return
	}

	err := cmd.Wait()

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.state == StateStopping || pm.state == StateRestarting {
		return
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			pm.lastError = fmt.Errorf("process exited with code %d", exitErr.ExitCode())
		} else {
			pm.lastError = fmt.Errorf("process error: %w", err)
		}
		pm.setState(StateFailed)
	} else {
		pm.setState(StateIdle)
	}
	pm.cmd = nil
}

func (pm *ProcessManager) Stop() error {
	pm.mu.Lock()
	if pm.cmd == nil || pm.cmd.Process == nil {
		pm.setState(StateIdle)
		pm.mu.Unlock()
		return nil
	}
	if pm.state == StateStopping {
		pm.mu.Unlock()
		return nil
	}
	pm.setState(StateStopping)
	cmd := pm.cmd
	pm.mu.Unlock()

	return pm.gracefulStop(cmd)
}

func (pm *ProcessManager) gracefulStop(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		pm.mu.Lock()
		pm.setState(StateIdle)
		pm.cmd = nil
		pm.mu.Unlock()
		return nil
	}

	done := make(chan error, 1)

	go func() {
		done <- cmd.Wait()
	}()

	if err := sendTermSignal(cmd); err != nil {
		if killErr := sendKillSignal(cmd); killErr != nil {
			logger.Debug("Failed to force kill process after SIGTERM failed", "error", killErr)
		}
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		if err := sendKillSignal(cmd); err != nil {
			logger.Debug("Failed to force kill process after timeout", "error", err)
		}
		select {
		case <-done:
		case <-time.After(1 * time.Second):
		}
	}

	pm.mu.Lock()
	pm.cmd = nil
	pm.setState(StateIdle)
	pm.mu.Unlock()

	return nil
}

func (pm *ProcessManager) Compile(ctx context.Context) error {
	pm.mu.Lock()
	pm.setState(StateCompiling)
	pm.mu.Unlock()

	executable, args := pm.buildCompileCommand()

	cmd := exec.CommandContext(ctx, executable, args...)
	cmd.Stdout = pm.stdout
	cmd.Stderr = pm.stderr

	err := cmd.Run()

	pm.mu.Lock()
	if err != nil {
		pm.lastError = fmt.Errorf("compilation failed: %w", err)
		pm.setState(StateFailed)
	}
	pm.mu.Unlock()

	return err
}

func (pm *ProcessManager) Restart(ctx context.Context) error {
	pm.mu.Lock()
	if pm.IsBusyLocked() {
		pm.mu.Unlock()
		return fmt.Errorf("cannot restart: process manager is busy (state: %s)", pm.state)
	}
	pm.setState(StateRestarting)
	pm.restartCount++
	pm.mu.Unlock()

	_, _ = fmt.Fprintln(pm.stdout, "\n\033[33m─────────────────────────────────────────\033[0m")
	_, _ = fmt.Fprintln(pm.stdout, "\033[33m→ Compiling...\033[0m")

	if err := pm.Compile(ctx); err != nil {
		_, _ = fmt.Fprintf(pm.stderr, "\n\033[31m✗ Compilation failed: %v\033[0m\n", err)
		_, _ = fmt.Fprintln(pm.stdout, "\033[33m→ Keeping current server running\033[0m")
		_, _ = fmt.Fprintln(pm.stdout, "\033[33m─────────────────────────────────────────\033[0m")
		pm.mu.Lock()
		pm.setState(StateRunning)
		pm.mu.Unlock()
		return err
	}

	_, _ = fmt.Fprintln(pm.stdout, "\n\033[32m✓ Compilation successful\033[0m")
	_, _ = fmt.Fprintln(pm.stdout, "\033[33m→ Stopping server...\033[0m")

	pm.mu.Lock()
	cmd := pm.cmd
	pm.mu.Unlock()

	if cmd != nil && cmd.Process != nil {
		if err := pm.gracefulStop(cmd); err != nil {
			_, _ = fmt.Fprintf(pm.stderr, "\033[31m✗ Failed to stop server: %v\033[0m\n", err)
		}
	}

	_, _ = fmt.Fprintln(pm.stdout, "\033[33m→ Starting server...\033[0m")
	_, _ = fmt.Fprintln(pm.stdout, "\033[33m─────────────────────────────────────────\033[0m")

	return pm.Start()
}

func (pm *ProcessManager) buildRunCommand() (string, []string) {
	switch pm.buildTool {
	case buildtool.Maven:
		return pm.buildMavenRunCommand()
	case buildtool.Gradle, buildtool.GradleKotln:
		return pm.buildGradleRunCommand()
	default:
		return pm.buildMavenRunCommand()
	}
}

func (pm *ProcessManager) buildMavenRunCommand() (string, []string) {
	executable := getMavenExecutable()
	args := []string{"spring-boot:run", "-DskipTests", "-B", "-ntp"}

	if pm.profile != "" {
		args = append(args, fmt.Sprintf("-Dspring-boot.run.profiles=%s", pm.profile))
	}
	if pm.debug {
		args = append(args, "-Dspring-boot.run.jvmArguments=-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005")
	}
	if pm.port > 0 {
		args = append(args, fmt.Sprintf("-Dspring-boot.run.arguments=--server.port=%d", pm.port))
	}

	return executable, args
}

func (pm *ProcessManager) buildGradleRunCommand() (string, []string) {
	executable := getGradleExecutable()
	args := []string{"bootRun", "-x", "test"}

	if pm.profile != "" {
		args = append(args, fmt.Sprintf("--args=--spring.profiles.active=%s", pm.profile))
	}
	if pm.debug {
		args = append(args, "--debug-jvm")
	}
	if pm.port > 0 {
		if pm.profile != "" {
			for i, arg := range args {
				if arg == fmt.Sprintf("--args=--spring.profiles.active=%s", pm.profile) {
					args[i] = fmt.Sprintf("--args=--spring.profiles.active=%s --server.port=%d", pm.profile, pm.port)
					break
				}
			}
		} else {
			args = append(args, fmt.Sprintf("--args=--server.port=%d", pm.port))
		}
	}

	return executable, args
}

func (pm *ProcessManager) buildCompileCommand() (string, []string) {
	switch pm.buildTool {
	case buildtool.Maven:
		return getMavenExecutable(), []string{"compile", "-DskipTests", "-q", "-B"}
	case buildtool.Gradle, buildtool.GradleKotln:
		return getGradleExecutable(), []string{"classes", "-x", "test", "-q"}
	default:
		return getMavenExecutable(), []string{"compile", "-DskipTests", "-q", "-B"}
	}
}

func (pm *ProcessManager) Kill() error {
	pm.mu.Lock()
	cmd := pm.cmd
	pm.mu.Unlock()

	if cmd != nil && cmd.Process != nil {
		if err := sendKillSignal(cmd); err != nil {
			logger.Debug("Failed to kill process", "error", err)
		}
	}

	pm.mu.Lock()
	pm.cmd = nil
	pm.setState(StateIdle)
	pm.mu.Unlock()

	return nil
}
