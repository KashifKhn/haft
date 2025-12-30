---
sidebar_position: 6
title: haft dev
description: Development commands for running, building, and testing
---

# haft dev

Development commands that wrap Maven/Gradle for a unified experience.

## Usage

```bash
haft dev <command>
haft d <command>      # Alias
```

## Description

The `dev` command provides a unified interface for common development tasks. It automatically detects your build tool (Maven or Gradle) and executes the appropriate underlying commands.

## Available Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `serve` | `run`, `start` | Start the application with hot-reload |
| `build` | `b`, `compile` | Build the project |
| `test` | `t` | Run tests |
| `clean` | - | Clean build artifacts |
| `restart` | - | Trigger restart of running dev server |

---

## haft dev serve

Start the Spring Boot application in **supervisor mode** with interactive restart support.

### Interactive Mode (Default)

When running in a terminal, `haft dev serve` runs as a supervisor that manages your Spring Boot process. You can use keyboard commands to control the server:

| Key | Action |
|-----|--------|
| `r` | **Restart** - Compiles first, then restarts (keeps old server if compile fails) |
| `q` | **Quit** - Gracefully stops the server and exits |
| `c` | **Clear** - Clears the screen |
| `h` | **Help** - Shows available commands |
| `Ctrl+C` | Same as `q` - Graceful shutdown |

```
╭─────────────────────────────────────────╮
│  Haft Dev Server                        │
│  Press r to restart, q to quit          │
│  Press h for more commands              │
╰─────────────────────────────────────────╯

INFO Starting application build-tool=Maven

  .   ____          _            __ _ _
 /\\ / ___'_ __ _ _(_)_ __  __ _ \ \ \ \
...
```

### Restart Behavior

The restart command (`r`) follows a **compile-first** strategy:

1. Runs compilation (`mvn compile -DskipTests` or `gradle classes -x test`)
2. **If compilation fails**: Shows error, keeps the old server running
3. **If compilation succeeds**: 
   - Gracefully stops the old server (SIGTERM → 2s wait → SIGKILL)
   - Starts a new server instance

This prevents the "dead server" situation where a syntax error leaves you with no running application.

### Usage

```bash
haft dev serve [flags]
haft dev run [flags]     # Alias
haft dev start [flags]   # Alias
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--profile` | `-p` | Spring profile to activate (e.g., dev, prod) |
| `--debug` | `-d` | Enable remote debugging on port 5005 |
| `--port` | - | Server port (overrides application config) |
| `--no-interactive` | - | Disable interactive mode (for CI/scripts) |

### Examples

```bash
# Start with default settings (interactive mode)
haft dev serve

# Start with dev profile
haft dev serve --profile dev

# Start with debug mode enabled
haft dev serve --debug

# Start on a specific port
haft dev serve --port 8081

# Combine options
haft dev serve -p dev -d --port 8081

# Non-interactive mode (for CI/scripts)
haft dev serve --no-interactive
```

### Build Tool Commands

| Build Tool | Run Command | Compile Command (for restart) |
|------------|-------------|-------------------------------|
| Maven | `mvn spring-boot:run -DskipTests` | `mvn compile -DskipTests -q` |
| Gradle | `./gradlew bootRun -x test` | `./gradlew classes -x test -q` |

---

## Plugin Integration

External tools (Neovim, VSCode, IntelliJ) can trigger a restart without keyboard interaction.

### Option 1: Using `haft dev restart` Command

The simplest way to trigger a restart from any tool:

```bash
haft dev restart
```

This command creates the trigger file that signals the running dev server to restart.

### Option 2: Creating the Trigger File Directly

When `haft dev serve` is running, it watches for the creation/modification of:

```
.haft/restart
```

To trigger a restart manually:

```bash
touch .haft/restart
```

### Neovim Integration

Add to your `init.lua`:

```lua
-- Auto-restart haft on save (using shell command)
vim.api.nvim_create_autocmd("BufWritePost", {
  pattern = { "*.java", "*.kt" },
  callback = function()
    vim.fn.jobstart("haft dev restart", { detach = true })
  end,
})

-- Manual restart keybinding
vim.keymap.set("n", "<leader>hr", function()
  vim.fn.jobstart("haft dev restart", { detach = true })
  print("Haft restart triggered")
end, { desc = "Trigger Haft restart" })
```

Alternative (direct file creation):

```lua
vim.keymap.set("n", "<leader>hr", function()
  local trigger = vim.fn.getcwd() .. "/.haft/restart"
  vim.fn.writefile({}, trigger)
  print("Haft restart triggered")
end, { desc = "Trigger Haft restart" })
```

