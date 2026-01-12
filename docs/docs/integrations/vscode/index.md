---
sidebar_position: 2
title: VS Code Extension
description: Haft integration for Visual Studio Code (Coming Soon)
---

# VS Code Extension

🔜 **Coming Soon**

The Haft VS Code extension will bring Spring Boot productivity tools directly into Visual Studio Code with native UI integration, command palette commands, and IntelliSense support.

## Planned Features

- ✨ **Command Palette Integration** - All Haft commands via `Cmd/Ctrl+Shift+P`
- 🎨 **Native UI** - Settings panel, sidebar views, notifications
- 📂 **File Explorer Integration** - Right-click context menus
- 💡 **IntelliSense** - Autocomplete for dependencies and commands
- ⚡ **Async Execution** - Non-blocking operations
- 🔔 **Progress Notifications** - Visual feedback for long-running tasks
- 📊 **Status Bar** - Quick project information
- 🎯 **Code Actions** - Quick fixes and refactorings

## Planned Commands

- `Haft: Generate Resource` - Generate complete CRUD resource
- `Haft: Add Dependency` - Interactive dependency picker
- `Haft: Remove Dependency` - Remove dependencies
- `Haft: View Project Info` - Show project details
- `Haft: View Routes` - List API endpoints
- `Haft: View Stats` - Code statistics

## Planned UI Elements

### Sidebar View

- Project information panel
- Dependency list (add/remove)
- Recent generations
- Quick actions

### Settings

- CLI binary path
- Generation preferences
- Notification preferences
- Template customization

### Status Bar

- Project type indicator
- Spring Boot version
- Quick action buttons

## Architecture

Like all Haft integrations, the VS Code extension will be a **thin wrapper** around the Haft CLI:

```
┌────────────────────────────────────┐
│       VS Code Extension            │
│  • Command Palette                 │
│  • UI Panels                       │
│  • Settings Integration            │
└────────────┬───────────────────────┘
             │
             ▼
┌────────────────────────────────────┐
│         Haft CLI (Core)            │
│  • All business logic              │
│  • Architecture detection          │
│  • Code generation                 │
└────────────────────────────────────┘
```

## Interested in Contributing?

If you'd like to help build the VS Code extension:

1. Check the [roadmap](/docs/roadmap) for timeline
2. Join the [discussion](https://github.com/KashifKhn/haft/discussions)
3. See [contributing guidelines](/docs/contributing)

## Meanwhile, Use the CLI

While the VS Code extension is in development, you can use the Haft CLI directly from the VS Code integrated terminal:

```bash
# In VS Code terminal (Ctrl+`)
haft generate resource User
haft add lombok validation
haft info
```

See [CLI Documentation](/docs/getting-started) to learn all CLI commands.

## Stay Updated

- **GitHub**: [KashifKhn/haft](https://github.com/KashifKhn/haft)
- **Discussions**: [Feature Requests](https://github.com/KashifKhn/haft/discussions)
- **Roadmap**: [What's Coming →](/docs/roadmap)

---

**Want to be notified when the VS Code extension launches?** Watch the [GitHub repository](https://github.com/KashifKhn/haft) for updates.
