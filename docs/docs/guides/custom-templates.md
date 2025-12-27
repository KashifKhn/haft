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

## Simple Template Syntax

Haft provides a user-friendly template syntax that's easy to learn:

### Variables

Use `${VariableName}` for simple variable substitution:

```java
package ${BasePackage}.controller;

public class ${Name}Controller {
    private final ${Name}Service ${nameCamel}Service;
    
    @GetMapping("/api/${namePlural}")
    public List<${Name}> getAll() {
        return ${nameCamel}Service.findAll();
    }
}
```

### Available Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `${Name}` | Resource name (PascalCase) | `User` |
| `${name}` | Resource name (lowercase) | `user` |
| `${nameCamel}` | Resource name (camelCase) | `user` |
| `${nameSnake}` | Resource name (snake_case) | `user_name` |
| `${nameKebab}` | Resource name (kebab-case) | `user-name` |
| `${namePlural}` | Pluralized lowercase name | `users` |
| `${NamePlural}` | Pluralized PascalCase name | `Users` |
| `${BasePackage}` | Base package path | `com.example.app` |
| `${Package}` | Full package path | `com.example.app.controller` |
| `${IDType}` | Entity ID type | `Long` or `UUID` |
| `${TableName}` | Database table name | `users` |

Run `haft template validate --vars` to see all available variables.

### Conditionals

Use comment-based conditionals with `// @if`, `// @else`, and `// @endif`:

```java
package ${BasePackage}.entity;

// @if HasLombok
import lombok.Data;
import lombok.Builder;
// @else
// Manual getters/setters required
// @endif

// @if HasLombok
@Data
@Builder
// @endif
public class ${Name} {
    
    // @if UsesUUID
    private UUID id;
    // @else
    private Long id;
    // @endif
    
    private String name;
}
```

### Available Conditions

| Condition | Description |
|-----------|-------------|
| `HasLombok` | True if Lombok is available |
| `HasJpa` | True if Spring Data JPA is available |
| `HasValidation` | True if Bean Validation is available |
| `HasMapStruct` | True if MapStruct is available |
| `HasSwagger` | True if Swagger/OpenAPI is available |
| `HasBaseEntity` | True if project has a base entity class |
| `UsesUUID` | True if entity uses UUID as ID type |
| `UsesLong` | True if entity uses Long as ID type |

Run `haft template validate --conditions` to see all available conditions.

## Managing Templates

### Initialize Custom Templates

Copy built-in templates to your project for customization:

```bash
# Initialize all templates
haft template init

# Initialize only resource templates
haft template init --category resource

# Initialize only test templates
haft template init --category test

# Initialize templates globally
haft template init --global

# Overwrite existing templates
haft template init --force
```

### List Templates

View all available templates and their sources:

```bash
# List all templates
haft template list

# List only custom templates
haft template list --custom

# List specific category
haft template list --category resource

# Show full paths
haft template list --paths
```

### Validate Templates

Check templates for errors before using them:

```bash
# Validate all project templates
haft template validate

# Validate a specific template
haft template validate .haft/templates/resource/layered/Controller.java.tmpl

# Validate a directory
haft template validate .haft/templates/resource/

# Show available variables
haft template validate --vars

# Show available conditions
haft template validate --conditions
```

Validation checks for:
- Unclosed placeholders (missing `}`)
- Unmatched `@if`/`@endif` directives
- Unknown variables (warnings)
- Template syntax errors

## Template Structure

Templates are organized by category:

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

## Advanced: Go Template Syntax

For complex logic, you can also use Go's `text/template` syntax directly:

```java
package {{.BasePackage}}.controller;

{{if .HasValidation}}
import jakarta.validation.Valid;
{{end}}

public class {{.Name}}Controller {
    {{range .Fields}}
    private {{.Type}} {{.Name}};
    {{end}}
}
```

The simple `${var}` syntax and Go template syntax can be mixed in the same file.

## Common Customizations

### Adding Copyright Headers

```java
/*
 * Copyright (c) 2024 MyCompany Inc.
 * All rights reserved.
 */
package ${BasePackage}.controller;
```

### Custom Annotations

```java
// @if HasLombok
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
// @endif
@Audited
@Cacheable
public class ${Name} {
```

### Using Records for DTOs

```java
package ${BasePackage}.dto;

// @if HasValidation
import jakarta.validation.constraints.*;
// @endif

public record ${Name}Request(
    // @if HasValidation
    @NotBlank
    // @endif
    String name,
    
    // @if HasValidation
    @Email
    // @endif
    String email
) {}
```

### Custom Repository Methods

```java
@Repository
public interface ${Name}Repository extends JpaRepository<${Name}, ${IDType}> {
    
    Optional<${Name}> findByEmail(String email);
    
    List<${Name}> findByActiveTrue();
    
    @Query("SELECT e FROM ${Name} e WHERE e.createdAt > :date")
    List<${Name}> findRecentlyCreated(@Param("date") LocalDateTime date);
}
```

## Tips

1. **Validate First**: Always run `haft template validate` after modifying templates
2. **Start Small**: Begin by customizing just one template to understand the process
3. **Use Version Control**: Track your custom templates in `.haft/templates/` with git
4. **Share Templates**: Copy project templates to `~/.haft/templates/` for use across projects
5. **Check Variables**: Use `haft template validate --vars` to see what's available
