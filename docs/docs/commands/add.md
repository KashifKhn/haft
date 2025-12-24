---
sidebar_position: 3
title: haft add
description: Add dependencies to your project
---

# haft add

Add dependencies to an existing Spring Boot project.

:::caution Coming Soon
The `add` command is currently under development and will be available in a future release.
:::

## Usage

```bash
haft add <dependency> [flags]
```

## Description

The `add` command modifies your `pom.xml` (or `build.gradle`) to add new dependencies. It automatically:

- Detects the correct dependency coordinates
- Adds required version properties
- Includes related dependencies when needed

## Examples

### Add a Single Dependency

```bash
haft add lombok
```

### Add Multiple Dependencies

```bash
haft add lombok validation mapstruct
```

### Add with Search

```bash
# Don't know the exact name? Search!
haft add --search security
```

## Common Dependencies

| Alias | Dependency |
|-------|------------|
| `web` | spring-boot-starter-web |
| `jpa` | spring-boot-starter-data-jpa |
| `security` | spring-boot-starter-security |
| `validation` | spring-boot-starter-validation |
| `lombok` | lombok |
| `mapstruct` | mapstruct + mapstruct-processor |
| `actuator` | spring-boot-starter-actuator |
| `devtools` | spring-boot-devtools |

## Flags

| Flag | Description |
|------|-------------|
| `--search` | Search for dependencies by name |
| `--scope` | Set dependency scope (compile, test, provided) |
| `--version` | Override default version |

## Database Drivers

```bash
# PostgreSQL
haft add postgresql

# MySQL
haft add mysql

# H2 (test database)
haft add h2
```

## What Gets Added

### Example: `haft add lombok`

```xml
<dependency>
    <groupId>org.projectlombok</groupId>
    <artifactId>lombok</artifactId>
    <scope>provided</scope>
</dependency>
```

### Example: `haft add mapstruct`

```xml
<dependency>
    <groupId>org.mapstruct</groupId>
    <artifactId>mapstruct</artifactId>
    <version>${mapstruct.version}</version>
</dependency>

<!-- Plus annotation processor configuration -->
```

## See Also

- [Dependencies](/docs/guides/dependencies) — Full dependency list
- [haft init](/docs/commands/init) — Add dependencies at project creation
