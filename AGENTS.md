# AGENTS.md - Coding Agent Guidelines for Haft

## Build/Test Commands
- `go build ./...` - Build all packages
- `go test ./...` - Run all tests
- `go test ./internal/generator/... -v` - Test specific package
- `go test -run TestFunctionName ./...` - Run single test
- `go test -cover ./...` - Run tests with coverage (target: 80%+)

## Code Standards
- **NO COMMENTS**: No code comments or doc comments unless absolutely necessary
- **Self-Documenting Code**: Use clear function/variable names; code explains itself
- **DRY Principle**: Never repeat code; extract common logic into reusable functions
- **Naming**: PascalCase (exported), camelCase (private), snake_case (files)
- **Errors**: Return errors with context, never panic; provide actionable messages
- **Functions**: Max 50 lines; single responsibility; split if larger
- **Files**: Max 300 lines; one concept per file

## CLI Documentation
- Every command MUST have Short, Long, and Example in Cobra definition
- CLI is the primary documentation; `--help` must be comprehensive
- Update Docusaurus docs incrementally with each feature

## Dependencies
- Cobra (CLI), Bubble Tea/Bubbles/Lip Gloss (TUI), charmbracelet/log (logging)
- testify (assertions), afero (filesystem mocking), teatest (TUI testing)

## Architecture
- CLI → TUI → Generator → Parser (dependencies flow inward)
- Templates embedded via Go embed; no hardcoded strings in generators
