<p align="center">
  <img src="assets/haft-logo.png" alt="Haft Logo" width="200"/>
</p>

<h1 align="center">Haft</h1>

<p align="center">
  <strong>The Spring Boot CLI that Spring forgot to build</strong>
</p>

<p align="center">
  <a href="https://github.com/KashifKhn/haft/releases"><img src="https://img.shields.io/github/v/release/KashifKhn/haft?style=flat-square" alt="Release"></a>
  <a href="https://github.com/KashifKhn/haft/blob/main/LICENSE"><img src="https://img.shields.io/github/license/KashifKhn/haft?style=flat-square" alt="License"></a>
  <a href="https://github.com/KashifKhn/haft/actions"><img src="https://img.shields.io/github/actions/workflow/status/KashifKhn/haft/ci.yml?style=flat-square" alt="Build Status"></a>
  <a href="https://goreportcard.com/report/github.com/KashifKhn/haft"><img src="https://goreportcard.com/badge/github.com/KashifKhn/haft?style=flat-square" alt="Go Report Card"></a>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#features">Features</a> •
  <a href="#commands">Commands</a> •
  <a href="#contributing">Contributing</a>
</p>

---

<p align="center">
  <img src="assets/haft-demo.gif" alt="Haft Demo" width="800"/>
</p>

## Why Haft?

**Spring Initializr** bootstraps your project. **Haft** is your lifecycle companion.

While Spring Initializr gets you started, you're on your own after that. Every new entity means manually creating 8+ files: Entity, Repository, Service, ServiceImpl, Controller, Request DTO, Response DTO, Mapper, and Exception handler. Every. Single. Time.

**Haft changes that.**

```bash
# Instead of creating 8 files manually...
haft generate resource User

# Done. All files generated with proper layered architecture.
```

## Features

- **Interactive TUI Wizard** - Beautiful terminal UI for project setup
- **Spring Initializr Integration** - All official dependencies with descriptions
- **Smart Dependency Detection** - Detects Lombok, MapStruct, Validation from pom.xml
- **Full CRUD Generation** - Entity, Repository, Service, Controller, DTOs in one command
- **Maven Support** - Parse, read, and modify pom.xml programmatically
- **Git Integration** - Optional repository initialization on project creation

## Installation

### Using Go

```bash
go install github.com/KashifKhn/haft/cmd/haft@latest
```

### From Source

```bash
git clone https://github.com/KashifKhn/haft.git
cd haft
make build
```

### Binary Releases

