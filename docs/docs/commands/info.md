---
sidebar_position: 7
title: haft info
description: Display project information and metadata
---

# haft info

Display detailed information about the current Spring Boot project.

## Usage

```bash
haft info [flags]
```

## Description

The `info` command analyzes your Spring Boot project and displays comprehensive information including project metadata, build configuration, dependency summary, and optionally code statistics.

## Flags

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON format |
| `--loc` | Include lines of code summary |
| `--deps` | Include full dependency list |

## Examples

```bash
# Show project info
haft info

# Output as JSON
haft info --json

# Include lines of code summary
haft info --loc

# Include full dependency list
haft info --deps

# JSON output with code stats
haft info --json --loc

# Full project analysis as JSON (useful for editor plugins)
haft info --json --loc --deps
```

## Output Sections

### Project Information

Displays basic project metadata:

- **Name**: Directory name of the project
- **Group ID**: Maven/Gradle group identifier
- **Artifact ID**: Maven/Gradle artifact identifier
- **Version**: Project version
- **Description**: Project description (if set)

### Build Configuration

Shows build tool and Java configuration:

- **Build Tool**: Maven or Gradle
- **Build File**: pom.xml or build.gradle
- **Java Version**: Target Java version
- **Spring Boot**: Spring Boot version
- **Packaging**: JAR or WAR (if specified)

### Dependencies

Summary of project dependencies:

- **Total**: Total number of dependencies
- **Spring Starters**: Count of spring-boot-starter-* dependencies
- **Spring Libraries**: Other Spring dependencies
- **Test Dependencies**: Dependencies with test scope

### Key Dependencies

Shows presence of common dependencies:

- Spring Web
- Spring Data JPA
- Lombok
- Validation
- MapStruct
- Spring Security

### Code Statistics (with --loc)

When `--loc` flag is used, displays:

- **Total Files**: Number of source files
- **Lines of Code**: Actual code lines (excluding comments/blanks)
- **Comments**: Comment lines
- **Blank Lines**: Empty lines

## Sample Output

```
  Project Information
──────────────────────────────────────────────────
  Name:              demo
  Group ID:          com.example
  Artifact ID:       demo
  Version:           0.0.1-SNAPSHOT
  Description:       Demo Spring Boot application

  Build Configuration
──────────────────────────────────────────────────
  Build Tool:        Maven
  Build File:        pom.xml
  Java Version:      17
  Spring Boot:       3.2.0

  Dependencies
──────────────────────────────────────────────────
  Total:             7
  Spring Starters:   5
  Spring Libraries:  0
  Test Dependencies: 1

  Key Dependencies
──────────────────────────────────────────────────
  Spring Web:        ✓
  Spring Data JPA:   ✓
  Lombok:            ✓
  Validation:        ✓
  MapStruct:         –
  Spring Security:   –
```

## JSON Output

With `--json` flag:

```json
{
  "name": "demo",
  "groupId": "com.example",
  "artifactId": "demo",
  "version": "0.0.1-SNAPSHOT",
  "description": "Demo Spring Boot application",
  "buildTool": "maven",
  "buildFile": "pom.xml",
  "javaVersion": "17",
  "springBootVersion": "3.2.0",
  "packaging": "",
  "dependencyCount": 7
}
```

With `--json --loc --deps`:

```json
{
  "name": "demo",
  "groupId": "com.example",
  "artifactId": "demo",
  "version": "0.0.1-SNAPSHOT",
  "description": "Demo Spring Boot application",
  "buildTool": "maven",
  "buildFile": "pom.xml",
  "javaVersion": "17",
  "springBootVersion": "3.2.0",
  "packaging": "",
  "dependencyCount": 7,
  "codeStats": {
    "totalFiles": 25,
    "linesOfCode": 915,
    "comments": 0,
    "blankLines": 225
  },
  "dependencies": [
    {
      "groupId": "org.springframework.boot",
      "artifactId": "spring-boot-starter-web",
      "version": "",
      "scope": ""
    },
    {
      "groupId": "org.springframework.boot",
      "artifactId": "spring-boot-starter-data-jpa",
      "version": "",
      "scope": ""
    },
    {
      "groupId": "org.projectlombok",
      "artifactId": "lombok",
      "version": "",
      "scope": ""
    }
  ]
}
```

## Editor Integration

Use this command from your editor:

- **Neovim**: `:HaftInfo` ([docs →](/docs/integrations/neovim/usage#project-information-commands))
- **VS Code**: Coming soon ([preview →](/docs/integrations/vscode))
- **IntelliJ IDEA**: Coming soon ([preview →](/docs/integrations/intellij))

## See Also

- [haft stats](/docs/commands/stats) - Full code statistics with language breakdown
- [haft routes](/docs/commands/routes) - List REST API endpoints
- [haft dev](/docs/commands/dev) - Development commands
