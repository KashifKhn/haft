package dev

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPackageCommand_HasCorrectDescription(t *testing.T) {
	cmd := newPackageCommand()

	assert.Equal(t, "package", cmd.Use)
	assert.Contains(t, cmd.Aliases, "pkg")
	assert.Contains(t, cmd.Aliases, "jar")
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestPackageCommand_HasCorrectFlags(t *testing.T) {
	cmd := newPackageCommand()

	skipTestsFlag := cmd.Flags().Lookup("skip-tests")
	assert.NotNil(t, skipTestsFlag)
	assert.Equal(t, "s", skipTestsFlag.Shorthand)
	assert.Equal(t, "true", skipTestsFlag.DefValue)

	cleanFlag := cmd.Flags().Lookup("clean")
	assert.NotNil(t, cleanFlag)
	assert.Equal(t, "c", cleanFlag.Shorthand)
	assert.Equal(t, "false", cleanFlag.DefValue)

	profileFlag := cmd.Flags().Lookup("profile")
	assert.NotNil(t, profileFlag)
	assert.Equal(t, "p", profileFlag.Shorthand)
	assert.Empty(t, profileFlag.DefValue)
}

func TestPackageCommand_SkipTestsDefaultTrue(t *testing.T) {
	cmd := newPackageCommand()

	skipTestsFlag := cmd.Flags().Lookup("skip-tests")
	assert.Equal(t, "true", skipTestsFlag.DefValue, "skip-tests should default to true")
}

func TestPackageCommand_FlagCombinations(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no flags (defaults to skip tests)",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "explicitly skip tests",
			args:    []string{"--skip-tests"},
			wantErr: false,
		},
		{
			name:    "run tests",
			args:    []string{"--skip-tests=false"},
			wantErr: false,
		},
		{
			name:    "clean flag",
			args:    []string{"--clean"},
			wantErr: false,
		},
		{
			name:    "clean short flag",
			args:    []string{"-c"},
			wantErr: false,
		},
		{
			name:    "profile flag",
			args:    []string{"--profile", "prod"},
			wantErr: false,
		},
		{
			name:    "profile short flag",
			args:    []string{"-p", "staging"},
			wantErr: false,
		},
		{
			name:    "combined flags",
			args:    []string{"-c", "-p", "prod"},
			wantErr: false,
		},
		{
			name:    "all flags",
			args:    []string{"-c", "-p", "prod", "--skip-tests=false"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newPackageCommand()
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

func TestPackageCommand_IsRegistered(t *testing.T) {
	devCmd := NewCommand()
	var found bool
	for _, cmd := range devCmd.Commands() {
		if cmd.Use == "package" {
			found = true
			break
		}
	}
	assert.True(t, found, "package command should be registered as a subcommand of dev")
}

func TestPackageCommand_Aliases(t *testing.T) {
	cmd := newPackageCommand()

	assert.Contains(t, cmd.Aliases, "pkg")
	assert.Contains(t, cmd.Aliases, "jar")
	assert.Len(t, cmd.Aliases, 2)
}

func TestPackageCommand_LongDescription(t *testing.T) {
	cmd := newPackageCommand()

	assert.Contains(t, cmd.Long, "deployable artifact")
	assert.Contains(t, cmd.Long, "JAR/WAR")
	assert.Contains(t, cmd.Long, "mvn package")
	assert.Contains(t, cmd.Long, "bootJar")
}

func TestPackageCommand_Examples(t *testing.T) {
	cmd := newPackageCommand()

	assert.Contains(t, cmd.Example, "haft dev package")
	assert.Contains(t, cmd.Example, "--clean")
	assert.Contains(t, cmd.Example, "--profile")
}
