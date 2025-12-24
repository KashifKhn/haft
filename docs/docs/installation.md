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

## Using Go

The easiest way to install Haft:

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
sudo mv haft /usr/local/bin/

# ARM64
curl -L https://github.com/KashifKhn/haft/releases/latest/download/haft-linux-arm64.tar.gz | tar xz
sudo mv haft /usr/local/bin/
```

### macOS

```bash
# Intel Mac
curl -L https://github.com/KashifKhn/haft/releases/latest/download/haft-darwin-amd64.tar.gz | tar xz
sudo mv haft /usr/local/bin/

# Apple Silicon
curl -L https://github.com/KashifKhn/haft/releases/latest/download/haft-darwin-arm64.tar.gz | tar xz
sudo mv haft /usr/local/bin/
```

### Windows

1. Download `haft-windows-amd64.zip` from the [releases page](https://github.com/KashifKhn/haft/releases)
2. Extract the ZIP file
3. Add the extracted folder to your system PATH

Or using PowerShell:

```powershell
# Download and extract
Invoke-WebRequest -Uri "https://github.com/KashifKhn/haft/releases/latest/download/haft-windows-amd64.zip" -OutFile "haft.zip"
Expand-Archive -Path "haft.zip" -DestinationPath "C:\Program Files\haft"

# Add to PATH (run as Administrator)
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\haft", "Machine")
```

## Verify Installation

```bash
haft version
```

You should see output like:

```
haft version 0.1.0
```

## Shell Completions

Haft supports shell completions for bash, zsh, fish, and PowerShell.

### Bash

```bash
# Add to ~/.bashrc
source <(haft completion bash)
```

### Zsh

```bash
# Add to ~/.zshrc
source <(haft completion zsh)
```

### Fish

```bash
haft completion fish | source
```

### PowerShell

```powershell
haft completion powershell | Out-String | Invoke-Expression
```

## Updating

To update to the latest version:

```bash
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

# Or if installed via Go
rm $(which haft)
```
