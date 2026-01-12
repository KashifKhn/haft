---
sidebar_position: 11
---

# upgrade

Upgrade Haft CLI to the latest version.

## Usage

```bash
haft upgrade [flags]
```

## Description

The `haft upgrade` command provides a safe and convenient way to update your Haft installation to the latest version. It automatically:

1. Checks if a newer version is available
2. Downloads the latest release for your platform
3. Verifies download integrity using SHA256 checksums
4. Creates a backup of your current installation
5. Installs the new version
6. Verifies the installation works correctly
7. Automatically rolls back if anything fails

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--check` | `-c` | Only check for updates without installing |
| `--force` | `-f` | Force upgrade even if already on latest version |
| `--version` | `-v` | Upgrade to a specific version |
| `--json` | | Output result as JSON |

## Examples

### Check and upgrade to latest version

```bash
haft upgrade
```

Output:
```
INFO  Checking for updates...
INFO  Current version version=v0.4.2
INFO  Latest version version=v0.5.0
INFO  Platform os=linux arch=amd64
INFO  Downloading version=v0.5.0
INFO  Extracting binary...
INFO  Installing new version...
INFO  Verifying installation...
SUCCESS  Successfully upgraded from v0.4.2 to v0.5.0
```

### Check for updates without installing

```bash
haft upgrade --check
```

Output:
```
INFO  Checking for updates...
INFO  Current version version=v0.4.2
INFO  Latest version version=v0.5.0
SUCCESS  Update available! Run 'haft upgrade' to install.
```

### Force reinstall

```bash
haft upgrade --force
```

Useful when you want to reinstall the current version (e.g., if the binary is corrupted).

### Upgrade to specific version

```bash
haft upgrade --version v0.4.0
```

### JSON output for scripting

```bash
haft upgrade --check --json
```

Output:
```json
{
  "current_version": "v0.4.2",
  "latest_version": "v0.5.0",
  "update_available": true,
  "upgraded": false,
  "platform": {
    "os": "linux",
    "arch": "amd64"
  }
}
```

## Platform Support

The upgrade command supports the following platforms:

| Platform | Architecture |
|----------|-------------|
| Linux | amd64, arm64 |
| macOS | amd64 (Intel), arm64 (Apple Silicon) |
| Windows | amd64 |

## Safety Features

### Automatic Backup

Before upgrading, Haft creates a backup of your current installation. If anything goes wrong during the upgrade, your original version is automatically restored.

### Checksum Verification

Downloaded files are verified against SHA256 checksums published with each release to ensure download integrity.

### Installation Verification

After installation, Haft runs a verification check to ensure the new binary works correctly. If verification fails, the backup is automatically restored.

### Multiple Installation Support

If Haft is installed in multiple locations (e.g., `/usr/local/bin` and `~/.local/bin`), the upgrade command will update all installations.

## Error Handling

| Scenario | Behavior |
|----------|----------|
| No internet connection | Shows error message |
| Already on latest version | Shows "already up to date" message |
| Download fails | Shows error, no changes made |
| Checksum mismatch | Shows error, no changes made |
| Installation fails | Automatically restores backup |
| Verification fails | Automatically restores backup |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success (upgraded or already up to date) |
| 1 | Error occurred |

## Editor Integration

Use this command from your editor:

- **Neovim**: Not yet available in haft.nvim (CLI only)
- **VS Code**: Coming soon ([preview →](/docs/integrations/vscode))
- **IntelliJ IDEA**: Coming soon ([preview →](/docs/integrations/intellij))

## See Also

- [Installation](../installation.md) - Initial installation instructions
