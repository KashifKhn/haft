package dev

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/cobra"
)

func newRestartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Trigger a restart of the running dev server",
		Long: `Trigger a restart of the running dev server (haft dev serve).

This command creates a trigger file that signals the running dev server to 
restart. Use this from scripts, editor plugins, or other tools that want to
trigger a restart without direct keyboard access.

The dev server must be running (haft dev serve) for this command to have effect.`,
		Example: `  # Trigger restart of running dev server
  haft dev restart

  # Use in a shell script
  #!/bin/bash
  # Edit some files...
  haft dev restart

  # Use with file watchers
  fswatch -o src/ | xargs -n1 -I{} haft dev restart`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRestart()
		},
	}

	return cmd
}

func runRestart() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current directory: %w", err)
	}

	dirPath := filepath.Join(cwd, triggerDir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("could not create trigger directory: %w", err)
	}

	triggerPath := filepath.Join(dirPath, triggerFileName)

	file, err := os.Create(triggerPath)
	if err != nil {
		return fmt.Errorf("could not create trigger file: %w", err)
	}
	file.Close()

	logger.Success("Restart triggered", "trigger", triggerPath)
	return nil
}
