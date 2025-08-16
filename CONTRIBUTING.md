# Contributing to MCProx

Thank you for your interest in contributing to MCProx! This document provides guidelines and information for contributors who want to help improve this Model Context Protocol proxy generator.

## Table of Contents

- [Project Overview](#project-overview)
- [Development Setup](#development-setup)
- [Development Workflow](#development-workflow)
- [Code Standards](#code-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Issue Reporting](#issue-reporting)
- [Project Structure](#project-structure)
- [Release Process](#release-process)
- [Getting Help](#getting-help)

## Project Overview

MCProx is a Go-based CLI tool that retrieves and parses OpenAPI/Swagger documentation to generate fully functional Model Context Protocol (MCP) proxy servers in Python. The tool acts as a bridge between LLMs and existing REST APIs without requiring modifications to your existing codebase.

### Key Components

- **OpenAPI Parser**: Fetches and analyzes Swagger/OpenAPI specifications
- **Code Generator**: Creates Python MCP server projects with proper structure
- **HTTP Client**: Enables real API integration with authentication support
- **CLI Interface**: User-friendly command-line interface using Cobra

## Development Setup

### Prerequisites

- **Go**: Version 1.23 or later (as specified in `go.mod`)
- **Git**: For version control
- **Make**: For build automation (optional but recommended)
- **golangci-lint**: For linting (optional but recommended)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/berkantay/mcprox.git
   cd mcprox
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Install development tools (optional):**
   ```bash
   # Install golangci-lint for comprehensive linting
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

4. **Build the project:**
   ```bash
   make build
   # or
   go build -o build/mcprox
   ```

5. **Verify installation:**
   ```bash
   ./build/mcprox --help
   ```

## Development Workflow

### Making Changes

1. **Fork the repository** on GitHub
2. **Create a feature branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes** following the code standards
4. **Test your changes** thoroughly
5. **Commit your changes** with descriptive messages
6. **Push to your fork** and create a pull request

### Available Make Commands

```bash
make help        # Display all available commands
make build       # Build the binary
make test        # Run tests
make lint        # Run linters
make fmt         # Format code
make clean       # Clean build artifacts
make run         # Run the application directly
```

### Running the Development Version

```bash
# Run directly with go run
go run cmd/mcprox/main.go generate --url https://example.com/openapi.json

# Or build and run
make build
./build/mcprox generate --url https://example.com/openapi.json
```

## Code Standards

### Go Code Style

We follow standard Go conventions and best practices:

- **Formatting**: Use `gofmt` or `make fmt` to format code
- **Linting**: Use `golangci-lint` or `make lint` for comprehensive checks
- **Naming**: Follow Go naming conventions (camelCase, PascalCase as appropriate)
- **Documentation**: All public functions and types should have comments
- **Error Handling**: Always handle errors appropriately, don't ignore them

### Code Organization

```
mcprox/
├── cmd/mcprox/          # Main application entry point
├── internal/            # Internal packages (not importable by other projects)
│   ├── config/         # Configuration handling
│   ├── mcp/            # MCP-related functionality
│   └── openapi/        # OpenAPI parsing and handling
├── pkg/util/           # Public utility packages
└── build/              # Build outputs
```

### Commit Message Format

Use clear, descriptive commit messages:

```
feat: add support for OAuth2 authentication
fix: handle empty OpenAPI descriptions gracefully
docs: update README with new configuration options
test: add unit tests for OpenAPI parser
refactor: simplify HTTP client initialization
```

### Code Guidelines

- **Keep functions small** and focused on a single responsibility
- **Use meaningful variable names** that clearly express purpose
- **Add unit tests** for new functionality
- **Handle edge cases** and error conditions
- **Use structured logging** with the existing zap logger
- **Follow the existing patterns** in the codebase

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/openapi -v
```

### Writing Tests

- **Unit tests**: Test individual functions and methods in isolation
- **Integration tests**: Test interactions between components
- **Table-driven tests**: Use for testing multiple scenarios
- **Mock external dependencies**: Don't rely on external services in tests

Example test structure:
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Submitting Changes

### Pull Request Process

1. **Ensure your code passes all checks:**
   ```bash
   make lint
   make test
   make build
   ```

2. **Update documentation** if you're changing APIs or adding features

3. **Write clear PR description** explaining:
   - What changes you made
   - Why you made them
   - How to test the changes
   - Any breaking changes

4. **Link related issues** if applicable

### PR Guidelines

- **Keep PRs focused** on a single feature or fix
- **Include tests** for new functionality
- **Update documentation** as needed
- **Ensure CI passes** before requesting review
- **Be responsive** to feedback and suggestions

### What We Look For

- **Code quality**: Clean, readable, maintainable code
- **Testing**: Adequate test coverage for new features
- **Documentation**: Clear documentation for new APIs
- **Backwards compatibility**: Avoid breaking changes when possible
- **Performance**: Consider performance implications of changes

## Issue Reporting

### Bug Reports

When reporting bugs, please include:

- **MCProx version**: `mcprox --version`
- **Go version**: `go version`
- **Operating system**: OS and version
- **Steps to reproduce**: Clear, minimal reproduction steps
- **Expected behavior**: What should happen
- **Actual behavior**: What actually happens
- **Error messages**: Full error output if applicable
- **Sample OpenAPI spec**: If related to parsing issues

### Feature Requests

For feature requests, please include:

- **Use case**: Describe your specific use case
- **Proposed solution**: Your ideas for implementation
- **Alternatives considered**: Other approaches you've thought about
- **Additional context**: Any other relevant information

### Security Issues

For security vulnerabilities, please:

- **Do not** open a public issue
- **Email** the maintainers directly
- **Include** detailed information about the vulnerability
- **Wait** for a response before disclosing publicly

## Project Structure

### Key Directories

- **`cmd/mcprox/`**: Main application entry point and CLI setup
- **`internal/config/`**: Configuration management and validation
- **`internal/mcp/`**: MCP server generation and protocol handling
- **`internal/openapi/`**: OpenAPI specification parsing and validation
- **`pkg/util/`**: Shared utilities that could be used by other projects

### Key Files

- **`go.mod/go.sum`**: Go module dependencies
- **`Makefile`**: Build automation and development tasks
- **`.gitignore`**: Files excluded from version control
- **`README.md`**: Project overview and usage instructions

### Adding New Dependencies

When adding new dependencies:

1. **Use `go get`** to add the dependency
2. **Run `go mod tidy`** to clean up the module file
3. **Justify the addition** in your PR description
4. **Prefer standard library** when possible
5. **Choose well-maintained** and popular packages

## Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):

- **Major**: Breaking changes
- **Minor**: New features (backwards compatible)
- **Patch**: Bug fixes (backwards compatible)

### Release Checklist

1. Update version information
2. Update CHANGELOG.md
3. Ensure all tests pass
4. Create and test release build
5. Tag the release
6. Update documentation
7. Announce the release

## Getting Help

### Communication Channels

- **Issues**: For bug reports and feature requests
- **Discussions**: For questions and general discussion
- **Pull Requests**: For code contributions

### Resources

- **Go Documentation**: https://golang.org/doc/
- **MCP Specification**: https://modelcontextprotocol.io/
- **OpenAPI Specification**: https://swagger.io/specification/
- **Project README**: [README.md](README.md)

### Before Asking for Help

1. **Check existing issues** and discussions
2. **Read the documentation** thoroughly  
3. **Try debugging** the issue yourself
4. **Prepare a minimal reproduction** case

## Code of Conduct

This project follows the [Go Community Code of Conduct](https://golang.org/conduct/). By participating, you agree to uphold this code. Please report any unacceptable behavior to the project maintainers.

## License

By contributing to MCProx, you agree that your contributions will be licensed under the same license as the project (MIT License).

---

Thank you for contributing to MCProx! Your help makes this tool better for everyone. 🚀