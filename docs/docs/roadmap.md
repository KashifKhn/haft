---
sidebar_position: 11
title: Roadmap
description: Haft development roadmap and future plans
---

# Roadmap

This page tracks the development progress of Haft.

## Current Version: 0.1.x

### Completed

- [x] **Project Initialization**
  - [x] Interactive TUI wizard
  - [x] Spring Initializr dependency integration
  - [x] Maven project generation
  - [x] Gradle project generation (basic)
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
  - [x] Dependency detection (Lombok, MapStruct, etc.)

### In Progress

- [ ] **Resource Generation** (`haft generate resource`)
  - [ ] Entity generation
  - [ ] Repository generation
  - [ ] Service generation
  - [ ] Controller generation
  - [ ] DTO generation
  - [ ] Mapper generation

## Upcoming: 0.2.0

### Planned Features

- [ ] **Individual Generators**
  - [ ] `haft generate controller`
  - [ ] `haft generate service`
  - [ ] `haft generate entity`
  - [ ] `haft generate repository`

- [ ] **Dependency Manager**
  - [ ] `haft add <dependency>`
  - [ ] `haft remove <dependency>`
  - [ ] Interactive dependency search

- [ ] **Shell Completions**
  - [ ] Bash completion
  - [ ] Zsh completion
  - [ ] Fish completion
  - [ ] PowerShell completion

## Future: 0.3.0+

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

- [ ] **Gradle Improvements**
  - [ ] Full Gradle support
  - [ ] Gradle Kotlin DSL

## Contributing

Want to help? Check the [GitHub Issues](https://github.com/KashifKhn/haft/issues) for tasks labeled:

- `good first issue` — Great for new contributors
- `help wanted` — We need assistance
- `enhancement` — Feature requests

See [Contributing](/docs/contributing) for guidelines.

## Changelog

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
