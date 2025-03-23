# MCP Generator Utilities

This package contains utility functions used by the MCP Generator SDK.

## Overview

The utilities package provides helper functions for various tasks related to MCP server generation, including:

1. String sanitization for tool IDs, parameter names, and package names
2. File generation for common project files
3. Script generation for setup and execution

## Functions

### String Sanitization

- `SanitizePathForToolID(path, method string) string`: Converts an OpenAPI path to a valid tool ID
- `SanitizeParamName(name string) string`: Converts an OpenAPI parameter name to a valid Python variable name
- `SanitizeForPackageName(name string) string`: Sanitizes a string to be used as a package name

### File Generation

- `GenerateRequirements(filePath string) error`: Generates a requirements.txt file for Python dependencies
- `GeneratePyprojectToml(filePath string, doc *openapi3.T) error`: Generates a pyproject.toml file for the project
- `GenerateGitignore(filePath string) error`: Generates a .gitignore file for the project
- `GenerateReadme(filePath string, doc *openapi3.T) error`: Generates a README.md file for the project
- `GenerateSetupScripts(outputDir string) error`: Generates setup scripts for the project
- `GenerateInitFiles(outputDir string) error`: Generates **init**.py files for Python package structure

## Usage Example

```go
package main

import (
    "github.com/berkantay/mcprox/internal/mcp/utils"
    "github.com/getkin/kin-openapi/openapi3"
)

func main() {
    // Example: Sanitize a path into a tool ID
    toolID := utils.SanitizePathForToolID("/users/{id}/posts", "GET")
    // Results in: "getUsers_idPosts"

    // Example: Sanitize a parameter name
    paramName := utils.SanitizeParamName("user-id")
    // Results in: "user_id"

    // Example: Generate project files
    doc := &openapi3.T{/* OpenAPI document */}
    utils.GenerateReadme("./README.md", doc)
    utils.GenerateRequirements("./requirements.txt")
}
```

## Development

When adding new utility functions:

1. Ensure the function has a clear, specific purpose
2. Add proper documentation with the function's purpose, parameters, and return values
3. Follow the existing naming and coding conventions
4. Consider adding tests for the new functionality

## License

MIT
