---
sidebar_position: 14
title: haft dockerize
description: Generate Docker configuration files
---

# haft dockerize

Generate optimized Docker configuration files for your Spring Boot project.

## Usage

```bash
haft dockerize                    # Auto-detect everything
haft dockerize --db postgresql    # Specify database
haft dockerize --no-compose       # Skip docker-compose.yml
haft dockerize --port 9000        # Override application port
haft dockerize --java 21          # Override Java version
```

## Description

The `dockerize` command generates production-ready Docker configuration files for your Spring Boot project:

- **Dockerfile** — Multi-stage build optimized for minimal image size
- **docker-compose.yml** — Service orchestration with database
- **.dockerignore** — Exclude unnecessary files from build context

## Auto-Detection

The command automatically detects:

| Detected | Source |
|----------|--------|
| Build tool | `pom.xml` → Maven, `build.gradle` → Gradle |
| Java version | `pom.xml` properties, `build.gradle` config |
| Application port | `application.properties` or `application.yml` |
| Application name | `artifactId` or project directory name |
| Database | Dependency analysis (PostgreSQL, MySQL, etc.) |
| Wrapper scripts | `mvnw` or `gradlew` presence |

## Examples

### Basic Usage

```bash
# Generate all Docker files with auto-detection
haft dockerize

# Output:
# ✓ Created Dockerfile
# ✓ Created .dockerignore
# ✓ Created docker-compose.yml
# ✓ Generated 3 Docker file(s)
```

### Specify Database

```bash
# Explicitly set PostgreSQL
haft dockerize --db postgresql

# Explicitly set MySQL
haft dockerize --db mysql

# Explicitly set no database
haft dockerize --db none
```

### Dockerfile Only

```bash
# Skip docker-compose.yml generation
haft dockerize --no-compose
```

### Override Defaults

```bash
# Override Java version
haft dockerize --java 17

# Override application port
haft dockerize --port 9000

# Combine options
haft dockerize --java 21 --port 9000 --db postgresql
```

### Non-Interactive Mode

```bash
# Skip all prompts
haft dockerize --no-interactive

# JSON output for scripting
haft dockerize --json
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--port` | `-p` | Application port (default: auto-detect or 8080) |
| `--java` | `-j` | Java version (default: auto-detect from build file) |
| `--db` | | Database type (postgresql, mysql, mariadb, mongodb, redis, none) |
| `--no-compose` | | Skip docker-compose.yml generation |
| `--no-interactive` | | Skip interactive prompts |
| `--force` | `-f` | Overwrite existing files |
| `--json` | | Output result as JSON |

## Supported Databases

| Database | Image | Port |
|----------|-------|------|
| PostgreSQL | `postgres:16-alpine` | 5432 |
| MySQL | `mysql:8` | 3306 |
| MariaDB | `mariadb:11` | 3306 |
| MongoDB | `mongo:7` | 27017 |
| Redis | `redis:7-alpine` | 6379 |
| Cassandra | `cassandra:4` | 9042 |

## Generated Files

### Dockerfile (Maven)

Multi-stage build with layer extraction for optimal caching:

```dockerfile
# Build stage
FROM eclipse-temurin:21-jdk-alpine AS builder
WORKDIR /app

# Copy Maven wrapper
COPY mvnw .
COPY .mvn .mvn
RUN chmod +x mvnw

# Copy pom.xml and download dependencies (cached layer)
COPY pom.xml .
RUN ./mvnw dependency:go-offline -B

# Copy source and build
COPY src src
RUN ./mvnw package -DskipTests -B

# Extract layers for better caching
RUN java -Djarmode=layertools -jar target/*.jar extract --destination target/extracted

# Runtime stage
FROM eclipse-temurin:21-jre-alpine AS runtime
WORKDIR /app

# Security: run as non-root user
RUN addgroup -g 1001 spring && adduser -u 1001 -G spring -D spring
USER spring:spring

# Copy layers (most stable to least stable)
COPY --from=builder --chown=spring:spring /app/target/extracted/dependencies/ ./
COPY --from=builder --chown=spring:spring /app/target/extracted/spring-boot-loader/ ./
COPY --from=builder --chown=spring:spring /app/target/extracted/snapshot-dependencies/ ./
COPY --from=builder --chown=spring:spring /app/target/extracted/application/ ./

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/actuator/health || exit 1

EXPOSE 8080

ENTRYPOINT ["java", "org.springframework.boot.loader.launch.JarLauncher"]
```

### Dockerfile (Gradle)

Similar multi-stage build for Gradle projects:

