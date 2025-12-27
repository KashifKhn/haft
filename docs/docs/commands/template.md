---
sidebar_position: 8
title: haft template
description: Manage custom code generation templates
---

# haft template

Manage custom templates for code generation.

## Usage

```bash
haft template <subcommand> [flags]
```

## Subcommands

| Command | Description |
|---------|-------------|
| `haft template init` | Copy embedded templates for customization |
| `haft template list` | List all available templates with sources |
| `haft template validate` | Validate custom template syntax |

## Template Locations

Haft loads templates from multiple locations with the following priority:

| Priority | Location | Description |
|----------|----------|-------------|
| 1 (Highest) | `.haft/templates/` | Project-specific templates |
| 2 | `~/.haft/templates/` | User-global templates |
| 3 (Lowest) | Built-in | Embedded default templates |

When a template exists in multiple locations, only the highest priority source is used during code generation.

---

## haft template init

Copy embedded templates to a local directory for customization.

```bash
# Copy all templates to project (.haft/templates/)
haft template init

# Copy specific category
haft template init --category resource

# Copy to global location (~/.haft/templates/)
haft template init --global

# Force overwrite existing templates
haft template init --force
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--category` | `-c` | Template category to copy: `resource`, `test`, `project` |
| `--global` | `-g` | Copy to global directory (`~/.haft/templates/`) instead of project |
| `--force` | `-f` | Overwrite existing templates |

### Examples

```bash
# Copy all templates for project customization
haft template init

# Copy only resource templates (Controller, Service, Entity, etc.)
haft template init --category resource

# Copy only test templates
haft template init --category test

# Set up global templates for all projects
haft template init --global

# Update templates to latest version
haft template init --force
```

### Template Categories

| Category | Templates |
|----------|-----------|
| `resource` | Controller, Service, ServiceImpl, Repository, Entity, Request, Response, Mapper |
| `test` | ServiceTest, ControllerTest, RepositoryTest, EntityTest |
| `project` | Application.java, pom.xml, build.gradle, application.yml |

---

## haft template list

List all available templates and their sources.

```bash
# List all templates
haft template list

# Show only custom templates
haft template list --custom

# Filter by category
haft template list --category resource

# Show full paths
haft template list --paths
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--custom` | | Show only custom (non-embedded) templates |
| `--category` | `-c` | Filter by category: `resource`, `test`, `project` |
| `--paths` | | Show full template file paths |

### Example Output

```
  Available Templates
────────────────────────────────────────────────────────

  Resource
    resource/layered/Controller.java.tmpl [project]
    resource/layered/Service.java.tmpl [embedded]
    resource/layered/Entity.java.tmpl [global]

  Test
    test/layered/ServiceTest.java.tmpl [embedded]
    test/layered/ControllerTest.java.tmpl [embedded]

  Templates: 25 total (project: 1) (global: 1) (embedded: 23)
```

Source indicators:
- `[project]` — Template from `.haft/templates/`
- `[global]` — Template from `~/.haft/templates/`
- `[embedded]` — Built-in default template

---

## haft template validate

Validate custom template syntax and check for errors.

```bash
# Validate all project templates
haft template validate

# Show available placeholder variables
haft template validate --vars

# Show available conditions for @if directives
haft template validate --conditions
```

### Flags

| Flag | Description |
|------|-------------|
| `--vars` | Display all available placeholder variables |
| `--conditions` | Display all available conditions for `@if` directives |

### Example: Show Variables

```bash
$ haft template validate --vars

Available Template Variables:
─────────────────────────────
  ${Name}         Resource name (PascalCase)      → User
  ${name}         Resource name (lowercase)       → user
  ${nameCamel}    Resource name (camelCase)       → user
  ${nameSnake}    Resource name (snake_case)      → user
  ${nameKebab}    Resource name (kebab-case)      → user
  ${namePlural}   Pluralized lowercase name       → users
  ${NamePlural}   Pluralized PascalCase name      → Users
  ${BasePackage}  Base package path               → com.example.app
  ${Package}      Full package path               → com.example.app.user
  ${IDType}       Entity ID type                  → Long or UUID
  ${TableName}    Database table name             → users
```

### Example: Show Conditions

```bash
$ haft template validate --conditions

Available Conditions:
─────────────────────
  HasLombok       Lombok dependency detected
  HasJpa          Spring Data JPA detected
  HasValidation   Bean Validation detected
  HasMapStruct    MapStruct mapper detected
  HasSwagger      Swagger/OpenAPI detected
  HasBaseEntity   Base entity class detected
  UsesUUID        UUID used for entity IDs
  UsesLong        Long used for entity IDs
```

### Example: Validate Templates

```bash
$ haft template validate

Validating templates in .haft/templates/...

✓ resource/layered/Controller.java.tmpl
✓ resource/layered/Service.java.tmpl
✗ resource/layered/Entity.java.tmpl
    Line 5: Unknown placeholder ${InvalidVar}
    Line 12: Unclosed @if directive

Validation: 2 passed, 1 failed
```

---

## Template Syntax

Haft supports two template syntaxes:

### 1. Simple Placeholder Syntax

Use `${variable}` for easy templating:

```java
package ${BasePackage}.controller;

@RestController
@RequestMapping("/api/${namePlural}")
public class ${Name}Controller {
    private final ${Name}Service ${nameCamel}Service;
}
```

### 2. Comment-Based Conditionals

Use `// @if`, `// @else`, `// @endif` for conditional blocks:

```java
// @if HasLombok
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
// @else
// TODO: Add getters and setters manually
// @endif
public class ${Name} {
    // @if UsesUUID
    private UUID id;
    // @else
    private Long id;
    // @endif
}
```

### 3. Standard Go Templates

Full Go template syntax is also supported:

```java
package {{.BasePackage}}.entity;

{{if .HasLombok}}
import lombok.*;
{{end}}

@Entity
@Table(name = "{{plural .NameLower}}")
public class {{.Name}} {
    @Id
    private {{.IDType}} id;
}
```

---

## Workflow Example

```bash
# 1. Initialize custom templates
haft template init --category resource

# 2. Edit templates to match your conventions
vim .haft/templates/resource/layered/Controller.java.tmpl

# 3. Validate your changes
haft template validate

# 4. Generate code using custom templates
haft generate resource User
```

## See Also

- [Templates Reference](/docs/reference/templates) — Template variables and functions
- [Custom Templates Guide](/docs/guides/custom-templates) — Detailed customization guide
- [haft generate](/docs/commands/generate) — Code generation commands
