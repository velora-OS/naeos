# Contributing to NAEOS

Thank you for your interest in contributing to NAEOS! This document provides guidelines and information for contributors.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Style](#code-style)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Architecture Overview](#architecture-overview)

## Getting Started

1. Fork the repository
2. Clone your fork
3. Create a feature branch
4. Make your changes
5. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.25 or later
- Git
- Make (for build targets)
- golangci-lint (required — runs in CI and will block PRs)
- Docker (optional, for container builds)

### Setup

```bash
# Clone the repository
git clone https://github.com/NAEOS-foundation/naeos.git
cd naeos

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run linter (optional)
golangci-lint run ./...
```

## Code Style

### Go Code

- Follow standard Go conventions
- Use `gofmt` and `goimports` for formatting
- Add comments for exported functions and types
- Keep functions focused and small
- Handle errors explicitly

### Documentation

- Use Markdown for all documentation
- Follow the existing document structure
- Include examples where appropriate
- Keep language clear and concise

## Testing

### Writing Tests

- Place test files alongside the code they test
- Use table-driven tests where appropriate
- Test both success and error cases
- Aim for good coverage of critical paths

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v ./internal/neir/model/...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Pull Request Process

1. **Create a feature branch** from `main`
2. **Make your changes** with clear, focused commits
3. **Add tests** for new functionality
4. **Update documentation** if needed
5. **Run the test suite** to ensure nothing is broken
6. **Submit your pull request** with a clear description

### Commit Messages

- Use clear, descriptive commit messages
- Start with a verb in imperative mood
- Keep the first line under 72 characters
- Reference issues when applicable

Example:
```
Add multi-language SDK generation support

- Implement OutputAdapter interface pattern
- Add Go, TypeScript, Python, Java, Rust adapters
- Update pipeline to dispatch to adapters based on generation config
- Add CLI --language flag for language override

Closes #42
```

### PR Description

Please include:
- What the PR does
- Why the change is needed
- How it was tested
- Any breaking changes

## Architecture Overview

### Project Structure

```
naeos/
├── cmd/naeos/          # CLI entry point
├── pkg/                # Public packages
│   ├── kernel/         # Core kernel
│   ├── pipeline/       # Pipeline orchestrator
│   └── config/         # Configuration
├── internal/           # Private packages
│   ├── generation/     # Code generation
│   │   ├── engine/     # Default engine
│   │   └── adapters/   # Language adapters
│   ├── neir/           # NEIR model
│   ├── specification/  # Parser, normalizer, resolver
│   └── governance/     # Policy, review
├── specification/      # NAEOS specifications
├── docs/               # Documentation
└── governance/         # Governance documents
```

### Key Concepts

- **NEIR** (Nusantara Enterprise Intermediate Representation): The canonical model that all engines consume
- **OutputAdapter**: Interface for language-specific code generation
- **Pipeline**: Orchestrates the entire transformation from specification to artifacts

### Adding a New Language Adapter

1. Create a new file in `internal/generation/adapters/`
2. Implement the `OutputAdapter` interface
3. Register via `init()` function
4. Add tests
5. Update documentation

```go
type MyAdapter struct{}

func init() {
    Register(MyAdapter{})
}

func (MyAdapter) Language() language.Language {
    return "mylang"
}

func (MyAdapter) GenerateProject(projectName string) []engine.Artifact {
    // Return project-level artifacts
}
// ... implement other methods
```

## Questions?

If you have questions about contributing, feel free to:
- Open an issue
- Start a discussion
- Reach out to maintainers

Thank you for contributing to NAEOS!
