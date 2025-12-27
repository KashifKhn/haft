---
sidebar_position: 11
title: Roadmap
description: Haft development roadmap and future plans
---

# Roadmap

This page tracks the development progress of Haft.

## Current Version: 0.3.x

### Completed

- [x] **Project Initialization**
  - [x] Interactive TUI wizard
  - [x] Spring Initializr dependency integration
  - [x] Maven project generation
  - [x] Gradle project generation (Groovy & Kotlin DSL)
  - [x] YAML and Properties config formats
  - [x] Git repository initialization
  - [x] **Offline operation** — No internet required

- [x] **Wizard Features**
  - [x] 12-step configuration wizard
  - [x] Back navigation (Esc key)
  - [x] Dynamic package name generation
  - [x] Dependency search (`/` key)
  - [x] Category filtering (Tab, 0-9 keys)

- [x] **Maven Parser**
  - [x] Read pom.xml files
  - [x] Write pom.xml files
  - [x] Dependency detection (Lombok, JPA, Validation)
  - [x] Add/Remove dependencies

- [x] **Gradle Parser**
  - [x] Read build.gradle (Groovy DSL)
  - [x] Read build.gradle.kts (Kotlin DSL)
  - [x] Write build.gradle files
  - [x] Write build.gradle.kts files
  - [x] Dependency detection (Lombok, JPA, Validation)
  - [x] Add/Remove dependencies

- [x] **Resource Generation** (`haft generate resource`)
  - [x] Controller generation with CRUD endpoints
  - [x] Service interface generation
  - [x] Service implementation generation
  - [x] Repository generation
  - [x] Entity generation (with/without Lombok)
  - [x] Request/Response DTO generation
  - [x] Mapper generation
  - [x] ResourceNotFoundException generation
  - [x] Smart dependency detection from pom.xml and build.gradle
  - [x] Interactive TUI wizard
  - [x] File safety (skip existing files)

- [x] **Individual Generators**
  - [x] `haft generate controller` (alias: `co`)
  - [x] `haft generate service` (alias: `s`)
  - [x] `haft generate entity` (alias: `e`)
  - [x] `haft generate repository` (alias: `repo`)
  - [x] `haft generate dto`

- [x] **Dependency Manager**
  - [x] `haft add <dependency>` — Add dependencies
  - [x] `haft add` — Interactive TUI search picker
  - [x] `haft add --browse` — Category browser
  - [x] Maven Central verification for coordinates
  - [x] Auto-fetch latest versions
  - [x] `haft remove <dependency>` — Remove dependencies
  - [x] `haft remove` — Interactive removal picker
  - [x] 330+ dependency shortcuts (jwt, guava, feign, etc.)

## Upcoming: 0.3.0

### Completed

- [x] **Shell Completions**
  - [x] Bash completion
  - [x] Zsh completion
  - [x] Fish completion
  - [x] PowerShell completion

- [x] **Development Commands** (`haft dev`)
  - [x] `haft dev serve` — Start with hot-reload
  - [x] `haft dev build` — Build project
  - [x] `haft dev test` — Run tests
  - [x] `haft dev clean` — Clean artifacts
  - [x] Auto-detect Maven/Gradle
  - [x] Wrapper support (mvnw/gradlew)

- [x] **Project Analysis**
  - [x] `haft info` — Show project information
  - [x] `haft info --loc` — Include lines of code summary
  - [x] `haft routes` — List REST endpoints (Java & Kotlin)
  - [x] `haft routes --files` — Show file locations
  - [x] `haft stats` — Code statistics with language breakdown
  - [x] `haft stats --cocomo` — COCOMO cost estimates
  - [x] JSON output support

### Planned Features

- [ ] **Custom Templates**
  - [ ] Local template directory
  - [ ] Project-level templates

## Future: 0.4.0+

### Editor Integration

We believe developers should stay in their editor. The first integration will be for Neovim.

- [ ] **Neovim Plugin** (Priority)
  - [ ] `:HaftInit` — Initialize project from Neovim
  - [ ] `:HaftGenerate` — Generate resources
  - [ ] `:HaftAdd` — Add dependencies
  - [ ] Telescope integration for dependency search
  - [ ] Floating window for wizard

- [ ] **VS Code Extension**
  - [ ] Command palette integration
  - [ ] Sidebar panel
  - [ ] Status bar integration

- [ ] **IntelliJ Plugin**
  - [ ] Tool window integration
  - [ ] Actions and shortcuts

### Advanced Features

- [ ] **Custom Templates**
  - [ ] Local template directory
  - [ ] Project-level templates
  - [ ] Template inheritance
  - [ ] Template variables

- [ ] **Architecture Support**
  - [ ] Hexagonal architecture
  - [ ] Clean architecture
  - [ ] Modular monolith

- [ ] **Additional Generators**
  - [ ] Exception handler generation
  - [ ] Configuration class generation
  - [ ] Test class generation
  - [ ] Security configuration

## Contributing

Want to help? Check the [GitHub Issues](https://github.com/KashifKhn/haft/issues) for tasks labeled:

- `good first issue` — Great for new contributors
- `help wanted` — We need assistance
- `enhancement` — Feature requests

See [Contributing](/docs/contributing) for guidelines.

## Changelog

### v0.3.0 (Current)

- Feature: `haft stats` command — Code statistics using SCC
- Feature: `haft stats --cocomo` — COCOMO cost estimates
- Feature: `haft info --loc` — Lines of code summary
- Feature: `haft routes --files` — Show file locations
- Feature: `haft routes` Kotlin support — Scans .kt files
- Feature: `haft info` command — Show project information
- Feature: `haft routes` command — List REST endpoints
- Feature: `haft dev` command for development workflow
- Feature: `haft dev serve` — Start application with hot-reload
- Feature: `haft dev build` — Build project with profiles
- Feature: `haft dev test` — Run tests with filtering
- Feature: `haft dev clean` — Clean build artifacts
- Feature: `haft completion` command for shell completions
- Feature: Bash, Zsh, Fish, PowerShell completion support
- Feature: Full Gradle support (Groovy & Kotlin DSL)
- Feature: Gradle parser for add/remove/generate commands
- Feature: Gradle project generation with wrapper

### v0.2.0

- Feature: `haft add` interactive TUI search picker
- Feature: `haft add --browse` category browser
- Feature: Maven Central API verification for coordinates
- Feature: Auto-fetch latest version for Maven coordinates
- Feature: `haft remove` command with interactive picker
- Feature: 330+ dependency shortcuts (jwt, guava, feign, resilience4j, etc.)
- Feature: Suffix matching for remove command

### v0.1.3

- Feature: `haft generate controller|service|entity|repository|dto` commands
- Feature: Individual component generation
- Feature: `haft add` basic command with shortcuts
- Feature: Dependency catalog with 30+ shortcuts

### v0.1.2

- Feature: `haft generate resource` command
- Feature: Interactive wizard for resource generation
- Feature: Auto-detect Lombok, JPA, Validation from pom.xml
- Feature: Smart code generation based on dependencies

### v0.1.1

- Fix: Config format default to YAML
- Fix: Version injection via ldflags
- Fix: Install script spinner animation

### v0.1.0

- Initial release
- `haft init` command with full wizard
- Spring Initializr integration
- Maven project generation
- Maven parser for pom.xml
- Offline operation

---

*Last updated: December 2024*
