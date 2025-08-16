# Contributing to MCProx

Thank you for your interest in contributing to MCProx! This document provides guidelines and information to help you contribute effectively to the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Style and Standards](#code-style-and-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Project Structure](#project-structure)
- [Release Process](#release-process)
- [Getting Help](#getting-help)

## Getting Started

### Prerequisites

Before contributing to MCProx, ensure you have the following installed:

- **Go 1.23 or later** (the project uses toolchain go1.24.1)
- **Git** for version control
- **Make** for build automation
- **golangci-lint** (optional but recommended for linting)

### First Contribution

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/mcprox.git
   cd mcprox
   ```
3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://github.com/berkantay/mcprox.git
   ```
4. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Building the Project

```bash
# Build the binary
make build

# Or build and run directly
make run

# Clean build artifacts
make clean
```

The built binary will be available at `./build/mcprox`.

### Available Make Targets

- `make help` - Display all available targets
- `make build` - Build the binary
- `make test` - Run tests
- `make lint` - Run linters
- `make fmt` - Format code
- `make clean` - Remove build artifacts

### Dependencies

MCProx uses Go modules for dependency management. Key dependencies include:

- **Cobra** (`github.com/spf13/cobra`) - CLI framework
- **Fiber** (`github.com/gofiber/fiber/v2`) - Web framework for generated servers
- **Zap** (`go.uber.org/zap`) - Structured logging
- **Viper** (`github.com/spf13/viper`) - Configuration management
- **kin-openapi** (`github.com/getkin/kin-openapi`) - OpenAPI parsing
- **mcp-go** (`github.com/mark3labs/mcp-go`) - MCP protocol implementation

To add or update dependencies:

```bash
go get <package>
go mod tidy
```

## Code Style and Standards

### Go Code Style

- Follow the standard Go formatting (`gofmt`)
- Use `make fmt` to format code before committing
- Follow Go naming conventions (PascalCase for exported, camelCase for unexported)
- Write clear, descriptive variable and function names
- Add comments for exported functions and types
- Keep functions focused and reasonably sized

### Code Quality

- Run `make lint` before submitting changes
- Address all linter warnings and errors
- Use meaningful error messages with context
- Prefer explicit error handling over panic
- Use structured logging with the zap logger

### Package Organization

The project follows the standard Go project layout:

```
mcprox/
├── cmd/mcprox/          # Main application entry point
├── internal/            # Private application packages
│   ├── config/         # Configuration handling
│   ├── mcp/           # MCP-related functionality
│   └── openapi/       # OpenAPI parsing and processing
├── pkg/               # Public library packages
│   └── util/         # Utility functions
└── ...
```

### Naming Conventions

- **Files**: Use snake_case (e.g., `mcp_server.go`)
- **Packages**: Use short, lowercase names (e.g., `config`, `mcp`)
- **Functions/Methods**: Use camelCase for unexported, PascalCase for exported
- **Constants**: Use PascalCase for exported, camelCase for unexported
- **Interfaces**: Use descriptive names, often ending with `-er` (e.g., `Generator`)

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Writing Tests

- Write tests for all new functionality
- Use table-driven tests where appropriate
- Place tests in `*_test.go` files in the same package
- Use descriptive test names that explain what is being tested
- Test both happy paths and error cases
- Mock external dependencies when necessary

### Test Structure Example

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Submitting Changes

### Before Submitting

1. **Sync with upstream**:
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```

2. **Rebase your feature branch**:
   ```bash
   git checkout feature/your-feature-name
   git rebase main
   ```

3. **Run all checks**:
   ```bash
   make fmt
   make lint
   make test
   make build
   ```

4. **Test your changes** thoroughly, including:
   - Build the binary and test basic functionality
   - Test with different OpenAPI specifications
   - Verify generated MCP servers work correctly

### Commit Guidelines

- Use clear, descriptive commit messages
- Follow conventional commit format when possible:
  ```
  type(scope): description
  
  [optional body]
  
  [optional footer]
  ```
  
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Keep commits focused on a single change
- Squash related commits before submitting

### Pull Request Process

1. **Create a Pull Request** with:
   - Clear title describing the change
   - Detailed description of what was changed and why
   - Reference to any related issues
   - Screenshots or examples if applicable

2. **PR Template** (if you see checkboxes, please complete them):
   - [ ] Code follows the project's style guidelines
   - [ ] Self-review of the code has been performed
   - [ ] Tests have been added/updated for the changes
   - [ ] All tests pass
   - [ ] Documentation has been updated if needed

3. **Code Review Process**:
   - Maintainers will review your PR
   - Address any feedback or requested changes
   - Once approved, maintainers will merge the PR

### PR Requirements

- All CI checks must pass
- At least one maintainer approval required
- No merge conflicts with the main branch
- Clear commit history (squash if necessary)

## Project Structure

### Core Components

1. **CLI Interface** (`cmd/mcprox/`): Main application entry point using Cobra
2. **OpenAPI Parser** (`internal/openapi/`): Handles fetching and parsing OpenAPI specs
3. **MCP Generator** (`internal/mcp/`): Generates MCP server code from parsed OpenAPI
4. **Configuration** (`internal/config/`): Application configuration management
5. **Utilities** (`pkg/util/`): Shared utility functions

### Generated Output

The tool generates Python MCP servers with the following structure:
- Modern Python project layout
- Virtual environment support
- Dependency management with `pyproject.toml`
- Setup scripts for different platforms
- Complete MCP implementation

## Release Process

### Versioning

MCProx follows semantic versioning (SemVer):
- **Major** (X.0.0): Breaking changes
- **Minor** (0.X.0): New features, backward compatible
- **Patch** (0.0.X): Bug fixes, backward compatible

### Release Checklist

1. Update version in relevant files
2. Update CHANGELOG.md with new features and fixes
3. Create a release tag
4. Build and test release binaries
5. Update documentation if needed

## Getting Help

### Communication

- **Issues**: Use GitHub issues for bug reports and feature requests
- **Discussions**: Use GitHub discussions for questions and general discussion
- **Email**: For security issues, contact maintainers directly

### Issue Templates

When reporting issues:

**Bug Reports** should include:
- MCProx version
- Go version
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs or error messages

**Feature Requests** should include:
- Use case description
- Proposed solution
- Alternative solutions considered
- Additional context

### Development Questions

For development-related questions:
- Check existing issues and discussions
- Look at the code and tests for examples
- Ask specific questions with context

## Code of Conduct

Please note that this project follows a Code of Conduct. Be respectful and inclusive in all interactions.

## License

By contributing to MCProx, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to MCProx! Your contributions help make this tool better for everyone.