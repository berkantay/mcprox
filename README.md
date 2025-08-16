# MCProx

[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/berkantay/mcprox)](https://goreportcard.com/report/github.com/berkantay/mcprox)

> **Transform any OpenAPI/Swagger API into an LLM-ready Model Context Protocol (MCP) proxy**

MCProx is a robust, production-ready tool that automatically retrieves and parses OpenAPI/Swagger documentation and generates a fully functional Model Context Protocol (MCP) server. It makes your existing APIs instantly accessible to LLMs without requiring any modifications to your codebase.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
- [Command Reference](#command-reference)
- [Architecture](#architecture)
- [Generated Server](#generated-server)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [Roadmap](#roadmap)
- [License](#license)

## Features

✨ **OpenAPI/Swagger Integration** - Automatically fetches and parses Swagger documentation from any URL  
🐍 **Python MCP Server Generation** - Creates production-ready MCP servers using modern Python best practices  
🔗 **Bridge Between LLMs and APIs** - Acts as middleware translating LLM function calls to REST API endpoints  
🚀 **Real API Integration** - Makes actual HTTP requests supporting all methods and authentication schemes  
📝 **Comprehensive Parsing** - Analyzes endpoints, methods, parameters, and schemas for complete MCP representation  
🏗️ **Modern Structure** - Generates well-structured projects with virtual environment support  
🛡️ **Production-Ready** - Built with robust error handling, logging, and configuration options  
🧪 **Testing Support** - Includes test structures and development tooling  

## Prerequisites

### For Using MCProx

- **Go 1.23+** - [Install Go](https://golang.org/doc/install)
- **Git** - For installation from source

### For Generated MCP Servers

- **Python 3.8+** - Required for running generated MCP servers
- **uv** (recommended) or **pip** - For Python package management
  - Install uv: `pip install uv` or see [uv documentation](https://docs.astral.sh/uv/)

## Installation

### Option 1: Install from Go (Recommended)

```bash
go install github.com/berkantay/mcprox@latest
```

### Option 2: Build from Source

```bash
git clone https://github.com/berkantay/mcprox.git
cd mcprox
make build
# Binary will be available at ./build/mcprox
```

### Option 3: Download Release

Visit the [releases page](https://github.com/berkantay/mcprox/releases) and download the binary for your platform.

## Quick Start

Here's how to get started in under 2 minutes:

```bash
# 1. Generate MCP proxy from a public API
mcprox generate --url https://petstore.swagger.io/v2/swagger.json \
                --service-url https://petstore.swagger.io/v2

# 2. Set up the Python environment
cd generated_mcp_server
chmod +x scripts/setup.sh
./scripts/setup.sh

# 3. Activate environment and run the server
source .venv/bin/activate
python -m src.mcp_server

# 4. Test the server (in another terminal)
curl -X POST http://localhost:8000/api/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/list"}'
```

## Usage Examples

### Example 1: JSONPlaceholder API

Generate an MCP proxy for the JSONPlaceholder testing API:

```bash
mcprox generate \
  --url https://jsonplaceholder.typicode.com/swagger.json \
  --service-url https://jsonplaceholder.typicode.com \
  --output ./jsonplaceholder-mcp
```

### Example 2: API with Authentication

For APIs requiring authentication:

```bash
mcprox generate \
  --url https://api.github.com/swagger.json \
  --service-url https://api.github.com \
  --service-auth "Bearer ghp_your_token_here" \
  --output ./github-mcp
```

### Example 3: Local Development API

Working with a local API during development:

```bash
mcprox generate \
  --url http://localhost:3000/api-docs/swagger.json \
  --service-url http://localhost:3000 \
  --timeout 60 \
  --output ./local-api-mcp
```

### Example 4: Corporate API with Custom Headers

```bash
mcprox generate \
  --url https://internal-api.company.com/swagger.json \
  --service-url https://internal-api.company.com \
  --service-auth "ApiKey your-api-key-here" \
  --timeout 120 \
  --output ./corporate-mcp
```

## Command Reference

### `mcprox generate`

Generates an MCP proxy from an OpenAPI/Swagger specification.

#### Required Flags

- `--url`, `-u` - URL to fetch OpenAPI documentation from

#### Optional Flags

- `--output`, `-o` - Output directory for generated server (default: `./generated_mcp_server`)
- `--service-url` - Base URL of the target API service  
- `--service-auth` - Authorization header value for API requests
- `--timeout`, `-t` - HTTP request timeout in seconds (default: 30)

#### Examples

```bash
# Minimal usage
mcprox generate -u https://api.example.com/swagger.json

# Full configuration
mcprox generate \
  --url https://api.example.com/swagger.json \
  --service-url https://api.example.com \
  --service-auth "Bearer token123" \
  --output ./my-mcp-server \
  --timeout 60
```

### `mcprox help`

Display help information for any command:

```bash
mcprox help generate
```

## Architecture

MCProx creates a bridge between LLMs and your existing APIs:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│       LLM       │◄──►│   MCP Proxy     │◄──►│  Original API   │
│  (Claude/GPT)   │    │  (Generated)    │    │ (Your Service)  │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                        │                        │
        │                        │                        │
        ▼                        ▼                        ▼
   Function calls         Parameter           HTTP requests
   via MCP protocol       validation &        to real API
                         transformation
```

### How it Works

1. **OpenAPI Analysis** 📋 - MCProx fetches and parses your OpenAPI/Swagger specification
2. **Code Generation** 🔧 - Generates a Python MCP server with tools mapped to API endpoints  
3. **LLM Integration** 🤖 - LLMs send function calls to the MCP proxy via standard protocol
4. **Request Translation** 🔄 - Proxy validates parameters and converts to HTTP requests
5. **API Communication** 🌐 - Makes real requests to your original API endpoints
6. **Response Formatting** 📤 - Formats API responses for LLM consumption

### Benefits

- **Zero API Modifications** - Your existing APIs work without changes
- **Type Safety** - Parameter validation based on OpenAPI schemas
- **Authentication Support** - Handles various auth methods transparently
- **Error Handling** - Robust error handling and informative error messages
- **Scalability** - Generated servers can handle production workloads

## Generated Server

The MCP proxy generates a complete Python application with the following structure:

```
generated_mcp_server/
├── 📄 README.md              # Auto-generated documentation
├── 📄 pyproject.toml          # Python project configuration  
├── 📄 .gitignore             # Git ignore patterns
├── 📁 scripts/               # Setup and utility scripts
│   ├── 🔧 setup.sh           # Unix environment setup
│   ├── 🔧 setup.bat          # Windows environment setup  
│   └── 🐍 run.py             # Development server runner
├── 📁 src/                   # Source code directory
│   ├── 📄 __init__.py        # Package initialization
│   └── 🐍 mcp_server.py      # Main MCP server implementation
└── 📁 tests/                 # Test directory structure
    └── 📄 __init__.py        # Test package initialization
```

### Server Features

#### 🌐 **HTTP Endpoints**
- `POST /api/mcp` - MCP protocol endpoint for LLM communication
- `GET /health` - Health check endpoint for monitoring

#### ⚙️ **Configuration**
Environment variables for runtime configuration:
- `PORT` - Server port (default: 8000)
- `SERVICE_URL` - Target API base URL
- `LOG_LEVEL` - Logging level (DEBUG, INFO, WARNING, ERROR)

#### 🔧 **Setup Scripts**
Platform-specific setup scripts handle:
- Virtual environment creation
- Dependency installation using uv (faster) or pip
- Development environment configuration

## Development

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/berkantay/mcprox.git
cd mcprox

# Install dependencies
go mod download

# Run tests
make test

# Format code
make fmt

# Run linting
make lint

# Build the project
make build
```

### Available Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build the binary to `./build/mcprox` |
| `make clean` | Remove build artifacts and clean Go cache |
| `make test` | Run all tests with verbose output |
| `make fmt` | Format code using gofmt |
| `make lint` | Run Go vet and golangci-lint |
| `make run` | Run from source (development) |
| `make help` | Display all available targets |

### Project Structure

```
mcprox/
├── 📁 cmd/mcprox/            # Main application entry point
├── 📁 internal/              # Internal packages (not imported externally)
│   ├── 📁 config/            # Configuration management
│   ├── 📁 mcp/               # MCP protocol implementation
│   └── 📁 openapi/           # OpenAPI parsing and processing
├── 📁 pkg/                   # Public packages (can be imported)
├── 📄 go.mod                 # Go module definition
├── 📄 go.sum                 # Dependency checksums
└── 📄 Makefile              # Build automation
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/openapi/

# Run tests in verbose mode
go test -v ./...
```

## Troubleshooting

### Common Issues

#### ❌ "Failed to fetch OpenAPI spec"

**Problem:** MCProx cannot retrieve the OpenAPI specification from the provided URL.

**Solutions:**
1. Verify the URL is accessible: `curl -I https://your-api.com/swagger.json`
2. Check if authentication is required for the documentation endpoint
3. Increase timeout: `--timeout 60`
4. Verify the URL returns valid JSON/YAML content

#### ❌ "Invalid OpenAPI specification"

**Problem:** The fetched specification is not valid OpenAPI/Swagger format.

**Solutions:**
1. Validate your OpenAPI spec using [Swagger Editor](https://editor.swagger.io/)
2. Ensure the Content-Type header is correct (application/json or application/yaml)
3. Check for malformed JSON/YAML syntax

#### ❌ "Generated server fails to start"

**Problem:** The generated MCP server encounters errors during startup.

**Solutions:**
1. Check Python version: `python --version` (requires 3.8+)
2. Verify virtual environment setup: `source .venv/bin/activate`
3. Install dependencies: `pip install -r requirements.txt` or use setup scripts
4. Check port availability: `lsof -i :8000`

#### ❌ "API requests return authentication errors"

**Problem:** The generated proxy cannot authenticate with the target API.

**Solutions:**
1. Verify the `--service-auth` flag format matches your API requirements
2. Check if the token/key has expired
3. Ensure the auth header format is correct (e.g., "Bearer token", "ApiKey key")
4. Test authentication directly: `curl -H "Authorization: Bearer token" https://api.example.com`

### Debug Mode

Enable debug logging for detailed troubleshooting:

```bash
# Set debug environment variable before running
export LOG_LEVEL=DEBUG
python -m src.mcp_server
```

### Getting Help

If you encounter issues not covered here:

1. Check [existing issues](https://github.com/berkantay/mcprox/issues)
2. Create a [new issue](https://github.com/berkantay/mcprox/issues/new) with:
   - MCProx version (`mcprox --version`)
   - Go version (`go version`)
   - Operating system
   - Complete command used
   - Error messages and logs

## Contributing

We welcome contributions! Here's how to get started:

### 🚀 Quick Contribution Guide

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/yourusername/mcprox.git`
3. **Create a branch**: `git checkout -b feature/amazing-feature`
4. **Make changes** and test them: `make test`
5. **Format code**: `make fmt`
6. **Commit changes**: `git commit -m "Add amazing feature"`
7. **Push to branch**: `git push origin feature/amazing-feature`  
8. **Open a Pull Request**

### 📝 Contribution Guidelines

- **Code Style**: Follow Go best practices and run `make fmt`
- **Testing**: Add tests for new features and ensure `make test` passes
- **Documentation**: Update README and code comments for significant changes
- **Commits**: Use clear, descriptive commit messages
- **Issues**: Reference issue numbers in PRs when applicable

### 🎯 Areas for Contribution

- **Authentication Methods** - Support for OAuth, API keys, etc.
- **Output Formats** - Additional language generators beyond Python
- **Schema Validation** - Enhanced parameter validation  
- **Error Handling** - Improved error messages and recovery
- **Documentation** - Examples, tutorials, and API documentation
- **Testing** - Integration tests and test coverage improvements

### 📋 Development Checklist

Before submitting a PR:

- [ ] Code builds successfully (`make build`)
- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] Documentation is updated
- [ ] New features have tests
- [ ] Breaking changes are documented

## Roadmap

### 🎯 Near Term (Next Release)

- [ ] **Configuration Files** - YAML/JSON config file support
- [ ] **Mock Mode** - Generate responses without calling real APIs
- [ ] **Schema Validation** - Enhanced request/response validation  
- [ ] **Auth Methods** - OAuth 2.0 and JWT token support
- [ ] **Local Files** - Support for local OpenAPI spec files

### 🚀 Medium Term

- [ ] **Multi-Language Support** - Generate servers in Go, Node.js, Rust
- [ ] **CLI Improvements** - Interactive mode and better error messages
- [ ] **Monitoring** - Built-in metrics and health monitoring
- [ ] **Caching** - Response caching and rate limiting
- [ ] **WebSockets** - Support for real-time API endpoints

### 🌟 Long Term

- [ ] **Visual Interface** - Web UI for proxy management
- [ ] **Plugin System** - Custom transformations and middleware
- [ ] **API Gateway** - Full-featured API gateway capabilities
- [ ] **Multi-API** - Single proxy for multiple APIs
- [ ] **Deployment** - Docker containers and cloud deployment tools

### 💡 Ideas & Suggestions

Have an idea? [Open an issue](https://github.com/berkantay/mcprox/issues/new) with the "enhancement" label!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Built with ❤️ by the MCProx team**

[⭐ Star us on GitHub](https://github.com/berkantay/mcprox) • [🐛 Report Issues](https://github.com/berkantay/mcprox/issues) • [💡 Request Features](https://github.com/berkantay/mcprox/issues/new)

</div>