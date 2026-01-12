---
sidebar_position: 4
title: Usage
description: Commands, workflows, and examples for haft.nvim
---

# Usage Guide

Learn how to use haft.nvim commands, Telescope integration, Lua API, and common workflows.

## Command Reference

All Haft CLI commands are available in Neovim with the `:Haft` prefix. Each command runs asynchronously and provides enhanced UI.

### Project Information Commands

| Neovim Command | CLI Equivalent | Description |
|----------------|----------------|-------------|
| `:HaftInfo` | `haft info` | Show project info in floating window ([CLI docs →](/docs/commands/info)) |
| `:HaftRoutes` | `haft routes` | Show API routes in floating window ([CLI docs →](/docs/commands/routes)) |
| `:HaftStats` | `haft stats` | Show code statistics in floating window ([CLI docs →](/docs/commands/stats)) |

**Example:**

```vim
:HaftInfo
" Output: Floating window with project details

:HaftRoutes
" Output: List of all @RestController endpoints

:HaftStats
" Output: Code statistics with COCOMO estimates
```

---

### Dependency Management Commands

| Neovim Command | CLI Equivalent | Description |
|----------------|----------------|-------------|
| `:HaftAdd [dep...]` | `haft add [dep...]` | Add dependencies ([CLI docs →](/docs/commands/add)) |
| `:HaftRemove [dep...]` | `haft remove [dep...]` | Remove dependencies ([CLI docs →](/docs/commands/remove)) |

**Interactive Mode** (no arguments):

```vim
:HaftAdd
" Opens Telescope picker with all available dependencies
" Press <Tab> to multi-select, <Enter> to confirm

:HaftRemove
" Opens Telescope picker with current project dependencies
" Select dependencies to remove
```

**Direct Mode** (with arguments):

```vim
:HaftAdd lombok validation
" Adds lombok and spring-boot-starter-validation

:HaftRemove lombok
" Removes lombok dependency
```

---

### Code Generation Commands

