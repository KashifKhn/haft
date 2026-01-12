---
sidebar_position: 1
title: Configuration
description: Haft configuration options and config files
---

# Configuration

Haft can be configured through command-line flags and configuration files.

## Configuration Files

Haft supports two configuration files:

| Config Type | File Name | Location | Purpose |
|-------------|-----------|----------|---------|
| **Project Config** | `.haft.json` | Project root directory | Project-specific settings |
| **Global Config** | `config.json` | `~/.config/haft/` | User-wide defaults |

### Project Configuration (`.haft.json`)

Create a `.haft.json` file in your project root to configure project-specific settings:

```json
{
  "version": "1",
  "project": {
    "name": "my-app",
    "group": "com.example",
    "artifact": "my-app",
    "description": "My Spring Boot application",
    "package": "com.example.myapp"
  },
  "spring": {
    "version": "3.4.0"
  },
  "java": {
    "version": "21"
  },
  "build": {
    "tool": "maven"
  },
  "architecture": {
    "style": "layered"
  },
  "database": {
    "type": "postgresql"
  },
  "generators": {
    "dto": {
      "style": "record"
    },
    "tests": {
      "enabled": true
    }
  }
}
```

#### Project Config Reference

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `version` | string | `"1"` | Config file version |
| `project.name` | string | `""` | Project name |
| `project.group` | string | `""` | Maven group ID |
| `project.artifact` | string | `""` | Maven artifact ID |
| `project.description` | string | `""` | Project description |
| `project.package` | string | `""` | Base Java package |
| `spring.version` | string | `"3.4.0"` | Spring Boot version |
| `java.version` | string | `"21"` | Java version |
| `build.tool` | string | `"maven"` | Build tool (maven, gradle, gradle-kotlin) |
| `architecture.style` | string | `"layered"` | Architecture style |
| `database.type` | string | `"postgresql"` | Database type |
| `generators.dto.style` | string | `"record"` | DTO generation style |
| `generators.tests.enabled` | bool | `true` | Generate test files |

### Global Configuration (`~/.config/haft/config.json`)

Create a global config file to set user-wide defaults:

```json
{
  "defaults": {
    "java_version": "21",
    "build_tool": "maven",
    "architecture": "layered",
    "spring_boot": "3.4.0"
  },
  "output": {
    "colors": true,
    "verbose": false
  }
}
```

#### Global Config Reference

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `defaults.java_version` | string | `"21"` | Default Java version for new projects |
| `defaults.build_tool` | string | `"maven"` | Default build tool |
| `defaults.architecture` | string | `"layered"` | Default architecture style |
| `defaults.spring_boot` | string | `"3.4.0"` | Default Spring Boot version |
| `output.colors` | bool | `true` | Enable colored output |
| `output.verbose` | bool | `false` | Enable verbose output |

## Command-Line Flags

### Global Flags

Available for all commands:

| Flag | Description |
|------|-------------|
| `--verbose` | Enable verbose output for debugging |
| `--no-color` | Disable colored output |
| `--help` | Show help for any command |
| `--version` | Show version information |

### Example

```bash
haft --verbose init my-app
haft --no-color init
```

## Default Values

When using the interactive wizard, these values are pre-selected:

### Project Defaults

| Field | Default Value |
|-------|---------------|
| Group ID | `com.example` |
| Version | `0.0.1-SNAPSHOT` |
| Java Version | `21` |
| Build Tool | `maven` |
| Packaging | `jar` |
| Config Format | `yaml` |

### Generation Defaults

| Setting | Default |
|---------|---------|
| Use Lombok | Auto-detect from build file |
| Use MapStruct | Auto-detect from build file |
| Use Validation | Auto-detect from build file |

## Overriding Defaults

### Per-Command

```bash
haft init my-app --java 17 --build gradle
```

### Non-Interactive Mode

Provide all required values via flags:

```bash
haft init my-app \
  --group com.company \
  --artifact my-service \
  --java 21 \
  --spring 3.4.0 \
  --build maven \
  --packaging jar \
  --config yaml \
  --deps web,jpa,lombok \
  --no-interactive
```

## Configuration Priority

Haft uses the following priority order (highest to lowest):

1. **Command-line flags** - Always take precedence
2. **Project config** (`.haft.json`) - Project-specific settings
3. **Global config** (`~/.config/haft/config.json`) - User defaults
4. **Built-in defaults** - Fallback values

## Dependency Detection

When generating code, Haft reads your build file (`pom.xml` or `build.gradle`) to detect:

### Lombok Detection

**Maven** - Looks for:
```xml
<dependency>
    <groupId>org.projectlombok</groupId>
    <artifactId>lombok</artifactId>
</dependency>
```

**Gradle** - Looks for:
```groovy
compileOnly 'org.projectlombok:lombok'
// or
annotationProcessor 'org.projectlombok:lombok'
```

**Effect**: Generated entities use `@Data`, `@Builder`, etc.

### MapStruct Detection

**Maven** - Looks for:
```xml
<dependency>
    <groupId>org.mapstruct</groupId>
    <artifactId>mapstruct</artifactId>
</dependency>
```

**Gradle** - Looks for:
```groovy
implementation 'org.mapstruct:mapstruct'
```

**Effect**: Generates mapper interfaces with `@Mapper` annotation.

### Validation Detection

**Maven** - Looks for:
```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-validation</artifactId>
</dependency>
```

**Gradle** - Looks for:
```groovy
implementation 'org.springframework.boot:spring-boot-starter-validation'
```

**Effect**: Controllers use `@Valid` on request parameters.

### Spring Data JPA Detection

**Maven** - Looks for:
```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-data-jpa</artifactId>
</dependency>
```

**Gradle** - Looks for:
```groovy
implementation 'org.springframework.boot:spring-boot-starter-data-jpa'
```

**Effect**: Repositories extend `JpaRepository`.

## CI/CD Configuration

For automated environments, always use `--no-interactive`:

```bash
# GitHub Actions example
- name: Generate project
  run: |
    haft init ${{ github.event.inputs.project_name }} \
      --group ${{ github.event.inputs.group_id }} \
      --java 21 \
      --deps web,jpa,actuator \
      --no-interactive
```

## Shell Completion

Enable shell completion for better CLI experience:

```bash
# Bash
echo 'source <(haft completion bash)' >> ~/.bashrc

# Zsh
echo 'source <(haft completion zsh)' >> ~/.zshrc

# Fish
haft completion fish > ~/.config/fish/completions/haft.fish
```
