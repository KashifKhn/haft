---
sidebar_position: 1
title: haft init
description: Initialize a new Spring Boot project
---

# haft init

Initialize a new Spring Boot project with an interactive wizard or command-line flags.

## Usage

```bash
haft init [name] [flags]
```

## Description

The `init` command creates a new Spring Boot project. When run without flags, it launches an interactive TUI wizard that guides you through the configuration.

The wizard presents all dependencies from Spring Initializr organized by category (Web, SQL, Security, etc.) with descriptions and search functionality.

## Examples

### Interactive Mode (Recommended)

```bash
haft init
```

Launches the full interactive wizard.

### With Project Name

```bash
haft init my-app
```

Starts the wizard with the project name pre-filled.

### Non-Interactive Mode

```bash
haft init my-app \
  --group com.example \
  --artifact my-app \
  --java 21 \
  --spring 3.4.1 \
  --build maven \
  --deps web,data-jpa,lombok \
  --no-interactive
```

Creates a project without any prompts.

### Quick Start with Defaults

```bash
haft init my-api -g com.company -j 21
```

Uses defaults for unspecified options.

## Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--group` | `-g` | Group ID (e.g., `com.example`) | `com.example` |
| `--artifact` | `-a` | Artifact ID | Project name |
| `--java` | `-j` | Java version (`17`, `21`, `25`) | `21` |
| `--spring` | `-s` | Spring Boot version | Latest stable |
| `--build` | `-b` | Build tool (`maven`, `gradle`, `gradle-kotlin`) | `maven` |
| `--deps` | | Dependencies (comma-separated) | None |
| `--package` | | Base package name | Auto-generated |
| `--packaging` | | Packaging type (`jar`, `war`) | `jar` |
| `--config` | | Config format (`properties`, `yaml`) | `yaml` |
| `--description` | | Project description | Empty |
| `--dir` | `-d` | Output directory | Current directory |
| `--no-interactive` | | Skip interactive wizard | `false` |

## Global Flags

| Flag | Description |
|------|-------------|
| `--verbose` | Enable verbose output |
| `--no-color` | Disable colored output |

## Wizard Steps

The interactive wizard guides you through 12 configuration steps:

| Step | Field | Description |
|------|-------|-------------|
| 1 | Project Name | Name of your application |
| 2 | Group ID | Maven group (e.g., `com.example`) |
| 3 | Artifact ID | Maven artifact name |
| 4 | Description | Project description |
| 5 | Package Name | Base package (auto-generated from group + artifact) |
| 6 | Java Version | 17 (LTS), 21 (LTS - Recommended), or 25 |
| 7 | Spring Boot | Select from available versions |
| 8 | Build Tool | Maven, Gradle (Groovy), or Gradle (Kotlin DSL) |
| 9 | Packaging | JAR or WAR |
| 10 | Config Format | Properties or YAML |
| 11 | Dependencies | Search and select from all Spring starters |
| 12 | Git Init | Initialize a Git repository |

## Dependency Selection

The dependency picker provides:

- **Search**: Press `/` to search by name or description
- **Categories**: Press `Tab` to cycle, or `0-9` to jump directly
- **Selection**: Press `Space` to toggle, `Enter` to confirm

### Available Categories

| Key | Category | Examples |
|-----|----------|----------|
| `0` | All | Show all dependencies |
| `1` | Developer Tools | DevTools, Lombok, Configuration Processor |
| `2` | Web | Spring Web, WebFlux, GraphQL |
| `3` | SQL | JPA, JDBC, MySQL, PostgreSQL |
| `4` | NoSQL | MongoDB, Redis, Elasticsearch |
| `5` | Security | Spring Security, OAuth2 |
| `6` | Messaging | Kafka, RabbitMQ |
| `7` | Cloud | Config, Discovery, Gateway |
| `8` | Observability | Actuator, Micrometer |
| `9` | Testing | Testcontainers |

## Generated Structure

### Maven Project

```
my-app/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/example/myapp/
│   │   │       └── MyAppApplication.java
│   │   └── resources/
│   │       └── application.yml
│   └── test/
│       └── java/
│           └── com/example/myapp/
│               └── MyAppApplicationTests.java
├── .gitignore
├── .haft.yaml
├── mvnw
├── mvnw.cmd
└── pom.xml
```

### Gradle Project

```
my-app/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/example/myapp/
│   │   │       └── MyAppApplication.java
│   │   └── resources/
│   │       └── application.yml
│   └── test/
│       └── java/
│           └── com/example/myapp/
│               └── MyAppApplicationTests.java
├── gradle/
│   └── wrapper/
│       └── gradle-wrapper.properties
├── .gitignore
├── .haft.yaml
├── build.gradle          # or build.gradle.kts for Kotlin DSL
├── settings.gradle       # or settings.gradle.kts for Kotlin DSL
├── gradlew
└── gradlew.bat
```

## Exit Codes

| Code | Description |
|------|-------------|
| `0` | Success |
| `1` | Error (invalid flags, generation failed) |

## See Also

- [Wizard Navigation](/docs/guides/wizard-navigation) — Keyboard shortcuts
- [Dependencies](/docs/guides/dependencies) — Dependency details
- [Project Structure](/docs/guides/project-structure) — Generated file details
