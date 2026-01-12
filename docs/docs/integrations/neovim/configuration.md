---
sidebar_position: 3
title: Configuration
description: Complete configuration reference for haft.nvim
---

# Configuration

haft.nvim is fully configurable with sensible defaults. You only need to specify options you want to change from the defaults.

## Setup Function

```lua
require("haft").setup(config)
```

Pass a table with your configuration options. Any omitted options use the defaults.

---

## Default Configuration

Here's the complete default configuration:

```lua
require("haft").setup({
  -- CLI Binary Path
  haft_path = "haft",  -- Auto-detected if in PATH

  -- Project Detection
  detection = {
    enabled = true,
    patterns = {
      ".haft.yaml",
      "pom.xml",
      "build.gradle",
      "build.gradle.kts",
    },
  },

  -- Notifications
  notifications = {
    enabled = true,
    level = "info",     -- "debug" | "info" | "warn" | "error"
    timeout = 3000,     -- milliseconds
  },

  -- Floating Windows
  float = {
    border = "rounded", -- "none" | "single" | "double" | "rounded" | "solid" | "shadow"
    width = 0.8,        -- 0-1 for percentage, >1 for fixed width
    height = 0.8,       -- 0-1 for percentage, >1 for fixed height
    title_pos = "center", -- "left" | "center" | "right"
  },

  -- Picker Settings
  picker = {
    provider = "auto",  -- "auto" | "telescope" | "native"
    telescope = {
      theme = "dropdown", -- "dropdown" | "cursor" | "ivy" | nil
      layout_config = {
        width = 0.8,
        height = 0.6,
      },
    },
  },

  -- Auto-Open Generated Files
  auto_open = {
    enabled = true,
    strategy = "first",  -- "first" | "all" | "none"
  },

  -- Quickfix Integration
  quickfix = {
    enabled = true,
    auto_open = false,
  },

  -- Terminal Settings (for future dev commands)
  terminal = {
    type = "auto",      -- "auto" | "float" | "split" | "vsplit" | "tab"
    float = {
      border = "rounded",
      width = 0.8,
      height = 0.8,
    },
    split = {
      size = 15,
      position = "below", -- "below" | "above"
    },
    persist = true,
    auto_close = false,
  },

  -- Dev Mode Settings (future)
  dev = {
    restart_on_save = false,
    save_patterns = { "*.java", "*.kt", "*.xml", "*.yaml", "*.yml", "*.properties" },
  },

  -- Command-Specific Settings
  commands = {
    generate = {
      refresh = false,  -- Force profile re-detection
    },
  },
})
```

---

## Configuration Sections

### CLI Path

Specify the path to the Haft CLI binary.

```lua
haft_path = "haft"  -- Default: auto-detect in PATH
```

**Examples:**

```lua
-- Custom path
haft_path = "/usr/local/bin/haft"

-- Expand home directory
haft_path = vim.fn.expand("~/bin/haft")

-- Relative to project
haft_path = vim.fn.getcwd() .. "/bin/haft"
```

**Auto-detection order:**
1. Value of `haft_path` config
2. `haft` command in system PATH
3. Error if not found

---

### Project Detection

Controls how haft.nvim detects Spring Boot projects.

```lua
detection = {
  enabled = true,
  patterns = { ".haft.yaml", "pom.xml", "build.gradle", "build.gradle.kts" },
}
```

**Options:**

- `enabled` (boolean): Enable/disable automatic detection
- `patterns` (table): File patterns to match for project root

**Disable detection:**
```lua
detection = { enabled = false }
```

**Custom patterns:**
```lua
detection = {
  enabled = true,
  patterns = { "pom.xml" },  -- Maven projects only
}
```

---

### Notifications

Control notification behavior.

```lua
notifications = {
  enabled = true,
  level = "info",
  timeout = 3000,
}
```

**Options:**

- `enabled` (boolean): Show notifications
- `level` (string): Minimum level to display
  - `"debug"` - All messages (verbose)
  - `"info"` - Info, warnings, errors (default)
  - `"warn"` - Warnings and errors only
  - `"error"` - Errors only
- `timeout` (number): Duration in milliseconds

**Examples:**

```lua
-- Quiet mode (errors only)
notifications = { level = "error" }

-- Disable all notifications
notifications = { enabled = false }

-- Longer timeout
notifications = { timeout = 5000 }  -- 5 seconds
```

**Integration with noice.nvim:**