Download pre-built binaries from the [Releases](https://github.com/KashifKhn/haft/releases) page.

<details>
<summary>Linux</summary>

```bash
curl -L https://github.com/KashifKhn/haft/releases/latest/download/haft-linux-amd64.tar.gz | tar xz
sudo mv haft /usr/local/bin/
```

</details>

<details>
<summary>macOS</summary>

```bash
curl -L https://github.com/KashifKhn/haft/releases/latest/download/haft-darwin-amd64.tar.gz | tar xz
sudo mv haft /usr/local/bin/
```

</details>

<details>
<summary>Windows</summary>

Download `haft-windows-amd64.zip` from the releases page and add to your PATH.

</details>

## Quick Start

### Initialize a New Project

```bash
# Interactive mode (recommended)
haft init

# With project name
haft init my-app

# Non-interactive with options
haft init my-app -g com.example -j 21 -s 3.4.1 --deps web,data-jpa,lombok
```

The interactive wizard guides you through:

| Step | Description |
|------|-------------|
| Project Name | Your application name |
| Group ID | Maven group (e.g., `com.example`) |
| Artifact ID | Maven artifact name |
| Description | Project description |
| Package Name | Base package (auto-generated) |
| Java Version | 17, 21 (LTS), or 25 |
| Spring Boot | Latest stable versions |
| Build Tool | Maven or Gradle |
| Packaging | JAR or WAR |
| Config Format | Properties or YAML |
| Dependencies | Search & select from all Spring starters |
| Git Init | Initialize git repository |

### Wizard Navigation

| Key | Action |
|-----|--------|
| `↑` `↓` | Navigate options |
| `Enter` | Select / Confirm |
| `Esc` | Go back to previous step |
| `Space` | Toggle selection (multi-select) |
| `/` | Search dependencies |
| `Tab` | Next category |
| `0-9` | Jump to category |

## Commands

### `haft init`

Initialize a new Spring Boot project.

```bash
haft init [name] [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--group` | `-g` | Group ID (e.g., com.example) |
| `--artifact` | `-a` | Artifact ID |
| `--java` | `-j` | Java version (17, 21, 25) |
| `--spring` | `-s` | Spring Boot version |
| `--build` | `-b` | Build tool (maven, gradle) |
| `--deps` | | Dependencies (comma-separated) |
| `--package` | | Base package name |
| `--packaging` | | Packaging type (jar, war) |
| `--config` | | Config format (properties, yaml) |
| `--dir` | `-d` | Output directory |
| `--no-interactive` | | Skip interactive wizard |

**Examples:**

```bash
# Full interactive wizard
haft init

# Quick start with defaults
haft init my-api -g com.company -j 21

# CI/CD friendly (non-interactive)
haft init user-service \
  --group com.example \
  --java 21 \
  --spring 3.4.1 \
  --build maven \
  --deps web,data-jpa,lombok,validation \
  --no-interactive
```

### `haft generate resource` (Coming Soon)

Generate a complete CRUD resource with all layers.

```bash
haft generate resource <name> [flags]
```

Generates:
- `Entity.java` - JPA entity with Lombok (if available)
- `Repository.java` - Spring Data JPA repository
- `Service.java` - Service interface
- `ServiceImpl.java` - Service implementation
- `Controller.java` - REST controller with CRUD endpoints
- `Request.java` - Request DTO
- `Response.java` - Response DTO
- `Mapper.java` - MapStruct mapper (if available)

### `haft add` (Coming Soon)

Add dependencies to your project.

```bash
haft add <dependency>
```

## Project Structure

Generated projects follow Spring Boot best practices:

```
my-app/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/example/myapp/
│   │   │       └── MyAppApplication.java
│   │   └── resources/
│   │       └── application.yml
│   └── test/
│       └── java/
│           └── com/example/myapp/
│               └── MyAppApplicationTests.java
├── .gitignore
├── mvnw
├── mvnw.cmd
└── pom.xml
```

## Dependency Categories

The dependency picker organizes all Spring starters by category:

| # | Category | Examples |
|---|----------|----------|
| 1 | Developer Tools | DevTools, Lombok, Configuration Processor |
| 2 | Web | Spring Web, WebFlux, GraphQL, REST Docs |
| 3 | SQL | JPA, JDBC, MySQL, PostgreSQL, H2 |
| 4 | NoSQL | MongoDB, Redis, Elasticsearch, Cassandra |
| 5 | Security | Spring Security, OAuth2, LDAP |
| 6 | Messaging | Kafka, RabbitMQ, ActiveMQ |
| 7 | Cloud | Config, Discovery, Gateway |
| 8 | Observability | Actuator, Micrometer, Zipkin |
| 9 | Testing | Testcontainers, Contract Testing |

## Requirements

- Go 1.21+ (for building from source)
- Java 17+ (for generated projects)
- Maven 3.6+ or Gradle 7+ (for generated projects)

## Roadmap

- [x] Interactive project initialization
- [x] Spring Initializr dependency integration
- [x] Maven pom.xml parser
- [x] Back navigation in wizard
- [x] Dependency search and filtering
- [ ] `haft generate resource` - Full CRUD generator
- [ ] `haft generate controller|service|entity` - Individual generators
- [ ] `haft add` - Dependency manager
- [ ] Gradle support improvements
- [ ] Shell completions (bash, zsh, fish)
- [ ] Custom templates support

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

```bash
# Clone the repository
git clone https://github.com/KashifKhn/haft.git
cd haft

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
make build

# Run
./bin/haft --help
```

## License

Haft is open-source software licensed under the [MIT License](LICENSE).

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Spring Initializr](https://start.spring.io) - Dependency metadata

---

<p align="center">
  <sub>Built with ❤️ for the Spring Boot community</sub>
</p>
