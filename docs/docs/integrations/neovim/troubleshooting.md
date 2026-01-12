---
sidebar_position: 5
title: Troubleshooting
description: Common issues and solutions for haft.nvim
---

# Troubleshooting

Common issues, solutions, and debugging tips for haft.nvim.

## Health Check

Always start with the health check to diagnose issues:

```vim
:checkhealth haft
```

Expected output when everything works:

```
haft: require("haft.health").check()

haft.nvim
- OK Neovim >= 0.9.0
- OK Haft CLI found: haft version v0.1.11
- OK plenary.nvim installed
- OK telescope.nvim installed (optional)
- INFO Haft project detected in current directory
```

---

## Common Issues

### 1. Haft CLI Not Found

**Symptom:**

```
[haft.nvim] ERROR: Haft CLI not found in PATH
```

Or health check shows:

```
- ERROR Haft CLI not found
```

**Solution:**

#### Check if CLI is installed:

```bash
haft --version
```

If not found, install the CLI first:

```bash
# See installation guide
curl -fsSL https://raw.githubusercontent.com/KashifKhn/haft/main/install.sh | bash
```

[CLI Installation Guide →](/docs/installation)

#### Verify CLI is in PATH:

```bash
which haft
# Should show: /usr/local/bin/haft (or similar)

echo $PATH
# Should include directory containing haft
```

#### Configure custom path:

If CLI is installed but not in PATH, specify the full path:

```lua
require("haft").setup({
  haft_path = "/usr/local/bin/haft",
  -- Or:
  haft_path = vim.fn.expand("~/bin/haft"),
})
```

---

### 2. Commands Not Available

**Symptom:**

- `:Haft<Tab>` shows nothing
- Commands like `:HaftInfo` show `E492: Not an editor command`

**Solution:**

#### Verify plugin is loaded:

```vim
:lua print(vim.inspect(package.loaded["haft"]))
```

If output is `nil`, the plugin didn't load.

#### Check plugin manager:

**lazy.nvim:**
```vim
:Lazy
" Check if haft.nvim is listed and loaded
```

**packer.nvim:**
```vim
:PackerStatus
" Check if haft.nvim is installed

:PackerCompile
:PackerSync
```

**vim-plug:**
```vim
:PlugStatus
```

#### Check for errors:

```vim
:messages
" Look for error messages related to haft
```

#### Reload plugin:

```vim
" lazy.nvim
:Lazy reload haft.nvim

" packer.nvim
:PackerCompile
:PackerLoad haft.nvim

" Manual
:lua package.loaded["haft"] = nil
:lua require("haft").setup()
```

---

### 3. Telescope Not Working

**Symptom:**

```
[haft.nvim] ERROR: Telescope not found but provider is set to 'telescope'
```

Or pickers don't open.

**Solution:**

#### Option 1: Install Telescope

Add to your plugin manager:

```lua
-- lazy.nvim
{
  "nvim-telescope/telescope.nvim",
  dependencies = { "nvim-lua/plenary.nvim" },
}
```

Then sync:

```vim
:Lazy sync
```

#### Option 2: Use native picker

```lua
require("haft").setup({
  picker = { provider = "native" },
})
```

This uses `vim.ui.select` instead of Telescope.

#### Option 3: Use auto mode

```lua
require("haft").setup({
  picker = { provider = "auto" },
})
```

Auto mode tries Telescope first, falls back to native.

---

### 4. Project Not Detected

**Symptom:**

- `:HaftInfo` shows error
- Health check shows `INFO No Haft project in current directory`

**Solution:**

#### Verify you're in a Spring Boot project:

```bash
ls
# Should see: pom.xml OR build.gradle OR .haft.yaml
```

#### Check detection patterns:

```lua
require("haft").setup({
  detection = {
    enabled = true,
    patterns = { ".haft.yaml", "pom.xml", "build.gradle", "build.gradle.kts" },
  },
})
```

#### Navigate to project root:

```vim
:cd /path/to/your/spring-boot-project
:HaftInfo
```

#### Verify with CLI:

```bash
cd /path/to/project
haft info
```

If CLI works but plugin doesn't, there's a path or detection issue.

---

### 5. Floating Windows Not Showing

**Symptom:**

- `:HaftInfo` runs but no window appears
- No error messages

**Solution:**

#### Check terminal size:

Floating windows require minimum terminal size. Try increasing window size.

#### Check border setting:

Some terminals don't support certain borders:

```lua
require("haft").setup({
  float = { border = "single" },  -- Try different border
})
```

#### Check for conflicting plugins:

Disable other float-heavy plugins temporarily to test.

#### Test with minimal config:

```lua
-- Minimal test
require("haft").setup({
  float = {
    border = "single",
    width = 0.5,
    height = 0.5,
  },
})
```

---

### 6. Generated Files Not Opening

**Symptom:**

- `:HaftGenerateResource User` succeeds
- Files are created but don't open

**Solution:**

