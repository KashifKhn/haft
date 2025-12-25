---
sidebar_position: 1
title: Configuration
description: Haft configuration options
---

# Configuration

Haft can be configured through command-line flags, environment variables, and configuration files.

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

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `HAFT_NO_COLOR` | Disable colors | `false` |
| `HAFT_VERBOSE` | Enable verbose mode | `false` |

### Example

```bash
export HAFT_NO_COLOR=true
haft init
```

## Default Values

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
  --spring 3.4.1 \
  --build maven \
  --packaging jar \
  --config yaml \
  --deps web,jpa,lombok \
  --no-interactive
```

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
