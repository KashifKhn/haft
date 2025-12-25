---
sidebar_position: 2
title: Installation
description: Install Haft on Linux, macOS, or Windows
---

# Installation

Haft can be installed on Linux, macOS, and Windows.

## Requirements

- **Go 1.21+** (for building from source)
- **Java 17+** (for generated projects)
- **Maven 3.6+** or **Gradle 7+** (for generated projects)

## Quick Install (Recommended)

The easiest way to install Haft:

```bash
curl -fsSL https://raw.githubusercontent.com/KashifKhn/haft/main/install.sh | bash
```

This automatically:
- Detects your OS and architecture
- Downloads the latest release
- Installs to `/usr/local/bin` or `~/.local/bin`

## Using Go

If you have Go installed:

```bash
go install github.com/KashifKhn/haft/cmd/haft@latest
```

This installs the `haft` binary to your `$GOPATH/bin` directory.

## From Source

Clone and build from source:

```bash
git clone https://github.com/KashifKhn/haft.git
cd haft
make build
```

The binary will be at `./bin/haft`. Move it to your PATH:

```bash
sudo mv ./bin/haft /usr/local/bin/
```

## Binary Releases

Download pre-built binaries from [GitHub Releases](https://github.com/KashifKhn/haft/releases).

### Linux

```bash
# AMD64
curl -L https://github.com/KashifKhn/haft/releases/latest/download/haft-linux-amd64.tar.gz | tar xz
sudo mv haft-linux-amd64 /usr/local/bin/haft

# ARM64
curl -L https://github.com/KashifKhn/haft/releases/latest/download/haft-linux-arm64.tar.gz | tar xz
sudo mv haft-linux-arm64 /usr/local/bin/haft
```

### macOS

```bash
# Intel Mac
curl -L https://github.com/KashifKhn/haft/releases/latest/download/haft-darwin-amd64.tar.gz | tar xz
sudo mv haft-darwin-amd64 /usr/local/bin/haft

# Apple Silicon
curl -L https://github.com/KashifKhn/haft/releases/latest/download/haft-darwin-arm64.tar.gz | tar xz
sudo mv haft-darwin-arm64 /usr/local/bin/haft
```

### Windows

1. Download `haft-windows-amd64.zip` from the [releases page](https://github.com/KashifKhn/haft/releases)
2. Extract the ZIP file
3. Add the extracted folder to your system PATH

Or using PowerShell:

```powershell
# Download and extract
Invoke-WebRequest -Uri "https://github.com/KashifKhn/haft/releases/latest/download/haft-windows-amd64.zip" -OutFile "haft.zip"
Expand-Archive -Path "haft.zip" -DestinationPath "."
Move-Item "haft-windows-amd64.exe" "$env:LOCALAPPDATA\Microsoft\WindowsApps\haft.exe"
```

## Verify Installation

```bash
haft version
```

You should see output like:

```
haft version v0.3.0
```

## Shell Completions

Haft supports shell completions for bash, zsh, fish, and PowerShell. Tab completions help you discover commands and flags without memorizing them.

### Bash

**Load for current session:**

```bash
source <(haft completion bash)
```

**Load permanently (Linux):**

```bash
# System-wide
haft completion bash > /etc/bash_completion.d/haft

# User-only
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
# Ensure completion is enabled in ~/.zshrc:
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
# Add `. /path/to/haft.ps1` to your profile
```

## Updating

To update to the latest version:

```bash
# Using the install script (recommended)
curl -fsSL https://raw.githubusercontent.com/KashifKhn/haft/main/install.sh | bash

# If installed via Go
go install github.com/KashifKhn/haft/cmd/haft@latest

# If installed from source
cd haft
git pull
make build
```

## Uninstalling

```bash
# Remove the binary
sudo rm /usr/local/bin/haft

# Or from ~/.local/bin
rm ~/.local/bin/haft

# Or if installed via Go
rm $(which haft)
```