#### Check auto_open setting:

```vim
:lua print(vim.inspect(require("haft.config").get().auto_open))
```

Should show:
```lua
{
  enabled = true,
  strategy = "first"
}
```

If `enabled = false`, auto-open is disabled:

```lua
require("haft").setup({
  auto_open = { enabled = true },
})
```

#### Manual file opening:

If auto-open doesn't work, use quickfix:

```vim
:HaftGenerateResource User
:copen    " Open quickfix
:cnext    " Go to first file
```

---

### 7. Quickfix Not Working

**Symptom:**

- `:copen` after generation shows empty list

**Solution:**

#### Check quickfix setting:

```lua
require("haft").setup({
  quickfix = { enabled = true },
})
```

#### Verify generation succeeded:

```vim
:messages
" Should show success message with file count
```

#### Check manually:

```vim
:cexpr system('haft generate controller Test --json')
:copen
```

---

### 8. Performance Issues / Slow Commands

**Symptom:**

- Commands take a long time to complete
- Editor feels sluggish

**Solution:**

#### Check profile caching:

Profile should be cached to `.haft/profile.yaml`. Check if it exists:

```bash
ls -la .haft/
# Should see: profile.yaml
```

If missing, first run will be slow (scanning project). Subsequent runs should be fast.

#### Disable refresh if enabled:

```lua
require("haft").setup({
  commands = { generate = { refresh = false } },  -- Use cache
})
```

#### Check CLI performance:

```bash
time haft info
# Should complete in < 1 second for cached profile
```

If CLI is slow, issue is with CLI not plugin.

---

### 9. Notifications Not Showing

**Symptom:**

- Commands run but no notifications appear

**Solution:**

#### Check notification setting:

```lua
require("haft").setup({
  notifications = { enabled = true },
})
```

#### Check notification level:

```lua
require("haft").setup({
  notifications = { level = "info" },  -- Not "error" only
})
```

#### Test notification:

```vim
:lua require("haft.ui.notify").info("Test notification")
```

If this works, notifications are fine. If not, check `vim.notify` setup.

---

### 10. JSON Parsing Errors

**Symptom:**

```
[haft.nvim] ERROR: Failed to parse JSON output
```

**Solution:**

#### Verify CLI JSON output:

```bash
haft info --json
# Should output valid JSON
```

If output is not JSON, CLI version is too old or broken.

#### Update CLI:

```bash
haft upgrade
# Or reinstall
curl -fsSL https://raw.githubusercontent.com/KashifKhn/haft/main/install.sh | bash
```

#### Check CLI version:

```bash
haft --version
# Should be >= 0.1.11
```

---

## Debugging

### Enable Debug Logging

Get detailed logs for troubleshooting:

```lua
require("haft").setup({
  notifications = { level = "debug" },
})
```

Then run commands and check `:messages`.

### Manual Testing

Test CLI integration manually:

```vim
" Test CLI execution
:lua vim.print(vim.fn.system("haft --version"))

" Test project detection
:lua vim.print(require("haft.detection").get_project_root())

" Test config
:lua vim.print(require("haft.config").get())
```

### Check Logs

View all messages:

```vim
:messages
```

Clear and re-run:

```vim
:messages clear
:HaftInfo
:messages
```

---

## Known Issues

### Issue 1: Telescope Preview Empty

**Status:** Known limitation

**Workaround:** Use arrow keys to navigate, preview updates on selection change.

### Issue 2: Commands Slow on First Run

**Status:** Expected behavior (profile caching)

**Explanation:** First run scans project and caches to `.haft/profile.yaml`. Subsequent runs are instant.

---

## Getting Help

If issues persist:

### 1. Gather Information

```vim
" Run health check
:checkhealth haft

" Check versions
:version
:lua print(vim.inspect(vim.version()))
!haft --version

" Check config
:lua vim.print(require("haft.config").get())

" Check messages
:messages
```

### 2. Minimal Reproduction

Test with minimal config:

```lua
-- minimal_init.lua
vim.opt.rtp:append("~/.local/share/nvim/site/pack/plugins/start/plenary.nvim")
vim.opt.rtp:append("~/.local/share/nvim/site/pack/plugins/start/haft.nvim")

require("haft").setup()
```

Run:

```bash
nvim -u minimal_init.lua
```

### 3. Report Issue

If problem persists, create an issue with:

- `:checkhealth haft` output
- Neovim version (`:version`)
- CLI version (`haft --version`)
- Plugin config
- Steps to reproduce
- Error messages from `:messages`

**Links:**

- [haft.nvim Issues](https://github.com/KashifKhn/haft.nvim/issues)
- [Haft CLI Issues](https://github.com/KashifKhn/haft/issues)
- [Discussions](https://github.com/KashifKhn/haft/discussions)

---

## Related

- [Configuration →](./configuration) - Adjust settings
- [Usage Guide →](./usage) - Learn commands
- [CLI Installation](/docs/installation) - CLI installation and setup
