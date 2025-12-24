<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="assets/logo-dark.png">
    <source media="(prefers-color-scheme: light)" srcset="assets/logo-light.png">
    <img src="assets/logo-light.png" alt="Haft" width="320"/>
  </picture>
</p>

<p align="center">
  <strong>The Spring Boot CLI that Spring forgot to build</strong>
</p>

<p align="center">
  <a href="https://github.com/KashifKhn/haft/releases"><img src="https://img.shields.io/github/v/release/KashifKhn/haft?style=for-the-badge&logo=github&color=blue" alt="Release"></a>
  <a href="https://github.com/KashifKhn/haft/blob/main/LICENSE"><img src="https://img.shields.io/github/license/KashifKhn/haft?style=for-the-badge" alt="License"></a>
  <a href="https://goreportcard.com/report/github.com/KashifKhn/haft"><img src="https://goreportcard.com/badge/github.com/KashifKhn/haft?style=for-the-badge" alt="Go Report"></a>
</p>

<p align="center">
  <a href="https://kashifkhn.github.io/haft">Documentation</a> ·
  <a href="https://github.com/KashifKhn/haft/releases">Releases</a> ·
  <a href="https://github.com/KashifKhn/haft/issues">Report Bug</a> ·
  <a href="https://github.com/KashifKhn/haft/discussions">Discussions</a>
</p>

---

## The Problem

You start a Spring Boot project with Spring Initializr. Great. Now what?

Every new feature means the same tedious ritual:
- Create `UserEntity.java`
- Create `UserRepository.java`  
- Create `UserService.java`
- Create `UserServiceImpl.java`
- Create `UserController.java`
- Create `UserRequest.java`
- Create `UserResponse.java`
- Create `UserMapper.java`

**8 files. Every. Single. Time.**

Copy-paste from existing code. Fix the class names. Fix the imports. Miss something. Debug. Repeat.

## The Solution

```bash
haft generate resource User
```

Done. All 8 files. Properly structured. Following your project's conventions.

## Why Haft?

| | Spring Initializr | Haft |
|---|---|---|
| Project Bootstrap | ✅ | ✅ |
| Works Offline | ❌ | ✅ |
| Resource Generation | ❌ | ✅ |
| Dependency Management | ❌ | ✅ |
| Interactive TUI | ❌ | ✅ |
| Lifecycle Companion | ❌ | ✅ |

**Haft works completely offline.** No web browser. No internet dependency. Just you and your terminal.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/KashifKhn/haft/main/install.sh | bash
```

<details>
<summary>Other installation methods</summary>

**Go**
```bash
go install github.com/KashifKhn/haft/cmd/haft@latest
```

**From Source**
```bash
git clone https://github.com/KashifKhn/haft.git && cd haft && make build
```

**Manual Download**

Download binaries from [Releases](https://github.com/KashifKhn/haft/releases) for Linux, macOS, or Windows.

</details>

## Quick Start

### Create a New Project

```bash
haft init
```

An interactive wizard guides you through project setup:

```
? Project name: my-api
? Group ID: com.example  
? Java version: 21
? Spring Boot: 3.4.1
? Dependencies: web, data-jpa, lombok, validation
```

### Non-Interactive Mode

Perfect for CI/CD and scripting:

```bash
haft init my-service \
  --group com.example \
  --java 21 \
  --deps web,data-jpa,lombok \
  --no-interactive
```

### Generate Resources (Coming Soon)

```bash
haft generate resource Product
```

Creates a complete CRUD stack with Entity, Repository, Service, Controller, and DTOs.

## Features

- **Interactive TUI** — Beautiful terminal interface with keyboard navigation
- **Offline First** — No internet required, all metadata bundled
- **Spring Initializr Parity** — All official starters and dependencies
- **Smart Defaults** — Sensible defaults that match industry standards
- **Back Navigation** — Made a mistake? Press `Esc` to go back
- **Dependency Search** — Find any dependency with `/`
- **Git Integration** — Initialize repository on project creation

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` `↓` | Navigate |
| `Enter` | Select |
| `Esc` | Go back |
| `Space` | Toggle |
| `/` | Search |
| `Tab` | Next category |
| `0-9` | Jump to category |

## Roadmap

- [x] Project initialization wizard
- [x] All Spring Initializr dependencies  
- [x] Maven support
- [x] Offline operation
- [ ] `haft generate resource` — Full CRUD generation
- [ ] `haft generate controller|service|entity` — Individual generators
- [ ] `haft add` — Dependency management
- [ ] Gradle improvements
- [ ] Neovim integration
- [ ] VS Code extension
- [ ] IntelliJ plugin
- [ ] Custom templates

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License. See [LICENSE](LICENSE) for details.

---

<p align="center">
  <sub>Built for developers who value their time</sub>
</p>
