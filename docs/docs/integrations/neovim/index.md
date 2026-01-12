---
sidebar_position: 1
title: Neovim Plugin
description: Seamless Haft CLI integration for Neovim
---

# haft.nvim

Neovim plugin for [Haft CLI](https://github.com/KashifKhn/haft) - bringing Spring Boot productivity tools directly into your editor with async commands, interactive Telescope pickers, and beautiful floating windows.

## Features

- ✅ **Full CLI Integration** - All Haft commands available as `:Haft*` commands
- ✅ **Async Execution** - Non-blocking operations using plenary.nvim
- ✅ **Telescope Pickers** - Interactive dependency selection and browsing
- ✅ **Floating Windows** - Beautiful output for info, routes, and stats
- ✅ **Auto-Open Files** - Generated files automatically opened in buffers
- ✅ **Quickfix Integration** - Navigate generated files with quickfix list
- ✅ **Project Detection** - Automatic detection of Haft/Spring Boot projects
- ✅ **Fully Configurable** - Customize every aspect, no default keybindings
- ✅ **Graceful Fallbacks** - Works with or without Telescope

## Quick Start

### Install with lazy.nvim

```lua
{
  "KashifKhn/haft.nvim",
  dependencies = {
    "nvim-lua/plenary.nvim",
    "nvim-telescope/telescope.nvim", -- optional but recommended
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

### Basic Usage

```vim
" Generate a complete CRUD resource
:HaftGenerateResource User

" Add dependencies with interactive picker
:HaftAdd

" View project information in floating window
:HaftInfo

" Show API routes
:HaftRoutes
```

## What Makes It Different?

haft.nvim is a **thin wrapper** around the Haft CLI, adding:

| Feature | CLI | haft.nvim |
|---------|-----|-----------|
| Execution | Blocking terminal | ✅ Async (non-blocking) |
| Dependency Selection | TUI picker | ✅ Telescope picker |
| Output Display | Terminal output | ✅ Floating windows |
| File Opening | Manual | ✅ Auto-open in buffers |
| Navigation | - | ✅ Quickfix integration |

> **Note:** All functionality (generation logic, detection, templates) comes from the CLI. The plugin provides convenience and enhanced UX. See [CLI Commands](/docs/commands/generate) to understand what each command does.

## Requirements

### Required

- **Neovim** >= 0.9.0
- **[Haft CLI](https://github.com/KashifKhn/haft)** >= 0.1.11
- **[plenary.nvim](https://github.com/nvim-lua/plenary.nvim)** - Async utilities

### Optional (Recommended)

- **[telescope.nvim](https://github.com/nvim-telescope/telescope.nvim)** - Interactive pickers (falls back to vim.ui.select)
- **[noice.nvim](https://github.com/folke/noice.nvim)** - Enhanced notifications

## Getting Started

1. **[Install the plugin →](./installation)** - Setup with lazy.nvim, packer, or vim-plug
2. **[Configure options →](./configuration)** - Customize behavior and UI
3. **[Learn commands →](./usage)** - Command reference and workflows
4. **[Troubleshoot →](./troubleshooting)** - Health checks and common issues

## Example Workflow

Generate a new feature in seconds:

```vim
" 1. Generate CRUD resource
:HaftGenerateResource Product

" Output:
" ✓ Generated 8 files (ProductController, ProductService, etc.)
" ✓ Opened ProductController.java
" ✓ Added files to quickfix list

" 2. Navigate generated files
:cnext    " ProductService.java
:cnext    " ProductRepository.java

" 3. Add dependencies if needed
:HaftAdd
" Select 'validation' and 'mapstruct' from Telescope picker
```

## Resources

- **GitHub Repository**: [KashifKhn/haft.nvim](https://github.com/KashifKhn/haft.nvim)
- **Vim Help**: `:help haft` (after installation)
- **Health Check**: `:checkhealth haft`
- **Report Issues**: [GitHub Issues](https://github.com/KashifKhn/haft.nvim/issues)

## Next Steps

- [Installation Guide →](./installation)
- [Configuration Reference →](./configuration)
- [Usage & Commands →](./usage)
- [Troubleshooting →](./troubleshooting)
