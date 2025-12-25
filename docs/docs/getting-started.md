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
curl -fsSL https://raw.githubusercontent.com/KashifKhn/haft/main/install.sh | bash
```

Or using Go:

```bash
go install github.com/KashifKhn/haft/cmd/haft@latest
```

See [Installation](/docs/installation) for more options.

### 2. Create a New Project

```bash
haft init
```

This launches an interactive wizard that guides you through:

- Project name and Maven coordinates
- Java and Spring Boot versions
- Build tool selection (Maven/Gradle)
- Dependency selection from all Spring starters

### 3. Generate Resources

```bash
cd my-app
haft generate resource User
```

This creates a complete CRUD structure:

- `UserController.java` — REST endpoints
- `UserService.java` — Service interface
- `UserServiceImpl.java` — Service implementation
- `UserRepository.java` — JPA repository
- `User.java` — Entity
- `UserRequest.java` / `UserResponse.java` — DTOs
- `UserMapper.java` — Entity-DTO mapping

Or generate individual components:

```bash
haft generate controller Product   # haft g co
haft generate service Order        # haft g s
haft generate repository Payment   # haft g repo
haft generate entity Customer      # haft g e
haft generate dto Invoice          # haft g dto
```

### 4. Run Your Project

```bash
./mvnw spring-boot:run
```

That's it! You have a fully configured Spring Boot project with CRUD endpoints.

## What's Next?

- [Installation](/docs/installation) — Detailed installation options
- [Why Haft?](/docs/why-haft) — Learn what makes Haft different
- [haft init](/docs/commands/init) — Project initialization reference
- [haft generate](/docs/commands/generate) — Resource generation reference
- [Wizard Navigation](/docs/guides/wizard-navigation) — Master the TUI wizard

## Example: Non-Interactive Mode

For CI/CD pipelines or scripting:

```bash
# Create project
haft init my-service \
  --group com.example \
  --java 21 \
  --spring 3.4.1 \
  --deps web,data-jpa,lombok,validation \
  --no-interactive

# Generate resources
cd my-service
haft generate resource User --no-interactive
haft generate resource Product --no-interactive

# Or generate individual components
haft generate controller Order --no-interactive
haft generate entity Customer --no-interactive
```

## Getting Help

```bash
# General help
haft --help

# Command-specific help
haft init --help
haft generate --help
haft generate resource --help
haft generate controller --help
```

Join the community:
- [GitHub Issues](https://github.com/KashifKhn/haft/issues) — Report bugs
- [GitHub Discussions](https://github.com/KashifKhn/haft/discussions) — Ask questions