### VS Code Integration

Add to `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Haft Restart",
      "type": "shell",
      "command": "haft dev restart",
      "problemMatcher": [],
      "presentation": {
        "reveal": "silent"
      }
    }
  ]
}
```

Add keybinding in `keybindings.json`:

```json
{
  "key": "ctrl+shift+r",
  "command": "workbench.action.tasks.runTask",
  "args": "Haft Restart"
}
```

### IntelliJ IDEA Integration

1. Go to **Settings → Tools → External Tools**
2. Add new tool:
   - **Name**: Haft Restart
   - **Program**: `haft`
   - **Arguments**: `dev restart`
   - **Working directory**: `$ProjectFileDir$`
3. Assign a keyboard shortcut in **Settings → Keymap**

---

## haft dev restart

Trigger a restart of the running dev server from the command line.

### Usage

```bash
haft dev restart
```

### Description

This command creates a trigger file that signals the running dev server (`haft dev serve`) to restart. It's useful for:

- Shell scripts that need to trigger restarts
- Editor plugins that prefer calling CLI commands
- CI/CD pipelines for hot-reload testing
- Any tool that can execute shell commands

### Examples

```bash
# Trigger restart of running dev server
haft dev restart

# Use in a shell script
#!/bin/bash
# Edit some files...
haft dev restart

# Use with file watchers
fswatch -o src/ | xargs -n1 -I{} haft dev restart

# Use with entr (run on file changes)
find src -name "*.java" | entr haft dev restart
```

---

## haft dev build

Build the Spring Boot project.

### Usage

```bash
haft dev build [flags]
haft dev b [flags]        # Alias
haft dev compile [flags]  # Alias
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--skip-tests` | `-s` | Skip running tests during build |
| `--profile` | `-p` | Maven/Gradle profile to activate |
| `--clean` | `-c` | Clean before building |

### Examples

```bash
# Build the project
haft dev build

# Build without running tests
haft dev build --skip-tests
haft dev build -s

# Clean and build
haft dev build --clean
haft dev build -c

# Build with production profile
haft dev build --profile prod

# Clean build without tests
haft dev build -c -s
```

### Build Tool Commands

| Build Tool | Executed Command |
|------------|------------------|
| Maven | `mvn package` (or `mvn clean package`) |
| Gradle | `./gradlew build` (or `./gradlew clean build`) |

---

## haft dev test

Run project tests.

### Usage

```bash
haft dev test [flags]
haft dev t [flags]    # Alias
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--filter` | `-f` | Filter tests by class or method name |
| `--verbose` | `-v` | Enable verbose test output |
| `--fail-fast` | - | Stop on first test failure |

### Examples

```bash
# Run all tests
haft dev test

# Run tests matching a pattern
haft dev test --filter UserService
haft dev test -f UserController

# Run with verbose output
haft dev test --verbose

# Stop on first failure
haft dev test --fail-fast

# Combine options
haft dev test -f UserService -v --fail-fast
```

### Build Tool Commands

| Build Tool | Executed Command |
|------------|------------------|
| Maven | `mvn test` |
| Gradle | `./gradlew test` |

---

## haft dev clean

Clean build artifacts and compiled files.

### Usage

```bash
haft dev clean
```

### Examples

```bash
# Clean build artifacts
haft dev clean
```

### Build Tool Commands

| Build Tool | Executed Command |
|------------|------------------|
| Maven | `mvn clean` |
| Gradle | `./gradlew clean` |

---

## Build Tool Detection

Haft automatically detects your build tool by looking for build files:

| File Found | Build Tool |
|------------|------------|
| `pom.xml` | Maven |
| `build.gradle.kts` | Gradle (Kotlin DSL) |
| `build.gradle` | Gradle (Groovy DSL) |

### Wrapper Detection

Haft prefers build tool wrappers when available:

- Uses `./mvnw` over `mvn` if `mvnw` exists
- Uses `./gradlew` over `gradle` if `gradlew` exists

This ensures consistent builds across different environments.

## Typical Workflow

```bash
# Start development server (interactive mode)
haft dev serve -p dev

# Press 'r' to restart after making changes
# Or trigger restart from your editor

# In another terminal, run tests
haft dev test

# Build for production
haft dev build -c -s --profile prod

# Clean up
haft dev clean
```

## See Also

- [haft init](/docs/commands/init) - Initialize new projects
- [haft generate](/docs/commands/generate) - Generate code
- [haft add](/docs/commands/add) - Add dependencies