```dockerfile
# Build stage
FROM eclipse-temurin:21-jdk-alpine AS builder
WORKDIR /app

# Copy Gradle wrapper
COPY gradlew .
COPY gradle gradle
RUN chmod +x gradlew

# Copy build files and download dependencies (cached layer)
COPY build.gradle* settings.gradle* ./
RUN ./gradlew dependencies --no-daemon || true

# Copy source and build
COPY src src
RUN ./gradlew bootJar --no-daemon -x test

# Extract layers
RUN java -Djarmode=layertools -jar build/libs/*.jar extract --destination build/extracted

# ... runtime stage same as Maven
```

### docker-compose.yml

Generated with database service and health checks:

```yaml
services:
  app:
    build: .
    container_name: myapp
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - SPRING_DATASOURCE_URL=jdbc:postgresql://postgres:5432/myapp
      - SPRING_DATASOURCE_USERNAME=postgres
      - SPRING_DATASOURCE_PASSWORD=postgres
    networks:
      - myapp-network
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    container_name: myapp-postgres
    environment:
      - POSTGRES_DB=myapp
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - myapp-network
    restart: unless-stopped

networks:
  myapp-network:
    driver: bridge

volumes:
  postgres_data:
```

### .dockerignore

Optimized exclusions for smaller build context:

```
# Build artifacts
target/
build/
*.jar
*.class

# IDE files
.idea/
*.iml
.vscode/

# Version control
.git/
.gitignore

# Docker files
Dockerfile*
docker-compose*.yml

# CI/CD
.github/
Jenkinsfile

# Documentation
*.md
docs/

# Environment files
.env
*.env

# Test files
src/test/

# Wrappers (kept in build stage)
.gradle/
.mvn/wrapper/maven-wrapper.jar

# Haft config
.haft/
.haft.json
```

## Build & Run

After generating Docker files:

```bash
# Build the Docker image
docker build -t myapp .

# Run with docker-compose
docker compose up -d

# View logs
docker compose logs -f app

# Stop services
docker compose down

# Stop and remove volumes
docker compose down -v
```

## JPA Without Database Driver

If your project has JPA dependency but no database driver detected, the command will interactively ask which database to include in docker-compose.yml:

```
JPA detected but no database driver found. Select a database for docker-compose:
  ➜ PostgreSQL (Recommended for most applications)
    MySQL (Popular relational database)
    MariaDB (MySQL-compatible database)
    MongoDB (Document-oriented NoSQL database)
    Redis (In-memory data store - caching)
    None (No database service)
```

Use `--no-interactive` to skip this prompt.

## Existing Files

If Docker files already exist, they are skipped (not overwritten):

```
⚠ File exists, skipping file=Dockerfile
✓ Created file=.dockerignore
✓ Created file=docker-compose.yml
✓ Generated 2 Docker file(s)
ℹ Skipped 1 existing file(s)
```

## JSON Output

For scripting and automation:

```bash
haft dockerize --json
```

```json
{
  "generated": ["Dockerfile", ".dockerignore", "docker-compose.yml"],
  "skipped": [],
  "config": {
    "appName": "myapp",
    "javaVersion": "21",
    "port": 8080,
    "buildTool": "maven",
    "databaseType": "postgres"
  }
}
```

## Best Practices

### Layer Optimization

The generated Dockerfile uses Spring Boot's layertools to extract layers:
1. **dependencies** — Third-party libraries (rarely change)
2. **spring-boot-loader** — Spring Boot launcher
3. **snapshot-dependencies** — Snapshot versions
4. **application** — Your compiled code (changes frequently)

This ordering ensures Docker caches the stable layers, making rebuilds faster.

### Security

The generated Dockerfile:
- Runs as non-root user (`spring:spring`)
- Uses Alpine-based images for smaller attack surface
- Separates build and runtime stages (no JDK in production)

### Health Checks

Both Dockerfile and docker-compose.yml include health checks:
- Docker: Uses wget to check `/actuator/health`
- Compose: Database-specific health checks for dependency ordering

Ensure Spring Boot Actuator is enabled for health endpoints:

```bash
haft add actuator
```

## Editor Integration

Use this command from your editor:

- **Neovim**: Not yet available in haft.nvim (CLI only)
- **VS Code**: Coming soon ([preview →](/docs/integrations/vscode))
- **IntelliJ IDEA**: Coming soon ([preview →](/docs/integrations/intellij))

## See Also

- [haft dev serve](/docs/commands/dev) — Development mode with hot-reload
- [haft info](/docs/commands/info) — Show project information
