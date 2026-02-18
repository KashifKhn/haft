package dev

import (
	"os"
	"path/filepath"
	"time"
)

const (
	triggerDir      = ".haft"
	triggerFileName = "restart"
)

type TriggerWatcher struct {
	triggerPath  string
	events       chan struct{}
	done         chan struct{}
	lastModTime  time.Time
	pollInterval time.Duration
}

func NewTriggerWatcher() *TriggerWatcher {
	return &TriggerWatcher{
		events:       make(chan struct{}, 1),
		done:         make(chan struct{}),
		pollInterval: 500 * time.Millisecond,
	}
}

func (tw *TriggerWatcher) Setup() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dirPath := filepath.Join(cwd, triggerDir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	tw.triggerPath = filepath.Join(dirPath, triggerFileName)

	tw.cleanup()

	go tw.watchLoop()

	return nil
}

func (tw *TriggerWatcher) cleanup() {
	_ = os.Remove(tw.triggerPath)
}

func (tw *TriggerWatcher) Events() <-chan struct{} {
	return tw.events
}

func (tw *TriggerWatcher) watchLoop() {
	ticker := time.NewTicker(tw.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-tw.done:
			return
		case <-ticker.C:
			tw.checkTrigger()
		}
	}
}

func (tw *TriggerWatcher) checkTrigger() {
	info, err := os.Stat(tw.triggerPath)
	if err != nil {
		return
	}

	modTime := info.ModTime()
	if modTime.After(tw.lastModTime) {
		tw.lastModTime = modTime

		_ = os.Remove(tw.triggerPath)

		select {
		case tw.events <- struct{}{}:
		default:
		}
	}
}

func (tw *TriggerWatcher) Cleanup() {
	close(tw.done)
	tw.cleanup()
}

func (tw *TriggerWatcher) TriggerPath() string {
	return tw.triggerPath
}

func GetTriggerPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, triggerDir, triggerFileName), nil
}
