---
sidebar_position: 5
title: haft completion
description: Generate shell completion scripts
---

# haft completion

Generate shell completion scripts for Haft.

## Usage

```bash
haft completion <shell>
haft completion bash
haft completion zsh
haft completion fish
haft completion powershell
```

## Description

The `completion` command generates shell-specific completion scripts that enable tab completion for Haft commands, subcommands, and flags. This improves productivity by allowing you to discover available options without memorizing them.

## Supported Shells

| Shell | Command |
|-------|---------|
| Bash | `haft completion bash` |
| Zsh | `haft completion zsh` |
| Fish | `haft completion fish` |
| PowerShell | `haft completion powershell` |

## Installation

### Bash

**Load for current session:**

```bash
source <(haft completion bash)
```

**Load permanently (Linux):**

```bash
# System-wide (requires sudo)
haft completion bash > /etc/bash_completion.d/haft

# User-only
mkdir -p ~/.local/share/bash-completion/completions
haft completion bash > ~/.local/share/bash-completion/completions/haft
```

**Load permanently (macOS with Homebrew):**

```bash
haft completion bash > $(brew --prefix)/etc/bash_completion.d/haft
```

### Zsh

**Load for current session:**

```bash
source <(haft completion zsh)
```

**Load permanently (Linux):**

```bash
# First, ensure completion is enabled in ~/.zshrc:
# autoload -U compinit; compinit

haft completion zsh > "${fpath[1]}/_haft"
```

**Load permanently (macOS with Homebrew):**

```bash
haft completion zsh > $(brew --prefix)/share/zsh/site-functions/_haft
```

Start a new shell session for changes to take effect.

### Fish

**Load for current session:**

```bash
haft completion fish | source
```

**Load permanently:**

```bash
haft completion fish > ~/.config/fish/completions/haft.fish
```

### PowerShell

**Load for current session:**

```powershell
haft completion powershell | Out-String | Invoke-Expression
```

**Load permanently:**

Add to your PowerShell profile (`$PROFILE`):

```powershell
haft completion powershell | Out-String | Invoke-Expression
```

Or save to a file and source it:

```powershell
haft completion powershell > haft.ps1
# Add to your profile: . /path/to/haft.ps1
```

## Examples

```bash
# Generate and redirect to file
haft completion bash > haft-completion.bash

# Generate and immediately source (bash/zsh)
source <(haft completion bash)

# Pipe to source (fish)
haft completion fish | source

# View completion script without installing
haft completion zsh | less
```

## What Gets Completed

Once installed, tab completion works for:

- **Commands**: `haft <TAB>` shows `init`, `generate`, `add`, `remove`, `completion`, `version`
- **Subcommands**: `haft generate <TAB>` shows `resource`, `controller`, `service`, etc.
- **Flags**: `haft add --<TAB>` shows `--browse`, `--list`, `--scope`, etc.
- **Flag values**: Context-aware completion for flag arguments

## Troubleshooting

### Bash completions not working

Ensure bash-completion is installed:

```bash
# Debian/Ubuntu
sudo apt install bash-completion

# macOS with Homebrew
brew install bash-completion@2
```

### Zsh completions not working

Ensure completion system is initialized in `~/.zshrc`:

```bash
autoload -U compinit; compinit
```

### Permission denied errors

Use user-local directories instead of system directories:

```bash
# Bash
mkdir -p ~/.local/share/bash-completion/completions
haft completion bash > ~/.local/share/bash-completion/completions/haft

# Zsh (add custom directory to fpath in ~/.zshrc)
mkdir -p ~/.zsh/completions
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
haft completion zsh > ~/.zsh/completions/_haft
```

## Editor Integration

Shell completion is a CLI feature. For editor integration:

- **Neovim**: haft.nvim provides native command completion via Telescope ([docs →](/docs/integrations/neovim/usage#telescope-integration))
- **VS Code**: Coming soon ([preview →](/docs/integrations/vscode))
- **IntelliJ IDEA**: Coming soon ([preview →](/docs/integrations/intellij))

## See Also

- [Installation](/docs/installation#shell-completions) - Installation guide with completion setup
- [haft init](/docs/commands/init) - Initialize new projects
- [haft generate](/docs/commands/generate) - Generate code
