package dev

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestOutdatedCommand_HasCorrectDescription(t *testing.T) {
	cmd := newOutdatedCommand()

	assert.Equal(t, "outdated", cmd.Use)
	assert.Contains(t, cmd.Aliases, "updates")
	assert.Contains(t, cmd.Aliases, "out")
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestOutdatedCommand_HasCorrectFlags(t *testing.T) {
	cmd := newOutdatedCommand()

	pluginsFlag := cmd.Flags().Lookup("plugins")
	assert.NotNil(t, pluginsFlag)
	assert.Equal(t, "p", pluginsFlag.Shorthand)
	assert.Equal(t, "false", pluginsFlag.DefValue)

	snapshotsFlag := cmd.Flags().Lookup("snapshots")
	assert.NotNil(t, snapshotsFlag)
	assert.Equal(t, "s", snapshotsFlag.Shorthand)
	assert.Equal(t, "false", snapshotsFlag.DefValue)
}

func TestOutdatedCommand_FlagCombinations(t *testing.T) {
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
			name:    "plugins flag",
			args:    []string{"--plugins"},
			wantErr: false,
		},
		{
			name:    "plugins short flag",
			args:    []string{"-p"},
			wantErr: false,
		},
		{
			name:    "snapshots flag",
			args:    []string{"--snapshots"},
			wantErr: false,
		},
		{
			name:    "snapshots short flag",
			args:    []string{"-s"},
			wantErr: false,
		},
		{
			name:    "combined flags",
			args:    []string{"-p", "-s"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newOutdatedCommand()
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

func TestOutdatedCommand_IsRegistered(t *testing.T) {
	devCmd := NewCommand()
	var found bool
	for _, cmd := range devCmd.Commands() {
		if cmd.Use == "outdated" {
			found = true
			break
		}
	}
	assert.True(t, found, "outdated command should be registered as a subcommand of dev")
}

func TestOutdatedCommand_Aliases(t *testing.T) {
	cmd := newOutdatedCommand()

	assert.Contains(t, cmd.Aliases, "updates")
	assert.Contains(t, cmd.Aliases, "out")
	assert.Len(t, cmd.Aliases, 2)
}

func TestOutdatedCommand_LongDescription(t *testing.T) {
	cmd := newOutdatedCommand()

	assert.Contains(t, cmd.Long, "newer versions")
	assert.Contains(t, cmd.Long, "versions:display-dependency-updates")
	assert.Contains(t, cmd.Long, "dependencyUpdates")
}

func TestOutdatedCommand_Examples(t *testing.T) {
	cmd := newOutdatedCommand()

	assert.Contains(t, cmd.Example, "haft dev outdated")
	assert.Contains(t, cmd.Example, "--plugins")
	assert.Contains(t, cmd.Example, "--snapshots")
}
