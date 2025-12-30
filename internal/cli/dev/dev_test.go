package dev

import (
	"runtime"
	"testing"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	assert.Equal(t, "dev", cmd.Use)
	assert.Equal(t, []string{"d"}, cmd.Aliases)
	assert.Equal(t, "Development commands for running, building, and testing", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestDevSubcommands(t *testing.T) {
	cmd := NewCommand()

	subcommands := make(map[string]*cobra.Command)
	for _, sub := range cmd.Commands() {
		subcommands[sub.Name()] = sub
	}

	expectedSubs := []string{"serve", "build", "test", "clean", "restart"}
	for _, name := range expectedSubs {
		assert.Contains(t, subcommands, name, "missing subcommand: %s", name)
	}
}

func TestServeCommand(t *testing.T) {
	cmd := NewCommand()

	var serveCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "serve" {
			serveCmd = sub
			break
		}
	}

	require.NotNil(t, serveCmd)
	assert.Equal(t, "serve", serveCmd.Use)
	assert.Equal(t, []string{"run", "start"}, serveCmd.Aliases)
	assert.NotEmpty(t, serveCmd.Short)
	assert.NotEmpty(t, serveCmd.Long)
	assert.NotEmpty(t, serveCmd.Example)

	profileFlag := serveCmd.Flags().Lookup("profile")
	require.NotNil(t, profileFlag)
	assert.Equal(t, "p", profileFlag.Shorthand)

	debugFlag := serveCmd.Flags().Lookup("debug")
	require.NotNil(t, debugFlag)
	assert.Equal(t, "d", debugFlag.Shorthand)

	portFlag := serveCmd.Flags().Lookup("port")
	require.NotNil(t, portFlag)
}

func TestBuildCommand(t *testing.T) {
	cmd := NewCommand()

	var buildCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "build" {
			buildCmd = sub
			break
		}
	}

	require.NotNil(t, buildCmd)
	assert.Equal(t, "build", buildCmd.Use)
	assert.Equal(t, []string{"b", "compile"}, buildCmd.Aliases)
	assert.NotEmpty(t, buildCmd.Short)
	assert.NotEmpty(t, buildCmd.Long)
	assert.NotEmpty(t, buildCmd.Example)

	skipTestsFlag := buildCmd.Flags().Lookup("skip-tests")
	require.NotNil(t, skipTestsFlag)
	assert.Equal(t, "s", skipTestsFlag.Shorthand)

	profileFlag := buildCmd.Flags().Lookup("profile")
	require.NotNil(t, profileFlag)
	assert.Equal(t, "p", profileFlag.Shorthand)

	cleanFlag := buildCmd.Flags().Lookup("clean")
	require.NotNil(t, cleanFlag)
	assert.Equal(t, "c", cleanFlag.Shorthand)
}

func TestTestCommand(t *testing.T) {
	cmd := NewCommand()

	var testCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "test" {
			testCmd = sub
			break
		}
	}

	require.NotNil(t, testCmd)
	assert.Equal(t, "test", testCmd.Use)
	assert.Equal(t, []string{"t"}, testCmd.Aliases)
	assert.NotEmpty(t, testCmd.Short)
	assert.NotEmpty(t, testCmd.Long)
	assert.NotEmpty(t, testCmd.Example)

	filterFlag := testCmd.Flags().Lookup("filter")
	require.NotNil(t, filterFlag)
	assert.Equal(t, "f", filterFlag.Shorthand)

	verboseFlag := testCmd.Flags().Lookup("verbose")
	require.NotNil(t, verboseFlag)
	assert.Equal(t, "v", verboseFlag.Shorthand)

	failFastFlag := testCmd.Flags().Lookup("fail-fast")
	require.NotNil(t, failFastFlag)
}

