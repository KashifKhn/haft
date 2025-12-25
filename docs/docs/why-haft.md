---
sidebar_position: 3
title: Why Haft?
description: What makes Haft different from Spring Initializr
---

# Why Haft?

## The Problem

Spring Initializr is great for bootstrapping a new project. You visit the website, select your dependencies, download a ZIP, and you're ready to go.

But then what?

Every time you need to add a new entity, you manually create:

1. `User.java` — Entity class
2. `UserRepository.java` — JPA repository
3. `UserService.java` — Service interface
4. `UserServiceImpl.java` — Service implementation
5. `UserController.java` — REST controller
6. `UserRequest.java` — Request DTO
7. `UserResponse.java` — Response DTO
8. `UserMapper.java` — MapStruct mapper

**That's 8 files. Every. Single. Time.**

And you need to:
- Remember the correct package structure
- Copy-paste from existing files
- Update imports
- Wire up dependencies
- Hope you didn't make a typo

## The Solution

Haft is your **lifecycle companion**. It doesn't just bootstrap — it stays with you.

```bash
# Instead of creating 8 files manually...
haft generate resource User

# Done. All files generated with proper architecture.
```

## What Makes Haft Different

### 1. Interactive TUI Wizard

Not a web form. A beautiful terminal interface with:

- Real-time search across 100+ dependencies
- Category filtering (Web, SQL, Security, etc.)
- Keyboard navigation
- Back button support (press `Esc`)

### 2. Spring Initializr Integration

Haft uses the same dependency metadata as Spring Initializr. You get access to all official Spring starters with their descriptions.

### 3. Smart Detection

When generating code, Haft reads your build file (`pom.xml` or `build.gradle`) to detect:

- **Lombok** → Uses `@Data`, `@Builder`, etc.
- **MapStruct** → Generates mapper interfaces
- **Validation** → Adds `@Valid` annotations
- **Spring Data JPA** → Configures repositories correctly

### 4. Architectural Consistency

Every generated file follows the same patterns:

- Consistent naming conventions
- Proper layered architecture
- DTOs for API boundaries
- Exception handling

### 5. CLI-First

Works in any terminal:

```bash
# Interactive mode
haft init

# Scripted mode (CI/CD friendly)
haft init my-app --group com.example --deps web,jpa --no-interactive
```

## Comparison

| Feature | Spring Initializr | Haft |
|---------|------------------|------|
| Project scaffolding | ✅ | ✅ |
| Web interface | ✅ | ❌ |
| CLI interface | ❌ | ✅ |
| Interactive TUI | ❌ | ✅ |
| Resource generation | ❌ | ✅ |
| Dependency management | ❌ | ✅ |
| Smart detection | ❌ | ✅ |
| Works offline | ❌ | ✅ |

## Use Cases

### Starting a New Project

```bash
haft init my-api
```

Interactive wizard guides you through all configuration options.

### Adding a New Resource

```bash
haft generate resource Product
```

Generates all CRUD layers with proper architecture.

### Adding a Dependency

```bash
haft add security
```

Adds Spring Security with proper configuration.

### CI/CD Pipelines

```bash
haft init user-service \
  --group com.company \
  --java 21 \
  --spring 3.4.1 \
  --deps web,data-jpa,security \
  --no-interactive
```

Fully scriptable, no interaction required.

## Philosophy

1. **Convention over Configuration** — Sensible defaults, override when needed
2. **CLI is Documentation** — `--help` tells you everything
3. **Offline First** — No internet required after installation
4. **Zero Lock-in** — Generated code is standard Spring Boot