If [noice.nvim](https://github.com/folke/noice.nvim) is installed, notifications automatically use enhanced rendering with better styling.

---

### Floating Windows

Customize floating window appearance for `:HaftInfo`, `:HaftRoutes`, `:HaftStats`.

```lua
float = {
  border = "rounded",
  width = 0.8,
  height = 0.8,
  title_pos = "center",
}
```

**Options:**

- `border` (string): Border style
  - `"none"` - No border
  - `"single"` - Single line `─│┌┐└┘`
  - `"double"` - Double line `═║╔╗╚╝`
  - `"rounded"` - Rounded corners `─│╭╮╰╯` (default)
  - `"solid"` - Solid block
  - `"shadow"` - Drop shadow
- `width` (number):
  - `0.0 - 1.0` - Percentage of screen (e.g., `0.8` = 80%)
  - `> 1` - Fixed columns (e.g., `120`)
- `height` (number): Same as width but for rows
- `title_pos` (string): `"left"`, `"center"`, `"right"`

**Examples:**

```lua
-- Minimal style
float = { border = "single", width = 0.7, height = 0.7 }

-- Full screen
float = { width = 0.95, height = 0.95 }

-- Fixed size
float = { width = 120, height = 40 }  -- 120 columns × 40 rows

-- Left-aligned title
float = { title_pos = "left" }
```

---

### Picker Provider

Choose how interactive pickers work (for `:HaftAdd`, `:HaftRemove`).

```lua
picker = {
  provider = "auto",
  telescope = {
    theme = "dropdown",
    layout_config = { width = 0.8, height = 0.6 },
  },
}
```

**Provider Options:**

- `"auto"` (default) - Use Telescope if available, fallback to `vim.ui.select`
- `"telescope"` - Force Telescope, error if not installed
- `"native"` - Always use `vim.ui.select`

**Telescope Themes:**

- `"dropdown"` (default) - Centered dropdown
- `"cursor"` - Near cursor position
- `"ivy"` - Bottom-aligned
- `nil` - Default Telescope layout

**Examples:**

```lua
-- Force Telescope
picker = { provider = "telescope" }

-- Use native picker (no Telescope required)
picker = { provider = "native" }

-- Custom Telescope theme
picker = {
  provider = "telescope",
  telescope = {
    theme = "ivy",
    layout_config = {
      width = 0.9,
      height = 0.8,
      preview_cutoff = 120,
    },
  },
}
```

---

### Auto-Open Generated Files

Control how generated files are opened after commands like `:HaftGenerateResource`.

```lua
auto_open = {
  enabled = true,
  strategy = "first",
}
```

**Options:**

- `enabled` (boolean): Enable auto-opening
- `strategy` (string):
  - `"first"` (default) - Open only the first file (usually Controller)
  - `"all"` - Open all generated files as buffers
  - `"none"` - Don't auto-open (files still added to quickfix)

**Examples:**

```lua
-- Open all generated files
auto_open = { strategy = "all" }

-- Disable auto-open
auto_open = { enabled = false }
```

---

### Quickfix Integration

Add generated files to the quickfix list for easy navigation.

```lua
quickfix = {
  enabled = true,
  auto_open = false,
}
```

**Options:**

- `enabled` (boolean): Add files to quickfix
- `auto_open` (boolean): Automatically open quickfix window

**Examples:**

```lua
-- Auto-open quickfix after generation
quickfix = { enabled = true, auto_open = true }

-- Then navigate with:
-- :cnext, :cprev, :cfirst, :clast
```

**Disable quickfix:**
```lua
quickfix = { enabled = false }
```

---

### Terminal Settings

Settings for future dev commands (`:HaftDevServe`, etc.). Currently unused.

```lua
terminal = {
  type = "auto",
  float = { border = "rounded", width = 0.8, height = 0.8 },
  split = { size = 15, position = "below" },
  persist = true,
  auto_close = false,
}
```

---

### Command-Specific Settings

Options for individual command behaviors.

```lua
commands = {
  generate = {
    refresh = false,  -- Force profile refresh on each generate
  },
}
```

**Options:**

- `generate.refresh` (boolean): Force re-detection of project profile
  - `false` (default) - Use cached profile
  - `true` - Re-scan project on every generate command

**Example:**

```lua
-- Always refresh profile (slower but detects changes)
commands = { generate = { refresh = true } }
```

---

## Configuration Examples

### Minimal Setup

Just the essentials:

```lua
require("haft").setup({
  haft_path = "haft",
  picker = { provider = "auto" },
})
```

### Power User Setup

Maximum features and customization:

```lua
require("haft").setup({
  -- Custom CLI path
  haft_path = vim.fn.expand("~/.local/bin/haft"),

  -- Quiet notifications
  notifications = { level = "warn", timeout = 2000 },

  -- Large floats
  float = {
    border = "double",
    width = 0.95,
    height = 0.95,
    title_pos = "left",
  },

  -- Force Telescope with custom theme
  picker = {
    provider = "telescope",
    telescope = {
      theme = "ivy",
      layout_config = { width = 0.9, height = 0.7 },
    },
  },

  -- Open all files + quickfix
  auto_open = { strategy = "all" },
  quickfix = { enabled = true, auto_open = true },

  -- Force profile refresh
  commands = { generate = { refresh = true } },
})
```

### No-Frills Setup

Minimal UI, maximum simplicity:

```lua
require("haft").setup({
  notifications = { enabled = false },
  picker = { provider = "native" },
  auto_open = { enabled = false },
  quickfix = { enabled = false },
})
```

### Telescope-First Setup

Optimized for Telescope users:

```lua
require("haft").setup({
  picker = {
    provider = "telescope",
    telescope = {
      theme = "dropdown",
      layout_config = {
        width = 0.9,
        height = 0.8,
        prompt_position = "top",
      },
    },
  },
})

-- Load Telescope extension
require("telescope").load_extension("haft")
```

---

## Related

- [Usage Guide →](./usage) - Learn how to use commands
- [CLI Configuration](/docs/reference/configuration) - CLI-level configuration
- [Troubleshooting →](./troubleshooting) - Common configuration issues