func TestCleanCommand(t *testing.T) {
	cmd := NewCommand()

	var cleanCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "clean" {
			cleanCmd = sub
			break
		}
	}

	require.NotNil(t, cleanCmd)
	assert.Equal(t, "clean", cleanCmd.Use)
	assert.NotEmpty(t, cleanCmd.Short)
	assert.NotEmpty(t, cleanCmd.Long)
	assert.NotEmpty(t, cleanCmd.Example)
}

func TestRestartCommand(t *testing.T) {
	cmd := NewCommand()

	var restartCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "restart" {
			restartCmd = sub
			break
		}
	}

	require.NotNil(t, restartCmd)
	assert.Equal(t, "restart", restartCmd.Use)
	assert.NotEmpty(t, restartCmd.Short)
	assert.NotEmpty(t, restartCmd.Long)
	assert.NotEmpty(t, restartCmd.Example)
}

func TestGetMavenExecutable(t *testing.T) {
	exec := getMavenExecutable()
	if runtime.GOOS == "windows" {
		assert.Contains(t, []string{"mvn", "mvnw.cmd"}, exec)
	} else {
		assert.Contains(t, []string{"mvn", "./mvnw"}, exec)
	}
}

func TestGetGradleExecutable(t *testing.T) {
	exec := getGradleExecutable()
	if runtime.GOOS == "windows" {
		assert.Contains(t, []string{"gradle", "gradlew.bat"}, exec)
	} else {
		assert.Contains(t, []string{"gradle", "./gradlew"}, exec)
	}
}

func TestServeCommandHasNoInteractiveFlag(t *testing.T) {
	cmd := NewCommand()

	var serveCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "serve" {
			serveCmd = sub
			break
		}
	}

	require.NotNil(t, serveCmd)

	noInteractiveFlag := serveCmd.Flags().Lookup("no-interactive")
	require.NotNil(t, noInteractiveFlag)
	assert.Equal(t, "false", noInteractiveFlag.DefValue)
}

func TestProcessState_String(t *testing.T) {
	tests := []struct {
		state    ProcessState
		expected string
	}{
		{StateIdle, "idle"},
		{StateStarting, "starting"},
		{StateRunning, "running"},
		{StateStopping, "stopping"},
		{StateRestarting, "restarting"},
		{StateCompiling, "compiling"},
		{StateFailed, "failed"},
		{ProcessState(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.state.String())
		})
	}
}

func TestKeyCommand_String(t *testing.T) {
	tests := []struct {
		cmd      KeyCommand
		expected string
	}{
		{KeyRestart, "restart"},
		{KeyQuit, "quit"},
		{KeyClear, "clear"},
		{KeyHelp, "help"},
		{KeyNone, "unknown"},
		{KeyUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.cmd.String())
		})
	}
}

func TestNewProcessManager(t *testing.T) {
	cfg := ProcessConfig{
		BuildTool: buildtool.Maven,
		Profile:   "dev",
		Debug:     true,
		Port:      8080,
	}

	pm := NewProcessManager(cfg)

	assert.NotNil(t, pm)
	assert.Equal(t, StateIdle, pm.State())
	assert.False(t, pm.IsRunning())
	assert.False(t, pm.IsBusy())
	assert.Equal(t, 0, pm.RestartCount())
	assert.Nil(t, pm.LastError())
}

func TestProcessManager_StateTransitions(t *testing.T) {
	var lastState ProcessState
	cfg := ProcessConfig{
		OnStateChange: func(s ProcessState) {
			lastState = s
		},
	}

	pm := NewProcessManager(cfg)

	pm.mu.Lock()
	pm.setState(StateStarting)
	pm.mu.Unlock()

	assert.Equal(t, StateStarting, pm.State())
	assert.Equal(t, StateStarting, lastState)
}

