package dev

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestVerifyCommand_HasCorrectDescription(t *testing.T) {
	cmd := newVerifyCommand()

	assert.Equal(t, "verify", cmd.Use)
	assert.Contains(t, cmd.Aliases, "vfy")
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestVerifyCommand_HasCorrectFlags(t *testing.T) {
	cmd := newVerifyCommand()

	skipTestsFlag := cmd.Flags().Lookup("skip-tests")
	assert.NotNil(t, skipTestsFlag)
	assert.Equal(t, "s", skipTestsFlag.Shorthand)
	assert.Equal(t, "false", skipTestsFlag.DefValue)

	skipIntegrationFlag := cmd.Flags().Lookup("skip-integration")
	assert.NotNil(t, skipIntegrationFlag)
	assert.Equal(t, "i", skipIntegrationFlag.Shorthand)
	assert.Equal(t, "false", skipIntegrationFlag.DefValue)

	profileFlag := cmd.Flags().Lookup("profile")
	assert.NotNil(t, profileFlag)
	assert.Equal(t, "p", profileFlag.Shorthand)
	assert.Empty(t, profileFlag.DefValue)
}

func TestVerifyCommand_FlagCombinations(t *testing.T) {
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
			name:    "skip-tests flag",
			args:    []string{"--skip-tests"},
			wantErr: false,
		},
		{
			name:    "skip-tests short flag",
			args:    []string{"-s"},
			wantErr: false,
		},
		{
			name:    "skip-integration flag",
			args:    []string{"--skip-integration"},
			wantErr: false,
		},
		{
			name:    "skip-integration short flag",
			args:    []string{"-i"},
			wantErr: false,
		},
		{
			name:    "profile flag",
			args:    []string{"--profile", "ci"},
			wantErr: false,
		},
		{
			name:    "profile short flag",
			args:    []string{"-p", "prod"},
			wantErr: false,
		},
		{
			name:    "combined flags",
			args:    []string{"-s", "-p", "ci"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newVerifyCommand()
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

func TestVerifyCommand_IsRegistered(t *testing.T) {
	devCmd := NewCommand()
	var found bool
	for _, cmd := range devCmd.Commands() {
		if cmd.Use == "verify" {
			found = true
			break
		}
	}
	assert.True(t, found, "verify command should be registered as a subcommand of dev")
}

func TestVerifyCommand_Aliases(t *testing.T) {
	cmd := newVerifyCommand()

	assert.Contains(t, cmd.Aliases, "vfy")
	assert.Len(t, cmd.Aliases, 1)
}

func TestVerifyCommand_LongDescription(t *testing.T) {
	cmd := newVerifyCommand()

	assert.Contains(t, cmd.Long, "integration tests")
	assert.Contains(t, cmd.Long, "quality checks")
	assert.Contains(t, cmd.Long, "mvn verify")
	assert.Contains(t, cmd.Long, "gradlew check")
}

func TestVerifyCommand_Examples(t *testing.T) {
	cmd := newVerifyCommand()

	assert.Contains(t, cmd.Example, "haft dev verify")
	assert.Contains(t, cmd.Example, "--skip-tests")
	assert.Contains(t, cmd.Example, "--skip-integration")
	assert.Contains(t, cmd.Example, "--profile")
}
