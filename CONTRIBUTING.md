# Contributing to Haft

Thank you for your interest in contributing to Haft! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Standards](#code-standards)
- [Making Changes](#making-changes)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Project Structure](#project-structure)

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment. Be kind, constructive, and professional in all interactions.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Fork and Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/haft.git
cd haft

# Add upstream remote
git remote add upstream https://github.com/KashifKhn/haft.git

# Install dependencies
go mod download
```

## Development Setup

### Building

```bash
# Using Make
make build

# Or directly with Go
go build -o bin/haft ./cmd/haft
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/generator/... -v

# Run a single test
go test -run TestFunctionName ./...
```

### Running the CLI

```bash
# After building
./bin/haft --help
./bin/haft init
```

## Code Standards

Haft follows strict coding standards to maintain quality and consistency.

### No Comments Policy

**Do not write code comments or doc comments** unless absolutely necessary. Code should be self-documenting through:

- Clear, descriptive function names
- Meaningful variable names
- Small, focused functions
- Logical code organization

```go
// Bad
// GetUserByID retrieves a user by their ID from the database
func GetUserByID(id int) (*User, error) { ... }

// Good - the function name is self-explanatory
func GetUserByID(id int) (*User, error) { ... }
```

### DRY Principle

Never repeat code. Extract common logic into reusable functions.

```go
// Bad - repeated validation logic
func CreateUser(u User) error {
    if u.Name == "" {
        return errors.New("name is required")
    }
    // ...
}

func UpdateUser(u User) error {
    if u.Name == "" {
        return errors.New("name is required")
    }
    // ...
}

// Good - extracted validation
func (u User) Validate() error {
    if u.Name == "" {
        return errors.New("name is required")
    }
    return nil
}

func CreateUser(u User) error {
    if err := u.Validate(); err != nil {
        return err
    }
    // ...
}
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Exported functions/types | PascalCase | `ParseConfig`, `UserService` |
| Private functions/types | camelCase | `parseFile`, `userCache` |
| Files | snake_case | `user_service.go`, `config_test.go` |
| Constants | PascalCase | `DefaultTimeout`, `MaxRetries` |

### Function Guidelines

- **Maximum 50 lines** per function
- **Single responsibility** - each function does one thing
- Split large functions into smaller, focused helpers

### File Guidelines

- **Maximum 300 lines** per file
- **One concept per file** - don't mix unrelated functionality
- Group related types and functions together

### Error Handling

- Always return errors with context
- Never use `panic()` for recoverable errors
- Provide actionable error messages

```go
// Bad
return errors.New("failed")

// Good
return fmt.Errorf("failed to parse pom.xml at %s: %w", path, err)
```

## Making Changes

### Branch Naming

Use descriptive branch names:

```
feature/add-gradle-support
fix/wizard-back-navigation
docs/update-readme
refactor/simplify-parser
```

### Commit Messages

Follow conventional commit format:

```
type: short description

[optional body]
[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**

```
feat: add dependency search in wizard

fix: handle empty pom.xml gracefully

docs: update installation instructions

refactor: extract validation logic to separate package
```

### Development Workflow

1. **Create a branch** from `dev`:
   ```bash
   git checkout dev
   git pull upstream dev
   git checkout -b feature/your-feature
   ```

2. **Make your changes** following code standards

3. **Test your changes**:
   ```bash
   go test ./...
   go build ./...
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: your feature description"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature
   ```

6. **Open a Pull Request** against `dev` branch

## Pull Request Process

### Before Submitting

- [ ] Code follows the [Code Standards](#code-standards)
- [ ] All tests pass (`go test ./...`)
- [ ] Build succeeds (`go build ./...`)
- [ ] No new linter warnings
- [ ] Commit messages follow conventions

### PR Description

Include in your PR description:

1. **What** - Brief description of changes
2. **Why** - Motivation for the change
3. **How** - Implementation approach (if complex)
4. **Testing** - How you tested the changes

### Review Process

1. Maintainers will review your PR
2. Address any requested changes
3. Once approved, your PR will be merged

## Testing

### Test Requirements

- **Target 80%+ code coverage** for new code
- Write unit tests for all exported functions
- Use table-driven tests where appropriate
- Mock external dependencies using `afero` for filesystem

### Test Structure

```go
func TestParser_ParseBytes(t *testing.T) {
    tests := []struct {
        name    string
        input   []byte
        want    *Project
        wantErr bool
    }{
        {
            name:  "valid pom.xml",
            input: []byte(`<project>...</project>`),
            want:  &Project{GroupId: "com.example"},
        },
        {
            name:    "invalid xml",
            input:   []byte(`not xml`),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := NewParser()
            got, err := parser.ParseBytes(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            require.NoError(t, err)
            assert.Equal(t, tt.want.GroupId, got.GroupId)
        })
    }
}
```

### Testing Libraries

- `github.com/stretchr/testify` - Assertions
- `github.com/spf13/afero` - Filesystem mocking
- `github.com/charmbracelet/x/exp/teatest` - TUI testing

## Project Structure

```
haft/
├── cmd/
│   └── haft/
│       └── main.go           # Entry point
├── internal/
│   ├── cli/
│   │   ├── init/             # haft init command
│   │   └── root/             # Root command setup
│   ├── config/               # Configuration management
│   ├── generator/
│   │   ├── templates/        # Go embed templates
│   │   │   ├── project/      # Project scaffolding
│   │   │   └── resource/     # CRUD resource templates
│   │   └── engine.go         # Template engine
│   ├── initializr/           # Spring Initializr metadata
│   ├── logger/               # Logging utilities
│   ├── maven/                # pom.xml parser
│   └── tui/
│       ├── components/       # Reusable TUI components
│       ├── styles/           # Lip Gloss styles
│       └── wizard/           # Wizard orchestration
├── AGENTS.md                 # AI coding guidelines
├── CONTRIBUTING.md           # This file
├── go.mod
├── Makefile
└── README.md
```

### Architecture

```
CLI Layer (Cobra)
    ↓
TUI Layer (Bubble Tea)
    ↓
Generator Layer (Templates)
    ↓
Parser Layer (Maven/Gradle)
```

Dependencies flow inward. Inner layers should not depend on outer layers.

## Getting Help

- **Issues**: Open a [GitHub Issue](https://github.com/KashifKhn/haft/issues) for bugs or feature requests
- **Discussions**: Use [GitHub Discussions](https://github.com/KashifKhn/haft/discussions) for questions

## Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes for significant contributions

---

Thank you for contributing to Haft!
