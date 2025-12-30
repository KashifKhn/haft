//go:build windows

package dev

import (
	"os/exec"
)

func setSysProcAttr(cmd *exec.Cmd) {
}

func sendTermSignal(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}

func sendKillSignal(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	return cmd.Process.Kill()
}
