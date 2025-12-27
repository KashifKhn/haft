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

---

## haft dev serve

Start the Spring Boot application with DevTools for hot-reload.

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

### Examples

```bash
# Start with default settings
haft dev serve

# Start with dev profile
haft dev serve --profile dev

# Start with debug mode enabled
haft dev serve --debug

# Start on a specific port
haft dev serve --port 8081

# Combine options
haft dev serve -p dev -d --port 8081
```

### Build Tool Commands

| Build Tool | Executed Command |
|------------|------------------|
| Maven | `mvn spring-boot:run` (or `./mvnw`) |
| Gradle | `./gradlew bootRun` (or `gradle bootRun`) |

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
# Start development server
haft dev serve -p dev

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
