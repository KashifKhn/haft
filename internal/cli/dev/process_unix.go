//go:build !windows

package dev

import (
	"os/exec"
	"syscall"

	"github.com/KashifKhn/haft/internal/logger"
)

func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func sendTermSignal(cmd *exec.Cmd) error {
	return cmd.Process.Signal(syscall.SIGTERM)
}

func sendKillSignal(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		if killErr := syscall.Kill(-pgid, syscall.SIGKILL); killErr != nil {
			logger.Debug("Failed to kill process group", "pgid", pgid, "error", killErr)
		}
	}
	return cmd.Process.Kill()
}
