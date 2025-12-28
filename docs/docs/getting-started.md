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

Haft's **intelligent detection engine** scans your project and generates code that matches your existing patterns:

- **Architecture-aware** — Generates files in the correct location based on your project structure (Layered, Feature-based, Hexagonal, Clean, or Modular)
- **Pattern-matching** — Uses your existing DTO naming (Request/Response vs DTO), ID types (Long vs UUID), and mapper style
- **Profile caching** — First run scans your project, subsequent runs are instant

This creates a complete CRUD structure:

- `UserController.java` — REST endpoints
- `UserService.java` — Service interface
- `UserServiceImpl.java` — Service implementation
- `UserRepository.java` — JPA repository
- `User.java` — Entity
- `UserRequest.java` / `UserResponse.java` — DTOs
- `UserMapper.java` — Entity-DTO mapping
- `UserServiceTest.java` — Unit tests with Mockito
- `UserControllerTest.java` — Integration tests with MockMvc

Or generate individual components:

```bash
haft generate controller Product   # haft g co
haft generate service Order        # haft g s
haft generate repository Payment   # haft g repo
haft generate entity Customer      # haft g e
haft generate dto Invoice
```

### 4. Add Security

```bash
# Interactive selection of security types
haft generate security

# Generate JWT authentication for REST APIs
haft generate security --jwt

# Generate session-based auth for web apps
haft generate security --session

# Generate OAuth2 (Google, GitHub, Facebook)
haft generate security --oauth2
```

This creates a complete security setup:

- `SecurityConfig.java` - Spring Security configuration
- `JwtUtil.java` - Token generation/validation (JWT)
- `AuthenticationController.java` - Login/register endpoints
- `CustomUserDetailsService.java` - User loading
- DTOs for authentication requests/responses

### 5. Add Dependencies

```bash
# Interactive search picker
haft add

# Browse by category
haft add --browse

# Add using shortcuts
haft add lombok validation jwt

# Add database driver
haft add postgresql

# List available shortcuts
haft add --list
```

### 6. Remove Dependencies

```bash
# Interactive picker
haft remove

# Remove by name
haft remove lombok
haft rm h2 validation

# Remove using suffix matching
haft rm jpa   # Removes spring-boot-starter-data-jpa
```

### 7. Run Your Project

```bash
# Maven
./mvnw spring-boot:run

# Gradle
./gradlew bootRun
```

That's it! You have a fully configured Spring Boot project with CRUD endpoints.

## What's Next?

- [Installation](/docs/installation) — Detailed installation options
- [Why Haft?](/docs/why-haft) — Learn what makes Haft different
- [haft init](/docs/commands/init) — Project initialization reference
- [haft generate](/docs/commands/generate) - Resource generation reference
- [haft generate security](/docs/commands/security) - Security configuration
- [haft add](/docs/commands/add) - Add dependencies
- [haft remove](/docs/commands/remove) — Remove dependencies
- [haft completion](/docs/commands/completion) — Shell completions setup
- [Wizard Navigation](/docs/guides/wizard-navigation) — Master the TUI wizard

## Example: Non-Interactive Mode

For CI/CD pipelines or scripting:

```bash
# Create project
haft init my-service \
  --group com.example \
  --java 21 \
  --spring 3.4.0 \
  --deps web,data-jpa,lombok,validation \
  --no-interactive

# Generate resources
cd my-service
haft generate resource User --no-interactive
haft generate resource Product --no-interactive

# Skip test generation
haft generate resource Order --skip-tests --no-interactive

# Force re-scan project (ignore cached profile)
haft generate resource Payment --refresh --no-interactive

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
haft add --help
haft remove --help
haft completion --help
```

## Enable Shell Completions

Set up tab completions for a better experience:

```bash
# Bash
source <(haft completion bash)

# Zsh
source <(haft completion zsh)

# Fish
haft completion fish | source

# PowerShell
haft completion powershell | Out-String | Invoke-Expression
```

See [haft completion](/docs/commands/completion) for permanent installation instructions.

Join the community:
- [GitHub Issues](https://github.com/KashifKhn/haft/issues) — Report bugs
- [GitHub Discussions](https://github.com/KashifKhn/haft/discussions) — Ask questions
