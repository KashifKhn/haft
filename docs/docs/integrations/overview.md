---
sidebar_position: 1
title: Editor Integrations
description: Use Haft from your favorite editor
---

# Editor Integrations

Haft CLI can be used directly from the terminal, or integrated into your favorite editor for a seamless development experience. All integrations are lightweight wrappers around the CLI, providing enhanced UI and ergonomics while maintaining full feature parity.

## Available Integrations

| Editor | Status | Repository | Description |
|--------|--------|------------|-------------|
| **🌙 Neovim** | ✅ **Stable** | [haft.nvim](https://github.com/KashifKhn/haft.nvim) | Async commands, Telescope pickers, floating windows |
| **💎 VS Code** | 🔜 **Coming Soon** | - | Command palette, sidebar integration, IntelliSense |
| **🧠 IntelliJ IDEA** | 🔜 **Coming Soon** | - | Tool windows, intention actions, native UI |

## Why Use an Editor Integration?

While the CLI is powerful on its own, editor integrations provide:

- **🚀 Faster Workflow** - Run commands without leaving your editor
- **🎨 Enhanced UI** - Interactive pickers, floating windows, notifications
- **⌨️ Keybindings** - Customize shortcuts for common operations
- **📂 File Integration** - Auto-open generated files, quickfix lists
- **🔄 Async Execution** - Non-blocking operations, continue editing while commands run

## Architecture

All editor integrations follow the same design philosophy:

```
┌─────────────────────────────────────────────┐
│         Editor Plugin (Thin Wrapper)        │
│  • UI Components (pickers, floats, etc.)    │
│  • Keybindings & Commands                   │
│  • Async execution                          │
└─────────────────┬───────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────┐
│              Haft CLI (Core)                │
│  • All business logic                       │
│  • Architecture detection                   │
│  • Code generation                          │
│  • Template rendering                       │
└─────────────────────────────────────────────┘
```

**Key Principle:** The CLI is the source of truth. Plugins provide convenience and UI, but all functionality lives in the CLI.

## Which Integration Should I Use?

Choose based on your primary editor:

### Neovim
✅ **Best for:** Vim power users, terminal-focused developers  
✅ **Features:** Telescope pickers, Lua API, floating windows, async execution  
✅ **Status:** Production-ready  
[Get Started →](./neovim/)

### VS Code
🔜 **Best for:** Modern IDE experience, GUI-focused developers  
🔜 **Planned Features:** Command palette, sidebar, settings UI, IntelliSense  
🔜 **Status:** Planned  
[Learn More →](./vscode/)

### IntelliJ IDEA
🔜 **Best for:** JetBrains users, Spring Boot developers  
🔜 **Planned Features:** Tool windows, intention actions, native UI integration  
🔜 **Status:** Planned  
[Learn More →](./intellij/)

## Common Workflow Example

Regardless of which integration you use, the workflow is the same:

**Goal:** Generate a new `Product` CRUD resource

### From CLI:
```bash
cd my-spring-app
haft generate resource Product
```

### From Neovim:
```vim
:HaftGenerateResource Product
" Auto-opens ProductController.java in current buffer
```

### From VS Code (Planned):
```
Cmd+Shift+P → "Haft: Generate Resource" → Type "Product" → Enter
```

### From IntelliJ (Planned):
```
Alt+Enter → "Generate Haft Resource" → Type "Product" → Enter
```

**Result:** All integrations produce the same 8 files, following your project's conventions.

## Learn More

- [CLI Documentation](../getting-started.md) - Understand what each command does
- [Template System](../guides/custom-templates.md) - Customize generated code

## Contributing

Interested in building an integration for another editor (Emacs, Sublime Text, etc.)? See our [Contributing Guide](../contributing.md).
