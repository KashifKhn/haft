# Contributing to Haft

Thank you for your interest in contributing to Haft.

## Getting Started

```bash
# Fork and clone
git clone https://github.com/YOUR_USERNAME/haft.git
cd haft

# Add upstream
git remote add upstream https://github.com/KashifKhn/haft.git

# Install and build
go mod download
make build

# Run tests
make test
```

## Development

```bash
make build        # Build binary
make test         # Run tests
make test-cover   # Run tests with coverage
make lint         # Run linter
make fmt          # Format code
```

## Code Standards

Read [AGENTS.md](AGENTS.md) for complete guidelines. Key points:

- **No comments** — Code should be self-documenting
- **DRY** — Don't repeat yourself
- **Max 50 lines** per function
- **Max 300 lines** per file
- **Return errors with context** — Never panic

## Pull Requests

1. Create a branch from `dev`
2. Make your changes
3. Run `make test` and `make lint`
4. Submit PR against `dev`

### Commit Format

```
type: description

feat: add dependency search
fix: handle empty pom.xml
docs: update installation guide
refactor: simplify parser logic
test: add generator tests
```

## Project Structure

```
haft/
├── cmd/haft/           # Entry point
├── internal/
│   ├── cli/            # Commands (init, generate, add)
│   ├── config/         # Configuration
│   ├── generator/      # Template engine
│   ├── maven/          # pom.xml parser
│   └── tui/            # Terminal UI components
├── docs/               # Documentation site
└── assets/             # Logos and images
```

## Need Help?

- [Documentation](https://kashifkhn.github.io/haft)
- [Issues](https://github.com/KashifKhn/haft/issues)
- [Discussions](https://github.com/KashifKhn/haft/discussions)
