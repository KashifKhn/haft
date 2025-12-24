---
sidebar_position: 2
title: haft generate
description: Generate CRUD resources and individual components
---

# haft generate

Generate boilerplate code for Spring Boot applications.

:::caution Coming Soon
The `generate` command is currently under development and will be available in a future release.
:::

## Usage

```bash
haft generate <type> <name> [flags]
```

## Subcommands

### haft generate resource

Generate a complete CRUD resource with all layers.

```bash
haft generate resource User
```

This generates:

| File | Description |
|------|-------------|
| `entity/User.java` | JPA entity |
| `repository/UserRepository.java` | Spring Data repository |
| `service/UserService.java` | Service interface |
| `service/impl/UserServiceImpl.java` | Service implementation |
| `controller/UserController.java` | REST controller |
| `dto/UserRequest.java` | Request DTO |
| `dto/UserResponse.java` | Response DTO |
| `mapper/UserMapper.java` | MapStruct mapper (if available) |
| `exception/UserNotFoundException.java` | Resource exception |

### haft generate controller

Generate only a REST controller.

```bash
haft generate controller Product
```

### haft generate service

Generate a service interface and implementation.

```bash
haft generate service Order
```

### haft generate entity

Generate a JPA entity.

```bash
haft generate entity Customer
```

### haft generate repository

Generate a Spring Data JPA repository.

```bash
haft generate repository Invoice
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--path` | `-p` | Custom output path |
| `--package` | | Override base package |
| `--no-lombok` | | Don't use Lombok annotations |
| `--no-mapper` | | Don't generate MapStruct mapper |

## Smart Detection

Haft reads your `pom.xml` to automatically detect and use:

- **Lombok** — Uses `@Data`, `@Builder`, `@NoArgsConstructor`, etc.
- **MapStruct** — Generates mapper interfaces with `@Mapper`
- **Validation** — Adds `@Valid` to controller parameters
- **Spring Data JPA** — Configures repository correctly

## Example Output

### With Lombok

```java
@Entity
@Table(name = "users")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
}
```

### Without Lombok

```java
@Entity
@Table(name = "users")
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }
}
```

## See Also

- [Project Structure](/docs/guides/project-structure) — Where files are generated
- [Templates](/docs/reference/templates) — Customize generated code
