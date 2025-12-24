---
sidebar_position: 1
slug: /getting-started
title: Getting Started
description: Get up and running with Haft in under 5 minutes
---

# Getting Started

Get up and running with Haft in under 5 minutes.

## What is Haft?

Haft is a command-line tool that supercharges Spring Boot development. While Spring Initializr gets you started with a new project, Haft serves as your **lifecycle companion** — helping you generate boilerplate code, manage dependencies, and maintain architectural consistency throughout your project.

## Quick Start

### 1. Install Haft

```bash
go install github.com/KashifKhn/haft/cmd/haft@latest
```

Or download from [GitHub Releases](https://github.com/KashifKhn/haft/releases).

### 2. Create a New Project

```bash
haft init
```

This launches an interactive wizard that guides you through:

- Project name and Maven coordinates
- Java and Spring Boot versions
- Build tool selection (Maven/Gradle)
- Dependency selection from all Spring starters

### 3. Run Your Project

```bash
cd my-app
./mvnw spring-boot:run
```

That's it! You have a fully configured Spring Boot project.

## What's Next?

- [Installation](/docs/installation) — Detailed installation options
- [Why Haft?](/docs/why-haft) — Learn what makes Haft different
- [haft init](/docs/commands/init) — Complete command reference
- [Wizard Navigation](/docs/guides/wizard-navigation) — Master the TUI wizard

## Example: Non-Interactive Mode

For CI/CD pipelines or scripting:

```bash
haft init my-service \
  --group com.example \
  --java 21 \
  --spring 3.4.1 \
  --deps web,data-jpa,lombok,validation \
  --no-interactive
```

## Getting Help

```bash
# General help
haft --help

# Command-specific help
haft init --help
```

Join the community:
- [GitHub Issues](https://github.com/KashifKhn/haft/issues) — Report bugs
- [GitHub Discussions](https://github.com/KashifKhn/haft/discussions) — Ask questions
