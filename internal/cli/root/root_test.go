package root

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommandExecutes(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{})

	err := rootCmd.Execute()

	assert.NoError(t, err)
}

func TestVersionCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"version"})

	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "haft version")
}

func TestRootHasExpectedSubcommands(t *testing.T) {
	commands := rootCmd.Commands()
	commandNames := make([]string, len(commands))
	for i, cmd := range commands {
		commandNames[i] = cmd.Name()
	}

	assert.Contains(t, commandNames, "version")
}

func TestGlobalFlagsExist(t *testing.T) {
	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	noColorFlag := rootCmd.PersistentFlags().Lookup("no-color")

	assert.NotNil(t, verboseFlag)
	assert.NotNil(t, noColorFlag)
	assert.Empty(t, verboseFlag.Shorthand)
}

func TestVersionFlag(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--version"})

	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "haft version")
}

func TestGetVersion(t *testing.T) {
	v := GetVersion()
	assert.NotEmpty(t, v)
	assert.Equal(t, "dev", v)
}

func TestSetVersionUpdatesRootCmd(t *testing.T) {
	originalVersion := version
	defer func() {
		version = originalVersion
		rootCmd.Version = originalVersion
	}()

	SetVersion("v1.2.3")

	assert.Equal(t, "v1.2.3", version)
	assert.Equal(t, "v1.2.3", rootCmd.Version)
}

func TestVerboseFlagSetsState(t *testing.T) {
	originalVerbose := verbose
	defer func() { verbose = originalVerbose }()

	rootCmd.SetArgs([]string{"--verbose", "version"})
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)

	_ = rootCmd.Execute()

	assert.True(t, IsVerbose())
}

func TestNoColorFlagSetsState(t *testing.T) {
	originalNoColor := noColor
	defer func() { noColor = originalNoColor }()

	rootCmd.SetArgs([]string{"--no-color", "version"})
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)

	_ = rootCmd.Execute()

	assert.True(t, IsNoColor())
}

func TestInitLoggerCalled(t *testing.T) {
	originalVerbose := verbose
	originalNoColor := noColor
	defer func() {
		verbose = originalVerbose
		noColor = originalNoColor
	}()

	rootCmd.SetArgs([]string{"version"})
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)

	err := rootCmd.Execute()

	assert.NoError(t, err)
}
