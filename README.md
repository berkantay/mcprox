# mcprox

A robust, production-ready tool that retrieves and parses OpenAPI/Swagger documentation and generates a fully functional Model Context Protocol (MCP) proxy. MCProx makes your existing APIs instantly accessible to LLMs without any modification to your codebase.

## Features

- **OpenAPI/Swagger Integration**: Automatically fetches and parses Swagger documentation from any URL
- **Python MCP Server Generation**: Creates a fully functional MCP server in Python using modern best practices
- **Bridge Between LLMs and APIs**: Acts as a middleware layer that translates between LLM function calls and REST API endpoints
- **Real API Integration**: Makes actual HTTP requests to the original API, supporting all HTTP methods and authentication
- **Comprehensive Parsing**: Analyzes endpoints, methods, parameters, and schemas to create a complete MCP representation
- **Modern Python Structure**: Generates a well-structured Python project with virtual environment support
- **Production-Ready**: Built with robust error handling, logging, and configuration options
- **Well-Structured**: Organized codebase following best practices for Go and Python projects

## Quick Start

```bash
# Install the tool
go install github.com/berkantay/mcprox@latest

# Generate an MCP proxy from an OpenAPI spec
mcprox generate --url https://api.example.com/openapi.json --service-url https://api.example.com

# Set up the Python environment (using uv for faster dependency management)
cd generated_mcp_server
./scripts/setup.sh  # or scripts/setup.bat on Windows

# Run the generated server
source .venv/bin/activate  # or .venv\Scripts\activate.bat on Windows
python -m src.mcp_server
```

## Usage

```bash
# Basic usage
mcprox generate --url <swagger-url>

# Connect to the original API service
mcprox generate --url <swagger-url> --service-url <api-base-url>

# Add authentication to API requests
mcprox generate --url <swagger-url> --service-url <api-base-url> --service-auth "Bearer token123"

# Configure the output directory
mcprox generate --url <swagger-url> --output ./my-mcp-server

# Set timeout for HTTP requests
mcprox generate --url <swagger-url> --timeout 60
```

All configuration is done through command line flags. The available options are:

- `--url`, `-u`: URL to fetch OpenAPI documentation (required)
- `--timeout`, `-t`: Timeout in seconds for HTTP requests (default: 30)
- `--output`, `-o`: Output directory for generated server (default: ./generated)
- `--service-url`: Base URL of your API service
- `--service-auth`: Authorization header for API requests

## Architecture

The generated MCP proxy serves as a bridge between LLMs (like Claude, GPT, etc.) and your existing API:

```
[LLM] → [MCP Proxy] → [Original API]
```

1. The LLM sends a function call request to the MCP proxy
2. The MCP proxy validates the parameters based on OpenAPI specifications
3. The proxy translates the function call into an HTTP request to your API
4. Your API processes the request and returns a response
5. The MCP proxy formats the response and sends it back to the LLM

This architecture allows you to:

- Make your existing APIs instantly accessible to LLMs without modifying them
- Add a validation layer that ensures LLM requests are properly formatted
- Maintain separation between LLM interactions and your core business logic

## How It Works

MCProx works in several stages:

1. **OpenAPI Parsing**: Fetches and parses OpenAPI/Swagger documentation from the provided URL
2. **Schema Analysis**: Analyzes endpoints, methods, parameters, and schemas
3. **Code Generation**: Generates Python code for an MCP server, including:
   - Tool definitions that map to API endpoints
   - Parameter validation based on OpenAPI schemas
   - HTTP client for making real API requests
   - Error handling and logging
4. **Project Structure**: Creates a complete Python project structure with all necessary files

## Generated Server

The generated MCP server is a Python application that includes:

- **Complete Project Structure**: With `src`, `tests`, and `scripts` directories
- **Modern Python Tooling**: Using `pyproject.toml` for dependency management
- **Virtual Environment Support**: Setup scripts for easy environment creation
- **MCP Protocol Endpoint**: Available at `POST /api/mcp`
- **Health Check Endpoint**: Available at `GET /health`
- **Configuration Options**: Port configuration via environment variables
- **Real API Integration**: Automatically forwards requests to your original API

## Project Structure

The generated project follows modern Python best practices:

```
generated_mcp_server/
├── pyproject.toml      # Project metadata and dependencies
├── README.md           # Auto-generated documentation
├── .gitignore          # Git ignore file
├── scripts/            # Utility scripts
│   ├── setup.sh        # Unix setup script
│   ├── setup.bat       # Windows setup script
│   └── run.py          # Server run script
├── src/                # Source code
│   ├── __init__.py     # Package marker
│   └── mcp_server.py   # MCP server implementation
└── tests/              # Test directory
    └── __init__.py     # Package marker
```

## Environment Variables

The generated MCP server respects the following environment variables:

- `SERVICE_URL`: Base URL of the API service (default: http://localhost:8080)
- `PORT`: Port for the MCP server to listen on (default: 8000)

## Roadmap

Future plans for MCProx include:

- Support for more authentication methods
- Generation of client libraries in multiple languages
- More customization options for generated MCP servers
- Integration with local OpenAPI spec files
- Support for generating mock responses
- Improved schema validation and error handling
- Configuration via YAML file

## License

This project is licensed under the MIT License - see the LICENSE file for details.