func TestProcessManager_IsBusy(t *testing.T) {
	pm := NewProcessManager(ProcessConfig{})

	assert.False(t, pm.IsBusy())

	busyStates := []ProcessState{StateStarting, StateStopping, StateRestarting, StateCompiling}
	for _, state := range busyStates {
		pm.mu.Lock()
		pm.state = state
		pm.mu.Unlock()
		assert.True(t, pm.IsBusy(), "should be busy in state: %s", state)
	}

	notBusyStates := []ProcessState{StateIdle, StateRunning, StateFailed}
	for _, state := range notBusyStates {
		pm.mu.Lock()
		pm.state = state
		pm.mu.Unlock()
		assert.False(t, pm.IsBusy(), "should not be busy in state: %s", state)
	}
}

func TestNewKeyboardListener(t *testing.T) {
	kl := NewKeyboardListener()

	assert.NotNil(t, kl)
	assert.NotNil(t, kl.commands)
	assert.NotNil(t, kl.done)
}

func TestKeyboardListener_ParseKey(t *testing.T) {
	kl := NewKeyboardListener()

	tests := []struct {
		name     string
		input    []byte
		expected KeyCommand
	}{
		{"restart lowercase", []byte{'r'}, KeyRestart},
		{"restart uppercase", []byte{'R'}, KeyRestart},
		{"quit lowercase", []byte{'q'}, KeyQuit},
		{"quit uppercase", []byte{'Q'}, KeyQuit},
		{"clear lowercase", []byte{'c'}, KeyClear},
		{"clear uppercase", []byte{'C'}, KeyClear},
		{"help lowercase", []byte{'h'}, KeyHelp},
		{"help uppercase", []byte{'H'}, KeyHelp},
		{"help question mark", []byte{'?'}, KeyHelp},
		{"ctrl+c", []byte{3}, KeyQuit},
		{"unknown key", []byte{'x'}, KeyNone},
		{"empty", []byte{}, KeyNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kl.parseKey(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewTriggerWatcher(t *testing.T) {
	tw := NewTriggerWatcher()

	assert.NotNil(t, tw)
	assert.NotNil(t, tw.events)
	assert.NotNil(t, tw.done)
}

func TestTriggerWatcher_GetTriggerPath(t *testing.T) {
	path, err := GetTriggerPath()

	assert.NoError(t, err)
	assert.Contains(t, path, ".haft")
	assert.Contains(t, path, "restart")
}

func TestProcessManager_BuildMavenRunCommand(t *testing.T) {
	pm := NewProcessManager(ProcessConfig{
		BuildTool: buildtool.Maven,
		Profile:   "dev",
		Debug:     true,
		Port:      8080,
	})

	executable, args := pm.buildMavenRunCommand()

	assert.Contains(t, []string{"mvn", "./mvnw", "mvnw.cmd"}, executable)
	assert.Contains(t, args, "spring-boot:run")
	assert.Contains(t, args, "-DskipTests")
	assert.Contains(t, args, "-Dspring-boot.run.profiles=dev")
	assert.Contains(t, args, "-Dspring-boot.run.arguments=--server.port=8080")
}

func TestProcessManager_BuildGradleRunCommand(t *testing.T) {
	pm := NewProcessManager(ProcessConfig{
		BuildTool: buildtool.Gradle,
		Profile:   "dev",
		Port:      8080,
	})

	executable, args := pm.buildGradleRunCommand()

	assert.Contains(t, []string{"gradle", "./gradlew", "gradlew.bat"}, executable)
	assert.Contains(t, args, "bootRun")
	assert.Contains(t, args, "-x")
	assert.Contains(t, args, "test")
}

func TestProcessManager_BuildCompileCommand(t *testing.T) {
	tests := []struct {
		name      string
		buildTool buildtool.Type
		wantGoal  string
	}{
		{"maven", buildtool.Maven, "compile"},
		{"gradle", buildtool.Gradle, "classes"},
		{"gradle-kotlin", buildtool.GradleKotln, "classes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := NewProcessManager(ProcessConfig{BuildTool: tt.buildTool})
			_, args := pm.buildCompileCommand()
			assert.Contains(t, args, tt.wantGoal)
		})
	}
}
