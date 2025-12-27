package completion

import (
	"bytes"
	"strings"
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

	assert.Equal(t, "completion", cmd.Use)
	assert.Equal(t, "Generate shell completion scripts", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestCompletionSubcommands(t *testing.T) {
	cmd := NewCommand()

	subcommands := make(map[string]*cobra.Command)
	for _, sub := range cmd.Commands() {
		subcommands[sub.Name()] = sub
	}

	expectedSubs := []string{"bash", "zsh", "fish", "powershell"}
	for _, name := range expectedSubs {
		assert.Contains(t, subcommands, name, "missing subcommand: %s", name)
	}
}

func TestBashCompletion(t *testing.T) {
	cmd := NewCommand()
	buf := new(bytes.Buffer)

	rootCmd := &cobra.Command{Use: "haft"}
	rootCmd.AddCommand(cmd)
	rootCmd.SetOut(buf)

	err := rootCmd.GenBashCompletion(buf)
	require.NoError(t, err)

	output := buf.String()
	assert.True(t, strings.Contains(output, "bash completion") || strings.Contains(output, "__haft") || strings.Contains(output, "_haft"),
		"expected bash completion script")
}

func TestZshCompletion(t *testing.T) {
	cmd := NewCommand()
	buf := new(bytes.Buffer)

	rootCmd := &cobra.Command{Use: "haft"}
	rootCmd.AddCommand(cmd)
	rootCmd.SetOut(buf)

	err := rootCmd.GenZshCompletion(buf)
	require.NoError(t, err)

	output := buf.String()
	assert.True(t, strings.Contains(output, "compdef") || strings.Contains(output, "#compdef") || strings.Contains(output, "_haft"),
		"expected zsh completion script")
}

func TestFishCompletion(t *testing.T) {
	cmd := NewCommand()
	buf := new(bytes.Buffer)

	rootCmd := &cobra.Command{Use: "haft"}
	rootCmd.AddCommand(cmd)
	rootCmd.SetOut(buf)

	err := rootCmd.GenFishCompletion(buf, true)
	require.NoError(t, err)

	output := buf.String()
	assert.True(t, strings.Contains(output, "complete") || strings.Contains(output, "fish"),
		"expected fish completion script")
}

func TestPowershellCompletion(t *testing.T) {
	cmd := NewCommand()
	buf := new(bytes.Buffer)

	rootCmd := &cobra.Command{Use: "haft"}
	rootCmd.AddCommand(cmd)
	rootCmd.SetOut(buf)

	err := rootCmd.GenPowerShellCompletionWithDesc(buf)
	require.NoError(t, err)

	output := buf.String()
	assert.True(t, strings.Contains(output, "PowerShell") || strings.Contains(output, "Register-ArgumentCompleter"),
		"expected powershell completion script")
}

func TestBashSubcommand(t *testing.T) {
	rootCmd := newTestRootCmd()
	rootCmd.SetArgs([]string{"completion", "bash"})

	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestZshSubcommand(t *testing.T) {
	rootCmd := newTestRootCmd()
	rootCmd.SetArgs([]string{"completion", "zsh"})

	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestFishSubcommand(t *testing.T) {
	rootCmd := newTestRootCmd()
	rootCmd.SetArgs([]string{"completion", "fish"})

	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestPowershellSubcommand(t *testing.T) {
	rootCmd := newTestRootCmd()
	rootCmd.SetArgs([]string{"completion", "powershell"})

	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestInvalidShell(t *testing.T) {
	rootCmd := newTestRootCmd()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"completion", "invalid"})

	err := rootCmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Available Commands")
}

func TestNoArgs(t *testing.T) {
	rootCmd := newTestRootCmd()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"completion"})

	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func TestTooManyArgs(t *testing.T) {
	rootCmd := newTestRootCmd()
	rootCmd.SetArgs([]string{"completion", "bash", "zsh"})

	err := rootCmd.Execute()
	assert.Error(t, err)
}

func TestSubcommandHelp(t *testing.T) {
	tests := []struct {
		name    string
		subcmd  string
		wantUse string
	}{
		{"bash", "bash", "bash"},
		{"zsh", "zsh", "zsh"},
		{"fish", "fish", "fish"},
		{"powershell", "powershell", "powershell"},
	}

	cmd := NewCommand()
	subcommands := make(map[string]*cobra.Command)
	for _, sub := range cmd.Commands() {
		subcommands[sub.Name()] = sub
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := subcommands[tt.subcmd]
			require.NotNil(t, sub)
			assert.Equal(t, tt.wantUse, sub.Use)
			assert.NotEmpty(t, sub.Short)
			assert.NotEmpty(t, sub.Long)
			assert.NotEmpty(t, sub.Example)
		})
	}
}
