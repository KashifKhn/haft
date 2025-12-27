package dev

import (
	"runtime"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "haft",
		Short: "Test root command",
	}
	rootCmd.AddCommand(NewCommand())
	return rootCmd
}

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

	expectedSubs := []string{"serve", "build", "test", "clean"}
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
