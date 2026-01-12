---
sidebar_position: 3
title: IntelliJ IDEA Plugin
description: Haft integration for IntelliJ IDEA (Coming Soon)
---

# IntelliJ IDEA Plugin

🔜 **Coming Soon**

The Haft IntelliJ IDEA plugin will bring Spring Boot productivity tools directly into JetBrains IDEs with native UI integration, tool windows, and intention actions.

## Planned Features

- 🛠️ **Tool Windows** - Dedicated panels for project info, routes, and statistics
- 💡 **Intention Actions** - Quick fixes and code generation via `Alt+Enter`
- 🎯 **Code Actions** - Context-aware generation in editor
- 📂 **Project View Integration** - Right-click menus in project tree
- ⚡ **Background Tasks** - Non-blocking operations with progress indicators
- 🔔 **Balloon Notifications** - Native IDE notifications
- ⚙️ **Settings Dialog** - Integrated preferences panel
- 🎨 **UI Theme Support** - Matches IntelliJ Light/Dark themes

## Planned Actions

- `Generate Haft Resource` - Generate complete CRUD resource
- `Add Dependency` - Interactive dependency selector
- `Remove Dependency` - Manage dependencies
- `View Project Information` - Show project details in tool window
- `View API Routes` - List endpoints in tool window
- `View Code Statistics` - Code metrics and COCOMO estimates

## Planned UI Elements

### Tool Windows

- **Haft Project** - Project information, dependencies, recent actions
- **Haft Routes** - API endpoint browser
- **Haft Stats** - Code statistics and metrics

### Intention Actions

Press `Alt+Enter` in Java files for:

- Generate Resource for class
- Generate Controller
- Generate Service
- Generate Repository
- Generate DTO

### Context Menus

Right-click in:

- Project view → `Haft > Generate...`
- Editor → `Haft > Generate for...`
- Package → `Haft > Add Resource`

### Settings Dialog

`Settings → Tools → Haft`:

- CLI binary path
- Generation preferences
- Notification settings
- Template customization

## Architecture

The IntelliJ plugin will be a **thin wrapper** around the Haft CLI, following the same pattern as other integrations:

```
┌────────────────────────────────────┐
│      IntelliJ IDEA Plugin          │
│  • Tool Windows                    │
│  • Intention Actions               │
│  • Native UI Integration           │
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

## Supported IDEs

The plugin will support all JetBrains IDEs:

- IntelliJ IDEA (Community & Ultimate)
- Android Studio
- GoLand (if working on Go + Spring Boot)
- WebStorm (if working on frontend + Spring Boot backend)

## Interested in Contributing?

If you'd like to help build the IntelliJ plugin:

1. Check the [roadmap](/docs/roadmap) for timeline
2. Join the [discussion](https://github.com/KashifKhn/haft/discussions)
3. See [contributing guidelines](/docs/contributing)
4. Experience with IntelliJ Platform SDK is helpful

## Meanwhile, Use the CLI

While the IntelliJ plugin is in development, you can use the Haft CLI directly from the IntelliJ terminal:

```bash
# In IntelliJ terminal (Alt+F12)
haft generate resource User
haft add lombok validation
haft info
```

Or add as an External Tool:

1. `Settings → Tools → External Tools`
2. Click `+` to add new tool
3. Configure:
   - **Name**: Haft Generate Resource
   - **Program**: `haft`
   - **Arguments**: `generate resource $Prompt$`
   - **Working directory**: `$ProjectFileDir$`

Then access via `Tools → External Tools → Haft Generate Resource`.

See [CLI Documentation](/docs/getting-started) to learn all CLI commands.

## Stay Updated

- **GitHub**: [KashifKhn/haft](https://github.com/KashifKhn/haft)
- **Discussions**: [Feature Requests](https://github.com/KashifKhn/haft/discussions)
- **Roadmap**: [What's Coming →](/docs/roadmap)

---

**Want to be notified when the IntelliJ plugin launches?** Watch the [GitHub repository](https://github.com/KashifKhn/haft) for updates.
