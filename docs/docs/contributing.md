---
sidebar_position: 10
title: Contributing
description: How to contribute to Haft
---

# Contributing

We welcome contributions to Haft! This page provides a quick overview. For detailed guidelines, see [CONTRIBUTING.md](https://github.com/KashifKhn/haft/blob/main/CONTRIBUTING.md).

## Quick Start

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

## Code Standards

Haft follows strict coding standards:

### No Comments

Code should be self-documenting. Use clear function and variable names instead of comments.

```go
// Bad
// GetUserByID retrieves a user by their ID
func GetUserByID(id int) (*User, error) { ... }

// Good - self-explanatory
func GetUserByID(id int) (*User, error) { ... }
```

### DRY Principle

Never repeat code. Extract common logic into reusable functions.

### Size Limits

- **Functions**: Maximum 50 lines
- **Files**: Maximum 300 lines

### Naming

| Type | Convention | Example |
|------|------------|---------|
| Exported | PascalCase | `ParseConfig` |
| Private | camelCase | `parseFile` |
| Files | snake_case | `user_service.go` |

## Development Workflow

1. Fork the repository
2. Create a feature branch from `dev`
3. Make your changes
4. Run tests: `go test ./...`
5. Submit a pull request

## Commit Messages

Use conventional commits:

```
feat: add dependency search
fix: handle empty pom.xml
docs: update installation guide
refactor: simplify parser logic
```

## Testing

- Target 80%+ code coverage
- Use table-driven tests
- Mock filesystem with `afero`

```go
func TestParser_Parse(t *testing.T) {
    tests := []struct {
        name    string
        input   []byte
        want    *Project
        wantErr bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

## Getting Help

- [GitHub Issues](https://github.com/KashifKhn/haft/issues) — Bug reports
- [GitHub Discussions](https://github.com/KashifKhn/haft/discussions) — Questions

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
