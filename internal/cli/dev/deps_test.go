package dev

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDepsCommand_HasCorrectDescription(t *testing.T) {
	cmd := newDepsCommand()

	assert.Equal(t, "deps", cmd.Use)
	assert.Contains(t, cmd.Aliases, "dependencies")
	assert.Contains(t, cmd.Aliases, "tree")
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestDepsCommand_HasCorrectFlags(t *testing.T) {
	cmd := newDepsCommand()

	configFlag := cmd.Flags().Lookup("configuration")
	assert.NotNil(t, configFlag)
	assert.Equal(t, "c", configFlag.Shorthand)
	assert.Empty(t, configFlag.DefValue)

	verboseFlag := cmd.Flags().Lookup("verbose")
	assert.NotNil(t, verboseFlag)
	assert.Equal(t, "v", verboseFlag.Shorthand)
	assert.Equal(t, "false", verboseFlag.DefValue)
}

func TestDepsCommand_FlagCombinations(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no flags",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "configuration flag",
			args:    []string{"--configuration", "compile"},
			wantErr: false,
		},
		{
			name:    "configuration short flag",
			args:    []string{"-c", "runtime"},
			wantErr: false,
		},
		{
			name:    "verbose flag",
			args:    []string{"--verbose"},
			wantErr: false,
		},
		{
			name:    "verbose short flag",
			args:    []string{"-v"},
			wantErr: false,
		},
		{
			name:    "combined flags",
			args:    []string{"-c", "compileClasspath", "-v"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newDepsCommand()
			cmd.RunE = func(cmd *cobra.Command, args []string) error {
				return nil
			}
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDepsCommand_IsRegistered(t *testing.T) {
	devCmd := NewCommand()
	var found bool
	for _, cmd := range devCmd.Commands() {
		if cmd.Use == "deps" {
			found = true
			break
		}
	}
	assert.True(t, found, "deps command should be registered as a subcommand of dev")
}

func TestDepsCommand_Aliases(t *testing.T) {
	cmd := newDepsCommand()

	assert.Contains(t, cmd.Aliases, "dependencies")
	assert.Contains(t, cmd.Aliases, "tree")
	assert.Len(t, cmd.Aliases, 2)
}

func TestDepsCommand_LongDescription(t *testing.T) {
	cmd := newDepsCommand()

	assert.Contains(t, cmd.Long, "dependency tree")
	assert.Contains(t, cmd.Long, "transitive")
	assert.Contains(t, cmd.Long, "mvn dependency:tree")
	assert.Contains(t, cmd.Long, "gradlew dependencies")
}

func TestDepsCommand_Examples(t *testing.T) {
	cmd := newDepsCommand()

	assert.Contains(t, cmd.Example, "haft dev deps")
	assert.Contains(t, cmd.Example, "--configuration")
	assert.Contains(t, cmd.Example, "--verbose")
}
