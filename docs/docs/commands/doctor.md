---
sidebar_position: 8
title: haft doctor
description: Check project health and best practices
---

# haft doctor

Analyze your Spring Boot project for issues, warnings, and suggestions.

## Usage

```bash
haft doctor [flags]
haft doc [flags]     # Alias
haft check [flags]   # Alias
haft health [flags]  # Alias
```

## Description

The `doctor` command performs comprehensive health checks on your Spring Boot project. It identifies problems early and suggests improvements to help you maintain a healthy codebase.

## Health Checks

### Build Checks

| Check | Severity | Description |
|-------|----------|-------------|
| `build_file` | Error | Build file exists (pom.xml or build.gradle) |
| `spring_boot_config` | Error | Spring Boot parent/plugin is configured |
| `java_version` | Warning | Java version is specified (warns for Java 8) |

### Source Checks

| Check | Severity | Description |
|-------|----------|-------------|
| `source_directory` | Error | Source directory exists (src/main/java or src/main/kotlin) |
| `test_directory` | Warning | Test directory exists (src/test/java or src/test/kotlin) |
| `main_class` | Error | Main class with @SpringBootApplication found |

### Configuration Checks

| Check | Severity | Description |
|-------|----------|-------------|
| `config_file` | Warning | Configuration file exists (application.yml/properties) |
| `hardcoded_secrets` | Error | No hardcoded passwords, secrets, or API keys |

### Security Checks

| Check | Severity | Description |
|-------|----------|-------------|
| `h2_scope` | Warning | H2 database not in compile scope |
| `devtools_scope` | Warning | DevTools marked as optional |

### Dependency Checks

| Check | Severity | Description |
|-------|----------|-------------|
| `test_dependencies` | Warning | spring-boot-starter-test is present |
| `suggest_actuator` | Suggestion | Consider adding Actuator for monitoring |
| `suggest_validation` | Suggestion | Consider adding Validation for input validation |
| `suggest_openapi` | Suggestion | Consider adding OpenAPI for API documentation |

### Best Practice Checks

| Check | Severity | Description |
|-------|----------|-------------|
| `lombok_config` | Info | Lombok configured with lombok.config |

## Flags

| Flag | Description |
|------|-------------|
| `--json` | Output results as JSON (for CI/CD pipelines) |
| `--strict` | Exit with code 1 on any warning |
| `--category` | Filter by category (build, source, config, security, dependencies, best-practice) |

## Examples

```bash
# Run full health check
haft doctor

# Output as JSON (for CI/CD)
haft doctor --json

# Strict mode (exit 1 on warnings)
haft doctor --strict

# Filter by category
haft doctor --category security
haft doctor --category build
haft doctor --category dependencies
```

## Output Format

### Standard Output

```
🏥 Haft Doctor - Project Health Check
=============================================

Project: my-app
Path: /path/to/my-app
Build Tool: Maven

Passed Checks:
  ✓ Build file exists (Maven)
  ✓ Spring Boot configured
  ✓ Java 17 configured
  ✓ Source directory exists (Java)
  ✓ Main class with @SpringBootApplication found
  ✓ Configuration file found (application.yml)
  ✓ No hardcoded secrets detected
  ✓ Test dependencies configured

Warnings:
  ⚠ H2 database in compile scope
    H2 should be test or runtime scope only
    → Change H2 scope to <scope>runtime</scope> or testImplementation

Suggestions:
  💡 Consider adding Actuator
    Actuator provides health checks, metrics, and monitoring endpoints
    → Run: haft add actuator

---------------------------------------------
Summary: 8 passed, 1 warnings, 2 suggestions
```

### JSON Output

```json
{
  "project_path": "/path/to/my-app",
  "project_name": "my-app",
  "build_tool": "Maven",
  "results": [
    {
      "name": "build_file",
      "category": "build",
      "passed": true,
      "severity": "error",
      "message": "Build file exists (Maven)"
    },
    {
      "name": "h2_scope",
      "category": "security",
      "passed": false,
      "severity": "warning",
      "message": "H2 database in compile scope",
      "details": "H2 should be test or runtime scope only",
      "fix_hint": "Change H2 scope to <scope>runtime</scope>"
    }
  ],
  "passed_count": 8,
  "error_count": 0,
  "warning_count": 1,
  "suggestion_count": 2
}
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All checks passed (or only warnings/suggestions in non-strict mode) |
| 1 | Health check failed (errors found, or warnings in strict mode) |

## Use Cases

### Pre-commit Hook

```bash
#!/bin/bash
haft doctor --strict || exit 1
```

### CI Pipeline

```yaml
# GitHub Actions
- name: Check project health
  run: haft doctor --json > health-report.json

- name: Fail on issues
  run: haft doctor --strict
```

### Quick Security Check

```bash
haft doctor --category security
```

### Check Before Release

```bash
haft doctor --strict
```

## Categories

Filter checks by category using `--category`:

- `build` - Build file and configuration checks
- `source` - Source code structure checks
- `config` - Configuration file checks
- `security` - Security-related checks
- `dependencies` - Dependency recommendations
- `best-practice` - Best practice suggestions

## Editor Integration

Use this command from your editor:

- **Neovim**: Not yet available in haft.nvim (CLI only)
- **VS Code**: Coming soon ([preview →](/docs/integrations/vscode))
- **IntelliJ IDEA**: Coming soon ([preview →](/docs/integrations/intellij))

## See Also

- [haft dev validate](/docs/commands/dev#haft-dev-validate) - Build tool validation
- [haft info](/docs/commands/info) - Show project information
- [haft add](/docs/commands/add) - Add dependencies