| Neovim Command | CLI Equivalent | Description |
|----------------|----------------|-------------|
| `:HaftGenerateResource [name]` | `haft generate resource [name]` | Generate complete CRUD resource ([CLI docs →](/docs/commands/generate#haft-generate-resource)) |
| `:HaftGenerateController [name]` | `haft generate controller [name]` | Generate REST controller ([CLI docs →](/docs/commands/generate#haft-generate-controller)) |
| `:HaftGenerateService [name]` | `haft generate service [name]` | Generate service layer ([CLI docs →](/docs/commands/generate#haft-generate-service)) |
| `:HaftGenerateRepository [name]` | `haft generate repository [name]` | Generate JPA repository ([CLI docs →](/docs/commands/generate#haft-generate-repository)) |
| `:HaftGenerateEntity [name]` | `haft generate entity [name]` | Generate JPA entity ([CLI docs →](/docs/commands/generate#haft-generate-entity)) |
| `:HaftGenerateDto [name]` | `haft generate dto [name]` | Generate Request/Response DTOs ([CLI docs →](/docs/commands/generate#haft-generate-dto)) |

**With Name Argument:**

```vim
:HaftGenerateResource Product
" Generates: ProductController, ProductService, ProductRepository,
"            Product (entity), ProductRequest, ProductResponse, etc.
" Auto-opens ProductController.java
```

**Without Name Argument:**

```vim
:HaftGenerateController
" Prompts: "Enter controller name:"
" Type: Order
" Generates: OrderController.java
```

> **Note:** All generation respects your project's architecture and conventions. Haft automatically adapts to your codebase structure.

---

## Telescope Integration

haft.nvim provides Telescope pickers for enhanced dependency management.

### Loading the Extension

```lua
require("telescope").load_extension("haft")
```

### Available Pickers

#### Dependencies Picker

Interactive picker for adding dependencies:

```vim
:Telescope haft dependencies

" Or via Lua:
:lua require("telescope").extensions.haft.dependencies()
```

**Features:**
- Search all Spring Boot dependencies
- Browse by category
- Multi-select with `<Tab>`
- Preview pane shows dependency details
- Press `/` to filter

#### Remove Dependencies Picker

Interactive picker for removing dependencies:

```vim
:Telescope haft remove

" Or via Lua:
:lua require("telescope").extensions.haft.remove()
```

**Features:**
- Shows current project dependencies
- Multi-select for batch removal
- Shows version and scope information

### Telescope Keybindings

Default Telescope keybindings in pickers:

| Key | Action |
|-----|--------|
| `<Enter>` | Select (single) or confirm (multi) |
| `<Tab>` | Toggle selection (multi-select mode) |
| `<C-x>` | Open in horizontal split |
| `<C-v>` | Open in vertical split |
| `/` | Search/filter |
| `<Esc>` | Close picker |
| `<C-c>` | Close picker |

---

## Lua API

Use haft.nvim programmatically from Lua.

### Setup

```lua
require("haft").setup(config)
```

### Project Detection

```lua
local haft = require("haft")

-- Check if current directory is a Haft project
if haft.is_haft_project() then
  print("Haft project detected")
end

-- Get project information
local info = haft.get_project_info()
if info then
  print("Project:", info.name)
  print("Type:", info.type)
end
```

### Commands (Programmatic)

```lua
local haft = require("haft")

-- Project information
haft.info()
haft.routes()
haft.stats()

-- Dependency management
haft.add()                -- Opens picker
haft.add({ "lombok", "validation" })  -- Direct add
haft.remove()             -- Opens picker
haft.remove({ "lombok" }) -- Direct remove

-- Code generation
haft.generate_resource("User")
haft.generate_controller("Order")
haft.generate_service("Payment")
haft.generate_repository("Product")
haft.generate_entity("Customer")
haft.generate_dto("Invoice")
```

---

## Keybindings

haft.nvim provides **no default keybindings**. Configure your own based on your workflow.

### Suggested Keybindings

#### Leader-Based (General Use)

```lua
-- Project info
vim.keymap.set("n", "<leader>hi", "<cmd>HaftInfo<cr>", { desc = "Haft: Info" })
vim.keymap.set("n", "<leader>hr", "<cmd>HaftRoutes<cr>", { desc = "Haft: Routes" })
vim.keymap.set("n", "<leader>hs", "<cmd>HaftStats<cr>", { desc = "Haft: Stats" })

-- Dependencies
vim.keymap.set("n", "<leader>ha", "<cmd>HaftAdd<cr>", { desc = "Haft: Add dependency" })
vim.keymap.set("n", "<leader>hR", "<cmd>HaftRemove<cr>", { desc = "Haft: Remove dependency" })

-- Generation
vim.keymap.set("n", "<leader>hg", "<cmd>HaftGenerateResource<cr>", { desc = "Haft: Generate resource" })
vim.keymap.set("n", "<leader>hc", "<cmd>HaftGenerateController<cr>", { desc = "Haft: Generate controller" })
vim.keymap.set("n", "<leader>he", "<cmd>HaftGenerateEntity<cr>", { desc = "Haft: Generate entity" })
```

#### LocalLeader (Java/Kotlin Files Only)

```lua
vim.api.nvim_create_autocmd("FileType", {
  pattern = { "java", "kotlin" },
  callback = function()
    local opts = { buffer = true }
    vim.keymap.set("n", "<localleader>i", "<cmd>HaftInfo<cr>", vim.tbl_extend("force", opts, { desc = "Project info" }))
    vim.keymap.set("n", "<localleader>r", "<cmd>HaftRoutes<cr>", vim.tbl_extend("force", opts, { desc = "API routes" }))
    vim.keymap.set("n", "<localleader>a", "<cmd>HaftAdd<cr>", vim.tbl_extend("force", opts, { desc = "Add dependency" }))
    vim.keymap.set("n", "<localleader>g", "<cmd>HaftGenerateResource<cr>", vim.tbl_extend("force", opts, { desc = "Generate resource" }))
  end,
})
```

#### With which-key.nvim

```lua
local wk = require("which-key")
wk.add({
  { "<leader>h", group = "Haft" },
  { "<leader>hi", "<cmd>HaftInfo<cr>", desc = "Info" },
  { "<leader>hr", "<cmd>HaftRoutes<cr>", desc = "Routes" },
  { "<leader>hs", "<cmd>HaftStats<cr>", desc = "Stats" },
  { "<leader>ha", "<cmd>HaftAdd<cr>", desc = "Add dependency" },
  { "<leader>hR", "<cmd>HaftRemove<cr>", desc = "Remove dependency" },
  { "<leader>hg", "<cmd>HaftGenerateResource<cr>", desc = "Generate resource" },
})
```

---

## Common Workflows

### Workflow 1: Generate New Feature

```vim
" 1. Generate complete CRUD resource
:HaftGenerateResource Product

" Output:
" ✓ Generated 8 files
" ✓ Opened ProductController.java

" 2. Navigate generated files via quickfix
:cnext    " ProductService.java
:cnext    " ProductServiceImpl.java
:cnext    " ProductRepository.java

" 3. Review all files
:copen    " Open quickfix window to see all files
```

### Workflow 2: Add Dependencies

```vim
" 1. Open dependency picker
:HaftAdd

" 2. In Telescope picker:
"    - Type 'valid' to search
"    - Press <Tab> on 'spring-boot-starter-validation'
"    - Type 'mapstr'
"    - Press <Tab> on 'mapstruct'
"    - Press <Enter> to confirm

" 3. Dependencies added to pom.xml/build.gradle
" ✓ Added 2 dependencies
```

### Workflow 3: Analyze Project

```vim
" View project metadata
:HaftInfo
" Shows: group, artifact, Java version, Spring Boot version, dependencies count

" View API endpoints
:HaftRoutes
" Shows: HTTP method, path, controller method

" View code statistics
:HaftStats
" Shows: LOC, files, packages, COCOMO estimates
```

### Workflow 4: Refactor with Context

```vim
" 1. Check existing routes
:HaftRoutes

" 2. See which controllers exist
" Output shows UserController, OrderController

" 3. Generate new matching controller
:HaftGenerateController Payment

" 4. Review and edit
:e PaymentController.java
```

---

## File Opening Behavior

After generation commands, haft.nvim can automatically open files based on your configuration.

### Strategy: "first" (Default)

```vim
:HaftGenerateResource User
" Opens: UserController.java (first file)
" Quickfix: All 8 files added
```

### Strategy: "all"

```lua
require("haft").setup({
  auto_open = { strategy = "all" },
})
```

```vim
:HaftGenerateResource User
" Opens: All 8 files as buffers
" Use: :bnext, :bprev to navigate
```

### Strategy: "none"

```lua
require("haft").setup({
  auto_open = { enabled = false },
})
```

```vim
:HaftGenerateResource User
" Opens: Nothing (files added to quickfix only)
" Use: :copen, :cnext to navigate
```

---

## Integration with Other Plugins

### nvim-tree / neo-tree

Generated files appear in file tree automatically. Refresh with `R` if needed.

### trouble.nvim

Quickfix entries from haft.nvim appear in Trouble:

```vim
:Trouble quickfix
```

### harpoon

Mark frequently generated file locations:

```lua
-- After generating User resource
require("harpoon.mark").add_file()  -- Mark UserController.java
```

---

## Tips & Tricks

### 1. Quick Generate from Visual Selection

```lua
-- Map visual selection to generate
vim.keymap.set("v", "<leader>hg", function()
  local text = vim.fn.getreg('"')
  vim.cmd("HaftGenerateResource " .. text)
end, { desc = "Generate resource from selection" })

-- Usage:
-- 1. Visual select "Product"
-- 2. Press <leader>hg
-- 3. Generates ProductController, etc.
```

### 2. Auto-Refresh Profile

If your project structure changes frequently:

```lua
require("haft").setup({
  commands = { generate = { refresh = true } },
})
```

### 3. Combine with LSP

After generating files, format with LSP:

```vim
:HaftGenerateResource Order
:cnext | lua vim.lsp.buf.format()  " Format OrderService
:cnext | lua vim.lsp.buf.format()  " Format OrderRepository
```

### 4. Custom Command Wrapper

Create convenience commands:

```vim
" Generate and format all
command! -nargs=1 HaftGen execute 'HaftGenerateResource ' .. <q-args> | cfdo lua vim.lsp.buf.format()

" Usage:
:HaftGen Product
```

---

## Related

- [CLI Commands Reference](/docs/commands/generate) - Understand what each command does
- [Configuration →](./configuration) - Customize behavior
- [Troubleshooting →](./troubleshooting) - Common issues
