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
