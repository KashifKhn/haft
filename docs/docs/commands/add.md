---
sidebar_position: 3
title: haft add
description: Add dependencies to your project
---

# haft add

Add dependencies to an existing Spring Boot project.

## Usage

```bash
haft add <dependency> [dependencies...]
haft add <groupId:artifactId>
haft add <groupId:artifactId:version>
```

## Description

The `add` command modifies your `pom.xml` to add new dependencies. It supports:

- **Shortcuts** - Common dependencies like `lombok`, `jpa`, `web`
- **Maven coordinates** - Any dependency as `groupId:artifactId`
- **With version** - Specify version as `groupId:artifactId:version`

## Examples

### Add Using Shortcuts

```bash
# Add Lombok
haft add lombok

# Add multiple dependencies
haft add jpa validation lombok

# Add database driver
haft add postgresql
```

### Add Using Maven Coordinates

```bash
# Without version (uses managed version)
haft add org.mapstruct:mapstruct

# With specific version
haft add io.jsonwebtoken:jjwt-api:0.12.3
```

### Override Scope

```bash
# Add as test dependency
haft add h2 --scope test

# Add as provided
haft add org.example:my-processor --scope provided
```

### List Available Shortcuts

```bash
haft add --list
```

## Flags

| Flag | Description |
|------|-------------|
| `--list` | List available dependency shortcuts |
| `--scope` | Set dependency scope (compile, runtime, test, provided) |
| `--version` | Override default version |

## Available Shortcuts

### Web

| Shortcut | Description |
|----------|-------------|
| `web` | Spring Boot Web (Spring MVC) |
| `webflux` | Spring WebFlux (reactive) |
| `graphql` | Spring GraphQL |
| `websocket` | WebSocket support |

### SQL

| Shortcut | Description |
|----------|-------------|
| `jpa` | Spring Data JPA |
| `jdbc` | Spring JDBC |
| `postgresql` | PostgreSQL driver |
| `mysql` | MySQL driver |
| `mariadb` | MariaDB driver |
| `h2` | H2 in-memory database |
| `flyway` | Flyway migrations |
| `liquibase` | Liquibase migrations |

### NoSQL

| Shortcut | Description |
|----------|-------------|
| `mongodb` | Spring Data MongoDB |
| `redis` | Spring Data Redis |
| `elasticsearch` | Spring Data Elasticsearch |

### Security

| Shortcut | Description |
|----------|-------------|
| `security` | Spring Security |
| `oauth2-client` | OAuth2 client |
| `oauth2-resource-server` | OAuth2 resource server |

### Messaging

| Shortcut | Description |
|----------|-------------|
| `amqp` | RabbitMQ (Spring AMQP) |
| `kafka` | Apache Kafka |

### Developer Tools

| Shortcut | Description |
|----------|-------------|
| `lombok` | Lombok annotations |
| `devtools` | Spring Boot DevTools |
| `mapstruct` | MapStruct bean mapping |
| `openapi` | SpringDoc OpenAPI (Swagger UI) |

### Ops & I/O

| Shortcut | Description |
|----------|-------------|
| `actuator` | Spring Boot Actuator |
| `validation` | Bean Validation |
| `mail` | Java Mail |
| `cache` | Spring Cache |
| `batch` | Spring Batch |
| `quartz` | Quartz Scheduler |

### Testing

| Shortcut | Description |
|----------|-------------|
| `test` | Spring Boot Test |
| `testcontainers` | Testcontainers |

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

Adds both the main library and annotation processor:

```xml
<dependency>
    <groupId>org.mapstruct</groupId>
    <artifactId>mapstruct</artifactId>
    <version>1.5.5.Final</version>
</dependency>
<dependency>
    <groupId>org.mapstruct</groupId>
    <artifactId>mapstruct-processor</artifactId>
    <version>1.5.5.Final</version>
    <scope>provided</scope>
</dependency>
```

### Example: `haft add postgresql`

```xml
<dependency>
    <groupId>org.postgresql</groupId>
    <artifactId>postgresql</artifactId>
    <scope>runtime</scope>
</dependency>
```

## Duplicate Detection

Haft automatically detects existing dependencies and skips them:

```
$ haft add lombok
WARN ⚠ Skipped (already exists) dependency=org.projectlombok:lombok
INFO ℹ No new dependencies added (all already exist)
```

## See Also

- [haft init](/docs/commands/init) — Add dependencies at project creation
- [Dependencies Guide](/docs/guides/dependencies) — Full dependency reference
