---
sidebar_position: 3
title: haft add
description: Add dependencies to your project
---

# haft add

Add dependencies to an existing Spring Boot project.

## Usage

```bash
haft add                              # Interactive search picker
haft add --browse                     # Browse by category
haft add <dependency> [dependencies...]
haft add <groupId:artifactId>
haft add <groupId:artifactId:version>
```

## Description

The `add` command modifies your `pom.xml` to add new dependencies. It supports:

- **Interactive mode** — Search and select from 50+ shortcuts
- **Browse mode** — Navigate dependencies by category
- **Shortcuts** — Common dependencies like `lombok`, `jpa`, `web`, `jwt`
- **Maven coordinates** — Any dependency as `groupId:artifactId`
- **Maven Central verification** — Auto-verify and fetch latest versions

## Interactive Modes

### Search Picker (Default)

```bash
haft add
```

Opens an interactive TUI where you can:
- Type to search/filter dependencies
- Select multiple with `Space`
- Navigate with `↑`/`↓`, `PgUp`/`PgDown`
- Select all visible with `a`, none with `n`
- Confirm with `Enter`, cancel with `Esc`

### Category Browser

```bash
haft add --browse
haft add -b
```

Opens a category-based browser similar to the init wizard:
- Jump to categories with `0-9` keys
- Cycle categories with `Tab`/`Shift+Tab`
- Search within category with `/`

## Examples

### Add Using Shortcuts

```bash
# Add Lombok
haft add lombok

# Add multiple dependencies
haft add jpa validation lombok

# Add JWT (adds all 3 JJWT artifacts)
haft add jwt

# Add database driver
haft add postgresql
```

### Add Using Maven Coordinates

```bash
# Without version (fetches latest from Maven Central)
haft add org.mapstruct:mapstruct

# With specific version
haft add io.jsonwebtoken:jjwt-api:0.12.5
```

Dependencies are automatically verified against Maven Central. If a dependency doesn't exist, you'll get an error:

```
ERROR ✗ dependency 'com.fake:nonexistent' not found on Maven Central
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

| Flag | Short | Description |
|------|-------|-------------|
| `--browse` | `-b` | Browse dependencies by category |
| `--list` | | List available dependency shortcuts |
| `--scope` | | Set dependency scope (compile, runtime, test, provided) |
| `--version` | | Override default version |

## Available Shortcuts (50+)

### Web

| Shortcut | Description |
|----------|-------------|
| `web` | Spring Boot Web (Spring MVC) |
| `webflux` | Spring WebFlux (reactive) |
| `graphql` | Spring GraphQL |
| `websocket` | WebSocket support |
| `hateoas` | Hypermedia-driven REST APIs |
| `data-rest` | Expose repositories as REST endpoints |
| `feign` | Declarative REST client (OpenFeign) |
| `resilience4j` | Fault tolerance (circuit breaker) |

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
| `cassandra` | Spring Data Cassandra |
| `neo4j` | Spring Data Neo4j |

### Security

| Shortcut | Description |
|----------|-------------|
| `security` | Spring Security |
| `oauth2-client` | OAuth2 client |
| `oauth2-resource-server` | OAuth2 resource server |
| `jwt` | JJWT library (api + impl + jackson) |

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
| `commons-lang` | Apache Commons Lang |
| `guava` | Google Guava |
| `modelmapper` | ModelMapper |

### I/O

| Shortcut | Description |
|----------|-------------|
| `validation` | Bean Validation |
| `mail` | Java Mail |
| `cache` | Spring Cache |
| `batch` | Spring Batch |
| `quartz` | Quartz Scheduler |
| `commons-io` | Apache Commons IO |
| `jackson-datatype` | Jackson Java 8 datatypes |

### Ops

| Shortcut | Description |
|----------|-------------|
| `actuator` | Spring Boot Actuator |
| `micrometer` | Prometheus metrics |

### Template Engines

| Shortcut | Description |
|----------|-------------|
| `thymeleaf` | Thymeleaf templates |
| `freemarker` | FreeMarker templates |
| `mustache` | Mustache templates |

### Testing

| Shortcut | Description |
|----------|-------------|
| `test` | Spring Boot Test |
| `testcontainers` | Testcontainers |
| `security-test` | Spring Security Test |
| `mockito` | Mockito mocking |
| `restdocs` | Spring REST Docs |

## What Gets Added

### Example: `haft add lombok`

```xml
<dependency>
    <groupId>org.projectlombok</groupId>
    <artifactId>lombok</artifactId>
    <scope>provided</scope>
</dependency>
```

### Example: `haft add jwt`

Adds all three JJWT artifacts:

```xml
<dependency>
    <groupId>io.jsonwebtoken</groupId>
    <artifactId>jjwt-api</artifactId>
    <version>0.12.5</version>
</dependency>
<dependency>
    <groupId>io.jsonwebtoken</groupId>
    <artifactId>jjwt-impl</artifactId>
    <version>0.12.5</version>
    <scope>runtime</scope>
</dependency>
<dependency>
    <groupId>io.jsonwebtoken</groupId>
    <artifactId>jjwt-jackson</artifactId>
    <version>0.12.5</version>
    <scope>runtime</scope>
</dependency>
```

### Example: `haft add org.mapstruct:mapstruct`

Auto-fetches latest version from Maven Central:

```xml
<dependency>
    <groupId>org.mapstruct</groupId>
    <artifactId>mapstruct</artifactId>
    <version>1.5.5.Final</version>
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

- [haft remove](/docs/commands/remove) — Remove dependencies
- [haft init](/docs/commands/init) — Add dependencies at project creation
- [Dependencies Guide](/docs/guides/dependencies) — Full dependency reference
