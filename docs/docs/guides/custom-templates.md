# Custom Templates

Haft supports custom code generation templates, allowing you to tailor generated code to match your project's coding standards, conventions, and requirements.

## Template Priority

When generating code, Haft looks for templates in the following order:

1. **Project-level** (`.haft/templates/`) - highest priority
2. **Global user-level** (`~/.haft/templates/`)
3. **Built-in embedded templates** - fallback

This allows you to:
- Override templates at the project level for project-specific customizations
- Set global defaults in your home directory for consistent personal preferences
- Fall back to Haft's built-in templates for anything you haven't customized

## Initializing Custom Templates

To start customizing templates, use the `haft template init` command:

```bash
# Initialize all templates for customization
haft template init

# Initialize only resource templates
haft template init --category resource

# Initialize only test templates
haft template init --category test

# Initialize templates globally (user-level)
haft template init --global

# Overwrite existing templates
haft template init --force
```

This copies the built-in templates to your project's `.haft/templates/` directory, where you can modify them.

## Listing Templates

To see all available templates and where they come from:

```bash
# List all templates
haft template list

# List only custom (overridden) templates
haft template list --custom

# List templates in a specific category
haft template list --category resource

# Show full paths
haft template list --paths
```

## Template Categories

Templates are organized into categories:

| Category | Description |
|----------|-------------|
| `resource` | Resource generation templates (Controller, Service, Entity, etc.) |
| `test` | Test file templates |
| `project` | Project scaffolding templates (for `haft init`) |

## Template Structure

Templates follow this directory structure:

```
.haft/templates/
├── resource/
│   ├── layered/
│   │   ├── Controller.java.tmpl
│   │   ├── Service.java.tmpl
│   │   ├── ServiceImpl.java.tmpl
│   │   ├── Repository.java.tmpl
│   │   ├── Entity.java.tmpl
│   │   ├── Request.java.tmpl
│   │   ├── Response.java.tmpl
│   │   └── Mapper.java.tmpl
│   └── feature/
│       └── ... (same files)
├── test/
│   ├── layered/
│   │   ├── ServiceTest.java.tmpl
│   │   ├── ControllerTest.java.tmpl
│   │   └── ...
│   └── feature/
│       └── ...
└── project/
    ├── Application.java.tmpl
    ├── pom.xml.tmpl
    └── ...
```

## Template Variables

Templates use Go's `text/template` syntax. The following variables are available:

### Common Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `{{.Name}}` | Resource name (PascalCase) | `User` |
| `{{.NameLower}}` | Resource name (lowercase) | `user` |
| `{{.NameCamel}}` | Resource name (camelCase) | `user` |
| `{{.BasePackage}}` | Base package path | `com.example.app` |
| `{{.Package}}` | Full package path | `com.example.app.controller` |
| `{{.HasLombok}}` | Whether Lombok is available | `true` |
| `{{.HasJpa}}` | Whether JPA is available | `true` |
| `{{.HasValidation}}` | Whether validation is available | `true` |

### Template Functions

| Function | Description | Example |
|----------|-------------|---------|
| `lower` | Convert to lowercase | `{{.Name \| lower}}` |
| `upper` | Convert to uppercase | `{{.Name \| upper}}` |
| `title` | Convert to title case | `{{.Name \| title}}` |
| `camelCase` | Convert to camelCase | `{{.Name \| camelCase}}` |
| `pascalCase` | Convert to PascalCase | `{{.Name \| pascalCase}}` |
| `snakeCase` | Convert to snake_case | `{{.Name \| snakeCase}}` |
| `kebabCase` | Convert to kebab-case | `{{.Name \| kebabCase}}` |
| `plural` | Pluralize a word | `{{.Name \| plural}}` |
| `singular` | Singularize a word | `{{.Name \| singular}}` |

## Common Customizations

### Adding Copyright Headers

Add a copyright header to all generated files:

```java
/*
 * Copyright (c) {{.Year}} MyCompany Inc.
 * All rights reserved.
 */
package {{.Package}};
// ... rest of template
```

### Adding Custom Annotations

Add project-specific annotations:

```java
@Audited
@Cacheable
@RestController
@RequestMapping("/api/{{.NameLower | plural}}")
public class {{.Name}}Controller {
    // ...
}
```

### Using Java Records for DTOs

Replace class-based DTOs with records:

```java
package {{.Package}};

public record {{.Name}}Request(
    {{if .HasValidation}}@NotBlank {{end}}String name,
    {{if .HasValidation}}@Email {{end}}String email
) {}
```

### Custom Repository Methods

Add common repository methods:

```java
@Repository
public interface {{.Name}}Repository extends JpaRepository<{{.Name}}, {{.IDType}}> {
    
    Optional<{{.Name}}> findByEmail(String email);
    
    List<{{.Name}}> findByActiveTrue();
    
    @Query("SELECT e FROM {{.Name}} e WHERE e.createdAt > :date")
    List<{{.Name}}> findRecentlyCreated(@Param("date") LocalDateTime date);
}
```

## Tips

1. **Start Small**: Begin by customizing just one template to understand the process
2. **Test Changes**: After modifying a template, generate a test resource to verify
3. **Use Version Control**: Track your custom templates in `.haft/templates/` with git
4. **Share Templates**: Copy your project templates to `~/.haft/templates/` for use across projects
