---
sidebar_position: 2
title: Templates
description: Understanding and customizing Haft templates
---

# Templates

Haft uses Go templates to generate code. This page explains how templates work and how they might be customized in the future.

## Built-in Templates

Haft includes embedded templates for all generated files:

### Project Templates

| Template | Output |
|----------|--------|
| `Application.java.tmpl` | Main application class |
| `ApplicationTests.java.tmpl` | Application test class |
| `pom.xml.tmpl` | Maven POM file |
| `application.yml.tmpl` | YAML configuration |
| `application.properties.tmpl` | Properties configuration |
| `gitignore.tmpl` | Git ignore file |

### Resource Templates

| Template | Output |
|----------|--------|
| `Entity.java.tmpl` | JPA entity |
| `Repository.java.tmpl` | Spring Data repository |
| `Service.java.tmpl` | Service interface |
| `ServiceImpl.java.tmpl` | Service implementation |
| `Controller.java.tmpl` | REST controller |
| `Request.java.tmpl` | Request DTO |
| `Response.java.tmpl` | Response DTO |
| `Mapper.java.tmpl` | MapStruct mapper |
| `ResourceNotFoundException.java.tmpl` | Exception class |

## Template Variables

Templates have access to these variables:

### Project Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `{{.Name}}` | Project name | `MyApp` |
| `{{.GroupId}}` | Maven group ID | `com.example` |
| `{{.ArtifactId}}` | Maven artifact ID | `my-app` |
| `{{.Version}}` | Project version | `0.0.1-SNAPSHOT` |
| `{{.Description}}` | Project description | `My Spring Boot App` |
| `{{.BasePackage}}` | Base package | `com.example.myapp` |
| `{{.JavaVersion}}` | Java version | `21` |
| `{{.SpringBootVersion}}` | Spring Boot version | `3.4.1` |

### Resource Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `{{.Name}}` | Resource name (PascalCase) | `User` |
| `{{.NameLower}}` | Resource name (lowercase) | `user` |
| `{{.NameCamel}}` | Resource name (camelCase) | `user` |
| `{{.BasePackage}}` | Base package | `com.example.myapp` |

### Boolean Flags

| Variable | Description |
|----------|-------------|
| `{{.HasLombok}}` | Lombok is available |
| `{{.HasMapStruct}}` | MapStruct is available |
| `{{.HasValidation}}` | Validation is available |

### Template Functions

| Function | Description | Example |
|----------|-------------|---------|
| `plural` | Pluralize a word | `{{plural .NameLower}}` → `users` |

## Example: Entity Template

```go
package {{.BasePackage}}.entity;

import jakarta.persistence.*;
{{if .HasLombok}}
import lombok.*;
{{end}}

@Entity
@Table(name = "{{plural .NameLower}}")
{{if .HasLombok}}
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
{{end}}
public class {{.Name}} {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
{{if not .HasLombok}}
    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }
{{end}}
}
```

## Example: Controller Template

```go
package {{.BasePackage}}.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
{{if .HasValidation}}
import jakarta.validation.Valid;
{{end}}
import {{.BasePackage}}.service.{{.Name}}Service;
import {{.BasePackage}}.dto.{{.Name}}Request;
import {{.BasePackage}}.dto.{{.Name}}Response;

import java.util.List;

@RestController
@RequestMapping("/api/{{plural .NameLower}}")
public class {{.Name}}Controller {

    private final {{.Name}}Service {{.NameCamel}}Service;

    public {{.Name}}Controller({{.Name}}Service {{.NameCamel}}Service) {
        this.{{.NameCamel}}Service = {{.NameCamel}}Service;
    }

    @GetMapping
    public ResponseEntity<List<{{.Name}}Response>> getAll() {
        return ResponseEntity.ok({{.NameCamel}}Service.findAll());
    }

    @PostMapping
    public ResponseEntity<{{.Name}}Response> create(
            {{if .HasValidation}}@Valid {{end}}@RequestBody {{.Name}}Request request) {
        return ResponseEntity.ok({{.NameCamel}}Service.create(request));
    }
}
```

## Custom Templates (Future)

:::info Planned Feature
Custom template support is planned for a future release. You'll be able to override built-in templates or add your own.
:::

Planned features:
- Local template directory (`~/.haft/templates/`)
- Project-level templates (`.haft/templates/`)
- Template inheritance
- Custom template functions

## See Also

- [Project Structure](/docs/guides/project-structure) — Generated file locations
- [haft generate](/docs/commands/generate) — Generation commands
