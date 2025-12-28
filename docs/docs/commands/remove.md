---
sidebar_position: 4
title: haft remove
description: Remove dependencies from your project
---

# haft remove

Remove dependencies from an existing Spring Boot project.

## Usage

```bash
haft remove                           # Interactive picker
haft remove <dependency> [dependencies...]
haft remove <groupId:artifactId>
```

## Aliases

```bash
haft rm                               # Short alias
```

## Description

The `remove` command modifies your build file (`pom.xml` or `build.gradle`) to remove dependencies. It supports:

- **Interactive mode** — Select dependencies to remove from a list
- **By artifact name** — Remove by artifact ID (e.g., `lombok`)
- **By coordinates** — Remove by full coordinates (e.g., `org.projectlombok:lombok`)
- **Suffix matching** — `jpa` matches `spring-boot-starter-jpa`
- **Multiple removal** — Remove several dependencies at once

## Interactive Mode

```bash
haft remove
# or
haft rm
```

Opens an interactive TUI showing all current dependencies where you can:
- Type to search/filter
- Select multiple with `Space`
- Navigate with `↑`/`↓`, `PgUp`/`PgDown`
- Select all with `a`, none with `n`
- Confirm removal with `Enter`
- Cancel with `Esc`

The picker shows:
- Artifact ID (highlighted)
- Group ID
- Version (if specified)
- Scope (if not compile)

## Flags

| Flag | Description |
|------|-------------|
| `--json` | Output result as JSON |
| `--no-interactive` | Skip interactive picker (requires dependency argument) |

## Examples

### Remove by Artifact Name

```bash
# Remove Lombok
haft remove lombok

# Remove Spring Web
haft remove spring-boot-starter-web
```

### Remove by Coordinates

```bash
haft remove org.projectlombok:lombok
haft remove io.jsonwebtoken:jjwt-api
```

### Remove Multiple Dependencies

```bash
haft remove lombok validation h2
haft rm jjwt-api jjwt-impl jjwt-jackson
```

### Suffix Matching

Haft supports suffix matching for convenience:

```bash
# These are equivalent:
haft remove spring-boot-starter-jpa
haft remove jpa

# Matches spring-boot-starter-web
haft remove web
```

## Output

### Successful Removal

```
$ haft remove lombok validation
INFO ✓ Removed dependency=org.projectlombok:lombok
INFO ✓ Removed dependency=org.springframework.boot:spring-boot-starter-validation
INFO ✓ Removed 2 dependencies from pom.xml
```

### Dependency Not Found

```
$ haft remove nonexistent
WARN ⚠ Dependency not found input=nonexistent
INFO ℹ No dependencies removed (none found)
```

### No Dependencies in Project

```
$ haft remove
INFO ℹ No dependencies found in pom.xml
```

## Keyboard Shortcuts (Interactive Mode)

| Key | Action |
|-----|--------|
| `↑` `↓` | Navigate |
| `Space` | Toggle selection |
| `a` | Select all visible |
| `n` | Select none |
| `Enter` | Confirm removal |
| `Esc` | Cancel |
| Type | Filter dependencies |
| `Backspace` | Clear filter |

## How It Works

1. Parses your build file (`pom.xml` or `build.gradle`) to find all dependencies
2. In interactive mode, displays them for selection
3. In CLI mode, matches input against:
   - Full coordinates (`groupId:artifactId`)
   - Exact artifact ID match
   - Suffix match (for Spring starters)
4. Removes matched dependencies from the build file
5. Writes the updated file

## Build Tool Detection

Haft automatically detects your build tool:

| File Found | Build Tool |
|------------|------------|
| `pom.xml` | Maven |
| `build.gradle.kts` | Gradle (Kotlin DSL) |
| `build.gradle` | Gradle (Groovy DSL) |

## Tips

### Removing JWT Dependencies

If you added JWT with `haft add jwt`, it added 3 artifacts. Remove them individually:

```bash
haft rm jjwt-api jjwt-impl jjwt-jackson
```

Or use the interactive picker to select all three.

### Removing Spring Starters

For Spring Boot starters, you can use the short name:

```bash
haft rm web jpa security
# Removes:
# - spring-boot-starter-web
# - spring-boot-starter-data-jpa
# - spring-boot-starter-security
```

### Double-Check with Git

Since `haft remove` modifies your build file, use Git to review changes:

```bash
haft rm lombok
git diff pom.xml      # Maven
git diff build.gradle # Gradle
```

## See Also

- [haft add](/docs/commands/add) — Add dependencies
- [Dependencies Guide](/docs/guides/dependencies) — Full dependency reference
