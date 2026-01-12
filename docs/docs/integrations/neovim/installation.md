---
sidebar_position: 2
title: Installation
description: Install haft.nvim with your preferred plugin manager
---

# Installation

haft.nvim can be installed with any Neovim plugin manager. Choose your preferred method below.

## Prerequisites

Before installing haft.nvim, ensure you have:

1. **Neovim >= 0.9.0**
   ```bash
   nvim --version
   # Should show v0.9.0 or higher
   ```

2. **Haft CLI >= 0.1.11** installed and in PATH
   ```bash
   haft --version
   # Install if needed: https://haft.kashifkhan.dev/docs/installation
   ```

3. **plenary.nvim** (required dependency)
   - Usually installed automatically by your plugin manager

4. **telescope.nvim** (optional but recommended)
   - Provides interactive pickers
   - Falls back to `vim.ui.select` if not available

---

## Using lazy.nvim (Recommended)

[lazy.nvim](https://github.com/folke/lazy.nvim) is the modern Neovim plugin manager with lazy loading and async execution.

### Basic Setup

```lua
{
  "KashifKhn/haft.nvim",
  dependencies = {
    "nvim-lua/plenary.nvim",
    "nvim-telescope/telescope.nvim", -- optional
  },
  cmd = {
    "HaftInfo", "HaftRoutes", "HaftStats",
    "HaftAdd", "HaftRemove",
    "HaftGenerateResource", "HaftGenerateController",
    "HaftGenerateService", "HaftGenerateRepository",
    "HaftGenerateEntity", "HaftGenerateDto",
  },
  opts = {},
}
```

### With Custom Configuration

```lua
{
  "KashifKhn/haft.nvim",
  dependencies = {
    "nvim-lua/plenary.nvim",
    "nvim-telescope/telescope.nvim",
  },
  cmd = {
    "HaftInfo", "HaftAdd", "HaftGenerateResource",
    -- Add other commands you use frequently
  },
  config = function()
    require("haft").setup({
      float = { border = "rounded", width = 0.9, height = 0.9 },
      picker = { provider = "telescope" },
      auto_open = { enabled = true, strategy = "first" },
    })
  end,
}
```

### Eager Loading (Always Active)

```lua
{
  "KashifKhn/haft.nvim",
  dependencies = { "nvim-lua/plenary.nvim" },
  lazy = false,  -- Load on startup
  config = function()
    require("haft").setup()
  end,
}
```

---

## Using packer.nvim

[packer.nvim](https://github.com/wbthomason/packer.nvim) is a popular plugin manager with good performance.

### Basic Setup

```lua
use {
  "KashifKhn/haft.nvim",
  requires = {
    "nvim-lua/plenary.nvim",
    "nvim-telescope/telescope.nvim",
  },
  config = function()
    require("haft").setup()
  end,
}
```

### With Command Lazy Loading

```lua
use {
  "KashifKhn/haft.nvim",
  requires = {
    "nvim-lua/plenary.nvim",
    "nvim-telescope/telescope.nvim",
  },
  cmd = { "HaftInfo", "HaftAdd", "HaftGenerateResource" },
  config = function()
    require("haft").setup({
      -- Your config here
    })
  end,
}
```

After adding to your `packer` config:

```vim
:PackerSync
```

---

## Using vim-plug

[vim-plug](https://github.com/junegunn/vim-plug) is a minimalist plugin manager.

### Basic Setup

Add to your `init.vim` or `init.lua`:

```vim
" init.vim
Plug 'nvim-lua/plenary.nvim'
Plug 'nvim-telescope/telescope.nvim'
Plug 'KashifKhn/haft.nvim'

" After plug#end(), configure:
lua << EOF
require("haft").setup()
EOF
```

Or in `init.lua`:

```lua
-- init.lua
vim.cmd [[
  call plug#begin()
  Plug 'nvim-lua/plenary.nvim'
  Plug 'nvim-telescope/telescope.nvim'
  Plug 'KashifKhn/haft.nvim'
  call plug#end()
]]

require("haft").setup()
```

Then install:

```vim
:PlugInstall
```

---

## Using rocks.nvim

[rocks.nvim](https://github.com/nvim-neorocks/rocks.nvim) is the new LuaRocks-based plugin manager.

```vim
:Rocks install haft.nvim
```

Note: rocks.nvim automatically handles dependencies (plenary.nvim).

Then configure in your `init.lua`:

```lua
require("haft").setup()
```

---

## Manual Installation

For development or testing:

```bash
# Clone the repository
git clone https://github.com/KashifKhn/haft.nvim.git ~/.local/share/nvim/site/pack/plugins/start/haft.nvim

# Also install plenary.nvim if not already installed
git clone https://github.com/nvim-lua/plenary.nvim.git ~/.local/share/nvim/site/pack/plugins/start/plenary.nvim
```

Then in your `init.lua`:

```lua
require("haft").setup()
```

---

## Verification

After installation, verify everything works:

### 1. Check Health

```vim
:checkhealth haft
```

Expected output:
```
haft: require("haft.health").check()

haft.nvim
- OK Neovim >= 0.9.0
- OK Haft CLI found: haft version v0.1.11
- OK plenary.nvim installed
- OK telescope.nvim installed (optional)
```

### 2. Test a Command

In any Spring Boot project:

```vim
:HaftInfo
```

You should see a floating window with project information.

### 3. Check Available Commands

```vim
:Haft<Tab>
```

Should show autocomplete with all available commands.

---

## Updating

### lazy.nvim
```vim
:Lazy sync
```

### packer.nvim
```vim
:PackerUpdate
```

### vim-plug
```vim
:PlugUpdate
```

### rocks.nvim
```vim
:Rocks update
```

---

## Troubleshooting Installation

### Haft CLI Not Found

**Error:**
```
[haft.nvim] Haft CLI not found in PATH
```

**Solution:**

1. Verify CLI is installed:
   ```bash
   haft --version
   ```

2. If not installed, see [CLI Installation Guide](/docs/installation)

3. If installed but not in PATH, configure custom path:
   ```lua
   require("haft").setup({
     haft_path = "/usr/local/bin/haft",
   })
   ```

### Telescope Not Working

**Error:**
```
[haft.nvim] Telescope not found
```

**Solution:**

Option 1: Install telescope.nvim
```lua
-- lazy.nvim
dependencies = { "nvim-telescope/telescope.nvim" }
```

Option 2: Use native picker
```lua
require("haft").setup({
  picker = { provider = "native" },
})
```

### Commands Not Available

**Symptom:** `:Haft<Tab>` shows nothing

**Solution:**

1. Ensure plugin is loaded:
   ```vim
   :lua print(vim.inspect(package.loaded["haft"]))
   ```

2. If `nil`, the plugin didn't load. Check:
   - Plugin manager loaded the plugin
   - No errors in `:messages`
   - Try `:PackerCompile` (packer) or `:Lazy reload haft.nvim` (lazy)

---

## Next Steps

- [Configuration →](./configuration) - Customize haft.nvim
- [Usage Guide →](./usage) - Learn commands and workflows
- [Troubleshooting →](./troubleshooting) - Common issues and solutions
