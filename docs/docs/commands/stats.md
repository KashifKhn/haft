---
sidebar_position: 9
title: haft stats
description: Display code statistics and metrics
---

# haft stats

Display detailed code statistics for your Spring Boot project.

## Usage

```bash
haft stats [flags]
```

## Description

The `stats` command uses [SCC (Sloc Cloc and Code)](https://github.com/boyter/scc) to analyze your codebase and provide detailed statistics including lines of code, comments, blank lines, and complexity metrics for each language in your project.

## Flags

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON format |
| `--cocomo` | Include COCOMO cost estimates |

## Examples

```bash
# Show code statistics
haft stats

# Output as JSON
haft stats --json

# Include COCOMO cost estimates
haft stats --cocomo

# JSON with COCOMO
haft stats --json --cocomo
```

## Output Columns

| Column | Description |
|--------|-------------|
| Language | Programming language detected |
| Files | Number of files |
| Lines | Total lines (code + comments + blanks) |
| Code | Lines of actual code |
| Comments | Comment lines |
| Blanks | Empty/whitespace lines |

## Sample Output

```
  Code Statistics
─────────────────────────────────────────────────────────────────────────────────────
  Language                  Files      Lines       Code   Comments     Blanks
─────────────────────────────────────────────────────────────────────────────────────
  Java                         19        534        378          0        156
  XML                           1         80         69          0         11
  YAML                          2         39         37          0          2
  Properties File               1          3          3          0          0
─────────────────────────────────────────────────────────────────────────────────────
  Total                        23        656        487          0        169

  Processed: 35.56 KB
```

## COCOMO Estimates

With `--cocomo` flag, the command includes software cost estimates based on the COCOMO (Constructive Cost Model) methodology:

```
  COCOMO Estimates
─────────────────────────────────────────────────────────────────────────────────────
  Estimated Cost:    $24,608
  Schedule Effort:   3.37 months
  People Required:   0.65
```

### COCOMO Metrics

| Metric | Description |
|--------|-------------|
| Estimated Cost | Development cost estimate in USD |
| Schedule Effort | Estimated development time in months |
| People Required | Average team size needed |

Note: COCOMO estimates use industry-standard parameters and should be used as rough guidelines, not precise predictions.

## JSON Output

With `--json` flag:

```json
{
  "languages": [
    {
      "name": "Java",
      "files": 19,
      "lines": 534,
      "code": 378,
      "comments": 0,
      "blanks": 156,
      "complexity": 18
    },
    {
      "name": "XML",
      "files": 1,
      "lines": 80,
      "code": 69,
      "comments": 0,
      "blanks": 11
    }
  ],
  "totalFiles": 20,
  "totalLines": 614,
  "totalCode": 447,
  "totalComments": 0,
  "totalBlanks": 167,
  "totalBytes": 35560
}
```

With `--json --cocomo`:

```json
{
  "languages": [...],
  "totalFiles": 20,
  "totalLines": 614,
  "totalCode": 447,
  "totalComments": 0,
  "totalBlanks": 167,
  "totalBytes": 35560,
  "estimatedCost": 24608.50,
  "estimatedMonths": 3.37,
  "estimatedPeople": 0.65
}
```

## Excluded Directories

The following directories are automatically excluded from analysis:

- `.git`
- `.svn`
- `.hg`
- `node_modules`
- `target` (Maven build output)
- `build` (Gradle build output)
- `.gradle`
- `.idea`

## Supported Languages

SCC supports 200+ programming languages. Common languages in Spring Boot projects:

- Java
- Kotlin
- XML (pom.xml, configuration files)
- YAML (application.yml)
- Properties (application.properties)
- SQL
- HTML/CSS/JavaScript (if web resources present)

## Quick Stats with haft info

For a quick lines of code summary without the full breakdown, use:

```bash
haft info --loc
```

This adds a condensed code statistics section to the project info output.

## Use Cases

- Track codebase growth over time
- Estimate project complexity
- Generate reports for stakeholders
- Compare language distribution
- Identify documentation coverage (comments ratio)

## Editor Integration

Use this command from your editor:

- **Neovim**: `:HaftStats` ([docs →](/docs/integrations/neovim/usage#project-information-commands))
- **VS Code**: Coming soon ([preview →](/docs/integrations/vscode))
- **IntelliJ IDEA**: Coming soon ([preview →](/docs/integrations/intellij))

## See Also

- [haft info](/docs/commands/info) - Project information with optional `--loc` flag
- [haft routes](/docs/commands/routes) - List REST API endpoints
- [SCC Documentation](https://github.com/boyter/scc) - Underlying statistics tool
