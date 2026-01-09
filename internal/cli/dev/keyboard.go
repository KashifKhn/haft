package dev

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"golang.org/x/term"
)

type KeyCommand int

const (
	KeyNone KeyCommand = iota
	KeyRestart
	KeyQuit
	KeyClear
	KeyHelp
	KeyUnknown
)

func (k KeyCommand) String() string {
	switch k {
	case KeyRestart:
		return "restart"
	case KeyQuit:
		return "quit"
	case KeyClear:
		return "clear"
	case KeyHelp:
		return "help"
	default:
		return "unknown"
	}
}

type KeyboardListener struct {
	fd       int
	oldState *term.State
	enabled  bool
	commands chan KeyCommand
	done     chan struct{}
}

func NewKeyboardListener() *KeyboardListener {
	return &KeyboardListener{
		fd:       int(os.Stdin.Fd()),
		commands: make(chan KeyCommand, 10),
		done:     make(chan struct{}),
	}
}

func (kl *KeyboardListener) IsInteractive() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}

func (kl *KeyboardListener) Start() error {
	if !kl.IsInteractive() {
		return nil
	}

	oldState, err := term.MakeRaw(kl.fd)
	if err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}

	kl.oldState = oldState
	kl.enabled = true

	go kl.readLoop()

	return nil
}

func (kl *KeyboardListener) Stop() error {
	if !kl.enabled || kl.oldState == nil {
		return nil
	}

	close(kl.done)
	kl.enabled = false

	return term.Restore(kl.fd, kl.oldState)
}

func (kl *KeyboardListener) Commands() <-chan KeyCommand {
	return kl.commands
}

func (kl *KeyboardListener) readLoop() {
	buf := make([]byte, 3)

	for {
		select {
		case <-kl.done:
			return
		default:
		}

		n, err := os.Stdin.Read(buf)
		if err != nil {
			continue
		}

		if n == 0 {
			continue
		}

		cmd := kl.parseKey(buf[:n])
		if cmd != KeyNone {
			select {
			case kl.commands <- cmd:
			case <-kl.done:
				return
			default:
			}
		}
	}
}

func (kl *KeyboardListener) parseKey(buf []byte) KeyCommand {
	if len(buf) == 0 {
		return KeyNone
	}

	if len(buf) == 1 && buf[0] == 3 {
		return KeyQuit
	}

	if len(buf) >= 1 {
		switch buf[0] {
		case 'r', 'R':
			return KeyRestart
		case 'q', 'Q':
			return KeyQuit
		case 'c', 'C':
			return KeyClear
		case 'h', 'H', '?':
			return KeyHelp
		}
	}

	return KeyNone
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func PrintKeyCommands() {
	fmt.Print("\r\n")
	fmt.Print("\033[36mHaft dev commands:\033[0m\r\n")
	fmt.Print("  \033[33mr\033[0m  Restart (compile first, then restart)\r\n")
	fmt.Print("  \033[33mq\033[0m  Quit\r\n")
	fmt.Print("  \033[33mc\033[0m  Clear screen\r\n")
	fmt.Print("  \033[33mh\033[0m  Help\r\n")
	fmt.Print("\r\n")
}

func PrintBanner() {
	fmt.Print("\r\n")
	fmt.Print("\033[36m╭─────────────────────────────────────╮\033[0m\r\n")
	fmt.Print("\033[36m│\033[0m  \033[1;35mHaft Dev Server\033[0m                    \033[36m│\033[0m\r\n")
	fmt.Print("\033[36m│\033[0m  Press \033[33mr\033[0m to restart, \033[33mq\033[0m to quit      \033[36m│\033[0m\r\n")
	fmt.Print("\033[36m│\033[0m  Press \033[33mh\033[0m for more commands          \033[36m│\033[0m\r\n")
	fmt.Print("\033[36m╰─────────────────────────────────────╯\033[0m\r\n")
	fmt.Print("\r\n")
}
