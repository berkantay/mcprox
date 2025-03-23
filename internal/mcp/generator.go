package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/berkantay/mcprox/internal/config"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// Generator handles the creation of MCP server from OpenAPI specs
type Generator struct {
	logger    *zap.Logger
	outputDir string
	document  *openapi3.T
}

// NewGenerator creates a new MCP generator
func NewGenerator(logger *zap.Logger) *Generator {
	return &Generator{
		logger:    logger,
		outputDir: config.GetString("output.dir"),
	}
}

// Generate generates an MCP server from an OpenAPI spec
func (g *Generator) Generate(ctx context.Context, doc *openapi3.T) error {
	g.logger.Info("Generating MCP server from OpenAPI documentation")

	// Store the document in the generator
	g.document = doc

	// Set up project directory
	projectDir := filepath.Join(g.outputDir, "generated_mcp_server")
	g.outputDir = projectDir

	// Create project directory structure
	if err := g.createProjectStructure(); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	// Create MCP server
	mcpServer := server.NewMCPServer(
		doc.Info.Title,
		doc.Info.Version,
	)

	// Process paths into tools
	if err := g.processPathsIntoTools(doc, mcpServer); err != nil {
		return err
	}

	// Generate server code
	serverPath := filepath.Join(g.outputDir, "src", "mcp_server.py")
	if err := g.generateServerCode(mcpServer, serverPath); err != nil {
		return fmt.Errorf("failed to generate server code: %w", err)
	}

	// Generate requirements.txt
	requirementsPath := filepath.Join(g.outputDir, "requirements.txt")
	if err := g.generateRequirements(requirementsPath); err != nil {
		return fmt.Errorf("failed to generate requirements.txt: %w", err)
	}

	// Generate pyproject.toml
	pyprojectPath := filepath.Join(g.outputDir, "pyproject.toml")
	if err := g.generatePyprojectToml(pyprojectPath, doc); err != nil {
		return fmt.Errorf("failed to generate pyproject.toml: %w", err)
	}

	// Generate .gitignore
	gitignorePath := filepath.Join(g.outputDir, ".gitignore")
	if err := g.generateGitignore(gitignorePath); err != nil {
		return fmt.Errorf("failed to generate .gitignore: %w", err)
	}

	// Generate README.md
	readmePath := filepath.Join(g.outputDir, "README.md")
	if err := g.generateReadme(readmePath, doc); err != nil {
		return fmt.Errorf("failed to generate README.md: %w", err)
	}

	// Generate setup scripts
	if err := g.generateSetupScripts(); err != nil {
		return fmt.Errorf("failed to generate setup scripts: %w", err)
	}

	// Generate __init__.py files for package structure
	if err := g.generateInitFiles(); err != nil {
		return fmt.Errorf("failed to generate __init__.py files: %w", err)
	}

	g.logger.Info("Successfully generated MCP server project",
		zap.String("project_dir", projectDir))

	return nil
}

// createProjectStructure creates the directory structure for the Python project
func (g *Generator) createProjectStructure() error {
	dirs := []string{
		g.outputDir,
		filepath.Join(g.outputDir, "src"),
		filepath.Join(g.outputDir, "tests"),
		filepath.Join(g.outputDir, "scripts"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateInitFiles generates __init__.py files for Python package structure
func (g *Generator) generateInitFiles() error {
	initFiles := []string{
		filepath.Join(g.outputDir, "src", "__init__.py"),
		filepath.Join(g.outputDir, "tests", "__init__.py"),
	}

	for _, file := range initFiles {
		if err := os.WriteFile(file, []byte("# Auto-generated by mcprox\n"), 0644); err != nil {
			return fmt.Errorf("failed to create __init__.py file at %s: %w", file, err)
		}
	}

	return nil
}

// generateSetupScripts generates setup scripts for the project
func (g *Generator) generateSetupScripts() error {
	// Generate setup.sh (for Unix-based systems)
	setupShPath := filepath.Join(g.outputDir, "scripts", "setup.sh")
	setupShContent := `#!/bin/bash
# Setup script for MCP server

# Check if uv is installed
if ! command -v uv &> /dev/null; then
    echo "uv not found, installing..."
    curl -LsSf https://astral.sh/uv/install.sh | sh
fi

# Create virtual environment and install dependencies
cd "$(dirname "$0")/.."
uv venv
uv pip install -e .
echo "Setup complete. Run 'source .venv/bin/activate' to activate the environment."
`
	if err := os.WriteFile(setupShPath, []byte(setupShContent), 0755); err != nil {
		return fmt.Errorf("failed to generate setup.sh: %w", err)
	}

	// Generate setup.bat (for Windows)
	setupBatPath := filepath.Join(g.outputDir, "scripts", "setup.bat")
	setupBatContent := `@echo off
REM Setup script for MCP server

REM Check if uv is installed
where uv >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo uv not found, please install it manually from https://astral.sh/uv
    exit /b 1
)

REM Create virtual environment and install dependencies
cd %~dp0\..
uv venv
uv pip install -e .
echo Setup complete. Run '.venv\Scripts\activate.bat' to activate the environment.
`
	if err := os.WriteFile(setupBatPath, []byte(setupBatContent), 0644); err != nil {
		return fmt.Errorf("failed to generate setup.bat: %w", err)
	}

	// Generate run script
	runScriptPath := filepath.Join(g.outputDir, "scripts", "run.py")
	runScriptContent := `#!/usr/bin/env python3
"""
Run script for MCP server.
"""
import os
import sys
import subprocess

def main():
    """Run the MCP server."""
    # Get the project root directory
    project_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    
    # Check if virtual environment exists
    venv_dir = os.path.join(project_dir, ".venv")
    if not os.path.exists(venv_dir):
        print("Virtual environment not found. Running setup...")
        setup_script = os.path.join(project_dir, "scripts", "setup.sh")
        if os.name == "nt":  # Windows
            setup_script = os.path.join(project_dir, "scripts", "setup.bat")
        
        subprocess.call(setup_script, shell=True)
    
    # Run the MCP server
    server_script = os.path.join(project_dir, "src", "mcp_server.py")
    
    # Determine python command (use venv python if available)
    python_cmd = "python"
    if os.name == "nt":  # Windows
        python_path = os.path.join(venv_dir, "Scripts", "python.exe")
    else:  # Unix-like
        python_path = os.path.join(venv_dir, "bin", "python")
    
    if os.path.exists(python_path):
        python_cmd = python_path
    
    # Run the server
    subprocess.call([python_cmd, server_script])

if __name__ == "__main__":
    main()
`
	if err := os.WriteFile(runScriptPath, []byte(runScriptContent), 0755); err != nil {
		return fmt.Errorf("failed to generate run.py: %w", err)
	}

	return nil
}

// generatePyprojectToml generates a pyproject.toml file for the project
func (g *Generator) generatePyprojectToml(filePath string, doc *openapi3.T) error {
	projectName := sanitizeForPackageName(doc.Info.Title)
	if projectName == "" {
		projectName = "mcp_server"
	}

	content := fmt.Sprintf(`[build-system]
requires = ["setuptools>=61.0"]
build-backend = "setuptools.build_meta"

[project]
name = "%s"
version = "%s"
authors = [
    {name = "Generated by mcprox", email = "example@example.com"},
]
description = "Model Context Protocol (MCP) server generated from OpenAPI specs"
readme = "README.md"
requires-python = ">=3.8"
classifiers = [
    "Programming Language :: Python :: 3",
    "License :: OSI Approved :: MIT License",
    "Operating System :: OS Independent",
]
dependencies = [
    "mcp",
    "requests",
]

[project.optional-dependencies]
dev = [
    "pytest",
    "black",
    "ruff",
]

[project.urls]
"Homepage" = "https://github.com/yourusername/mcprox"
"Bug Tracker" = "https://github.com/yourusername/mcprox/issues"

[tool.setuptools]
package-dir = {"" = "src"}

[tool.ruff]
line-length = 100
target-version = "py38"

[tool.black]
line-length = 100
target-version = ["py38"]
`, projectName, doc.Info.Version)

	return os.WriteFile(filePath, []byte(content), 0644)
}

// generateGitignore generates a .gitignore file for the project
func (g *Generator) generateGitignore(filePath string) error {
	content := `# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
env/
build/
develop-eggs/
dist/
downloads/
eggs/
.eggs/
lib/
lib64/
parts/
sdist/
var/
*.egg-info/
.installed.cfg
*.egg

# Virtual Environment
.env
.venv
venv/
ENV/
.uv/

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
*.log
`
	return os.WriteFile(filePath, []byte(content), 0644)
}

// generateReadme generates a README.md file for the project
func (g *Generator) generateReadme(filePath string, doc *openapi3.T) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s MCP Server\n\n", doc.Info.Title))
	sb.WriteString(fmt.Sprintf("This is an auto-generated Model Context Protocol (MCP) server for %s (version %s).\n\n", doc.Info.Title, doc.Info.Version))

	sb.WriteString("## Description\n\n")
	sb.WriteString(doc.Info.Description)
	sb.WriteString("\n\n")

	sb.WriteString("## Installation\n\n")
	sb.WriteString("### Using uv (recommended)\n\n")
	sb.WriteString("This project uses [uv](https://astral.sh/uv) for dependency management and virtual environments.\n\n")

	sb.WriteString("1. Install uv (if not already installed):\n")
	sb.WriteString("   ```bash\n")
	sb.WriteString("   curl -LsSf https://astral.sh/uv/install.sh | sh\n")
	sb.WriteString("   ```\n\n")

	sb.WriteString("2. Run the setup script:\n")
	sb.WriteString("   ```bash\n")
	sb.WriteString("   # On Unix/Linux/MacOS\n")
	sb.WriteString("   ./scripts/setup.sh\n")
	sb.WriteString("   \n")
	sb.WriteString("   # On Windows\n")
	sb.WriteString("   scripts\\setup.bat\n")
	sb.WriteString("   ```\n\n")

	sb.WriteString("3. Activate the virtual environment:\n")
	sb.WriteString("   ```bash\n")
	sb.WriteString("   # On Unix/Linux/MacOS\n")
	sb.WriteString("   source .venv/bin/activate\n")
	sb.WriteString("   \n")
	sb.WriteString("   # On Windows\n")
	sb.WriteString("   .venv\\Scripts\\activate.bat\n")
	sb.WriteString("   ```\n\n")

	sb.WriteString("### Using pip\n\n")
	sb.WriteString("Alternatively, you can use pip:\n\n")

	sb.WriteString("1. Create a virtual environment:\n")
	sb.WriteString("   ```bash\n")
	sb.WriteString("   python -m venv .venv\n")
	sb.WriteString("   ```\n\n")

	sb.WriteString("2. Activate the virtual environment:\n")
	sb.WriteString("   ```bash\n")
	sb.WriteString("   # On Unix/Linux/MacOS\n")
	sb.WriteString("   source .venv/bin/activate\n")
	sb.WriteString("   \n")
	sb.WriteString("   # On Windows\n")
	sb.WriteString("   .venv\\Scripts\\activate.bat\n")
	sb.WriteString("   ```\n\n")

	sb.WriteString("3. Install dependencies:\n")
	sb.WriteString("   ```bash\n")
	sb.WriteString("   pip install -e .\n")
	sb.WriteString("   ```\n\n")

	sb.WriteString("## Running the Server\n\n")
	sb.WriteString("You can run the server using the provided script:\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("python scripts/run.py\n")
	sb.WriteString("```\n\n")

	sb.WriteString("Or directly:\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("python src/mcp_server.py\n")
	sb.WriteString("```\n\n")

	sb.WriteString("## Configuration\n\n")
	sb.WriteString("Set the following environment variables to configure the server:\n\n")
	sb.WriteString("- `SERVICE_URL`: The base URL of the service to proxy (default: http://localhost:8080)\n")
	sb.WriteString("- `PORT`: The port to run the MCP server on (default: 8000)\n\n")

	sb.WriteString("## License\n\n")
	sb.WriteString("MIT\n")

	return os.WriteFile(filePath, []byte(sb.String()), 0644)
}

// processPathsIntoTools converts OpenAPI paths to MCP tools
func (g *Generator) processPathsIntoTools(doc *openapi3.T, s *server.MCPServer) error {
	g.document = doc

	for path, pathItem := range doc.Paths.Map() {
		// Process each HTTP method
		for method, opRef := range pathItem.Operations() {
			if opRef == nil {
				continue
			}

			op := opRef
			toolID := sanitizePathForToolID(path, method)
			toolDesc := op.Summary
			if toolDesc == "" {
				toolDesc = op.Description
			}

			// Create tool options
			toolOpts := []mcp.ToolOption{mcp.WithDescription(toolDesc)}

			// Process parameters into tool options
			for _, paramRef := range op.Parameters {
				if paramRef == nil || paramRef.Value == nil {
					continue
				}

				param := paramRef.Value
				if param.Schema == nil || param.Schema.Value == nil {
					continue
				}

				schema := param.Schema.Value
				propOpts := []mcp.PropertyOption{}

				if param.Required {
					propOpts = append(propOpts, mcp.Required())
				}

				if param.Description != "" {
					propOpts = append(propOpts, mcp.Description(param.Description))
				}

				switch schema.Type {
				case "string":
					// Add enum values if available
					if len(schema.Enum) > 0 {
						enumValues := make([]string, 0, len(schema.Enum))
						for _, v := range schema.Enum {
							if s, ok := v.(string); ok {
								enumValues = append(enumValues, s)
							}
						}
						if len(enumValues) > 0 {
							propOpts = append(propOpts, mcp.Enum(enumValues...))
						}
					}

					toolOpts = append(toolOpts, mcp.WithString(param.Name, propOpts...))
				case "integer", "number":
					toolOpts = append(toolOpts, mcp.WithNumber(param.Name, propOpts...))
				case "boolean":
					toolOpts = append(toolOpts, mcp.WithBoolean(param.Name, propOpts...))
				default:
					// Handle arrays and objects as strings for simplicity
					toolOpts = append(toolOpts, mcp.WithString(param.Name, propOpts...))
				}
			}

			// Process request body
			if op.RequestBody != nil && op.RequestBody.Value != nil {
				reqBody := op.RequestBody.Value

				for _, mediaType := range reqBody.Content {
					if mediaType.Schema != nil && mediaType.Schema.Value != nil {
						propOpts := []mcp.PropertyOption{}

						if reqBody.Required {
							propOpts = append(propOpts, mcp.Required())
						}

						desc := "Request body"
						if reqBody.Description != "" {
							desc = reqBody.Description
						}

						propOpts = append(propOpts, mcp.Description(desc))
						toolOpts = append(toolOpts, mcp.WithString("body", propOpts...))
						break
					}
				}
			}

			// Create the tool with all options
			tool := mcp.NewTool(toolID, toolOpts...)

			// Add tool to server with handler
			s.AddTool(tool, g.createToolHandler(op, path, method))

			g.logger.Debug("Added tool",
				zap.String("id", toolID),
				zap.String("path", path),
				zap.String("method", method))
		}
	}

	return nil
}

// createToolHandler returns a handler function for an MCP tool
func (g *Generator) createToolHandler(op *openapi3.Operation, path, method string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get the service URL from config
		serviceURL := config.GetString("service.url")
		if serviceURL == "" {
			// If no service URL is provided, return a mock response
			resultText := fmt.Sprintf("Mock response for %s %s\nParams: %v",
				method,
				path,
				request.Params.Arguments)
			return mcp.NewToolResultText(resultText), nil
		}

		// Create the full URL
		fullURL := buildURL(serviceURL, path, request.Params.Arguments, op.Parameters)

		// Create HTTP request
		httpReq, err := createHTTPRequest(ctx, method, fullURL, request.Params.Arguments, op)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Add authorization header if provided
		authHeader := config.GetString("service.authorization")
		if authHeader != "" {
			httpReq.Header.Set("Authorization", authHeader)
		}

		// Set common headers
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "application/json")

		// Create HTTP client with timeout
		timeout := config.GetDuration("client.timeout")
		if timeout == 0 {
			timeout = 30 * time.Second
		}
		client := &http.Client{
			Timeout: timeout,
		}

		// Execute the request
		g.logger.Debug("Executing API request",
			zap.String("method", method),
			zap.String("url", fullURL),
		)

		resp, err := client.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("API request failed: %w", err)
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		// Check if response is successful
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("API returned error status: %d - %s", resp.StatusCode, string(body))
		}

		// Return the response
		return mcp.NewToolResultText(string(body)), nil
	}
}

// buildURL constructs the full URL with path parameters and query parameters
func buildURL(baseURL, path string, args map[string]interface{}, parameters []*openapi3.ParameterRef) string {
	// Replace path parameters
	for _, paramRef := range parameters {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}

		param := paramRef.Value
		if param.In == "path" {
			if val, ok := args[param.Name]; ok {
				placeholder := fmt.Sprintf("{%s}", param.Name)
				path = strings.Replace(path, placeholder, fmt.Sprintf("%v", val), -1)
			}
		}
	}

	// Normalize base URL and path
	if !strings.HasSuffix(baseURL, "/") && !strings.HasPrefix(path, "/") {
		baseURL += "/"
	} else if strings.HasSuffix(baseURL, "/") && strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	// Create URL instance
	u, err := url.Parse(baseURL + path)
	if err != nil {
		return baseURL + path
	}

	// Add query parameters
	q := u.Query()
	for _, paramRef := range parameters {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}

		param := paramRef.Value
		if param.In == "query" {
			if val, ok := args[param.Name]; ok {
				q.Add(param.Name, fmt.Sprintf("%v", val))
			}
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// createHTTPRequest creates an HTTP request with the appropriate method and body
func createHTTPRequest(ctx context.Context, method, url string, args map[string]interface{}, op *openapi3.Operation) (*http.Request, error) {
	var body []byte
	var err error

	// Add request body for methods that support it
	if method == "POST" || method == "PUT" || method == "PATCH" {
		// Check if there's a body parameter in the arguments
		if bodyArg, ok := args["body"]; ok {
			// If body is a string, use it directly
			if bodyStr, ok := bodyArg.(string); ok {
				body = []byte(bodyStr)
			} else {
				// Otherwise, marshal it to JSON
				body, err = json.Marshal(bodyArg)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal request body: %w", err)
				}
			}
		} else {
			// If no body parameter is found, use all arguments that are not used in path or query
			bodyMap := make(map[string]interface{})
			for name, value := range args {
				isPathOrQuery := false
				for _, paramRef := range op.Parameters {
					if paramRef != nil && paramRef.Value != nil {
						param := paramRef.Value
						if (param.In == "path" || param.In == "query") && param.Name == name {
							isPathOrQuery = true
							break
						}
					}
				}
				if !isPathOrQuery {
					bodyMap[name] = value
				}
			}

			if len(bodyMap) > 0 {
				body, err = json.Marshal(bodyMap)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal request body: %w", err)
				}
			}
		}
	}

	// Create the request
	if body != nil {
		return http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	}
	return http.NewRequestWithContext(ctx, method, url, nil)
}

// generateServerCode writes the MCP server code to a file
func (g *Generator) generateServerCode(s *server.MCPServer, filePath string) error {
	var toolRegistrations strings.Builder

	// Get the OpenAPI document from the Generator context
	doc := g.document

	// Write Python imports and server setup
	fmt.Fprintf(&toolRegistrations, `#!/usr/bin/env python3
"""
MCP Server generated from OpenAPI specification.
"""
from mcp.server.fastmcp import FastMCP
import os
import requests
import logging
import json
from urllib.parse import urlencode
from typing import Dict, Any, Optional, Union

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Create MCP server
mcp = FastMCP("%s")

# Get service URL from environment
service_url = os.getenv("SERVICE_URL", "http://localhost:8080")
logger.info(f"Using service URL: {service_url}")

def build_url(base_url: str, path: str, params: Dict[str, Any] = None) -> str:
    """Build URL with path parameters and query parameters."""
    # Handle path parameters
    if params:
        for key, value in params.items():
            if "{" + key + "}" in path:
                path = path.replace("{" + key + "}", str(value))
    
    # Normalize URL joining
    if base_url.endswith("/") and path.startswith("/"):
        path = path[1:]
    elif not base_url.endswith("/") and not path.startswith("/"):
        base_url += "/"
    
    url = base_url + path
    
    # Add query parameters
    if params:
        query_params = {k: v for k, v in params.items() if "{" + k + "}" not in path}
        if query_params:
            url += "?" + urlencode(query_params)
    
    return url
`, doc.Info.Title)

	// Iterate over all paths in the OpenAPI document
	for path, pathItem := range doc.Paths.Map() {
		for method, op := range pathItem.Operations() {
			if op == nil {
				continue
			}

			toolID := sanitizePathForToolID(path, method)
			description := op.Summary
			if description == "" {
				description = op.Description
			}
			if description == "" {
				description = fmt.Sprintf("%s %s", method, path)
			}

			// Start building tool registration code
			fmt.Fprintf(&toolRegistrations, "\n@mcp.tool()\ndef %s(", toolID)

			// Add parameters
			var params []string
			var requiredParams []string
			var optionalParams []string

			for _, paramRef := range op.Parameters {
				if paramRef == nil || paramRef.Value == nil {
					continue
				}

				param := paramRef.Value
				paramName := sanitizeParamName(param.Name)
				paramType := "str" // Default to string type

				if param.Schema != nil && param.Schema.Value != nil {
					switch param.Schema.Value.Type {
					case "integer":
						paramType = "int"
					case "number":
						paramType = "float"
					case "boolean":
						paramType = "bool"
					}
				}

				if param.Required {
					requiredParams = append(requiredParams, fmt.Sprintf("%s: %s", paramName, paramType))
				} else {
					optionalParams = append(optionalParams, fmt.Sprintf("%s: Optional[%s] = None", paramName, paramType))
				}
			}

			// Add body parameter if needed
			if op.RequestBody != nil && op.RequestBody.Value != nil {
				if op.RequestBody.Value.Required {
					requiredParams = append(requiredParams, "body: Union[str, Dict[str, Any]]")
				} else {
					optionalParams = append(optionalParams, "body: Optional[Union[str, Dict[str, Any]]] = None")
				}
			}

			// Combine parameters with required ones first, then optional ones
			params = append(requiredParams, optionalParams...)

			fmt.Fprintf(&toolRegistrations, "%s) -> str:\n", strings.Join(params, ", "))
			fmt.Fprintf(&toolRegistrations, "    \"\"\"%s\"\"\"\n", description)

			// Build parameters dictionary
			fmt.Fprintf(&toolRegistrations, "    params: Dict[str, Any] = {}\n")
			for _, paramRef := range op.Parameters {
				if paramRef == nil || paramRef.Value == nil {
					continue
				}

				param := paramRef.Value
				paramName := sanitizeParamName(param.Name)
				fmt.Fprintf(&toolRegistrations, "    if %s is not None:\n", paramName)
				fmt.Fprintf(&toolRegistrations, "        params[\"%s\"] = %s\n", param.Name, paramName)
			}

			// Build URL
			fmt.Fprintf(&toolRegistrations, "    url = build_url(service_url, \"%s\", params)\n", path)
			fmt.Fprintf(&toolRegistrations, "    logger.info(f\"Making request to: {url}\")\n\n")

			// Add headers setup
			fmt.Fprintf(&toolRegistrations, "    headers = {\"Content-Type\": \"application/json\"}\n")
			for _, paramRef := range op.Parameters {
				if paramRef == nil || paramRef.Value == nil {
					continue
				}

				param := paramRef.Value
				if param.In == "header" {
					paramName := sanitizeParamName(param.Name)
					fmt.Fprintf(&toolRegistrations, "    if %s is not None:\n", paramName)
					fmt.Fprintf(&toolRegistrations, "        headers[\"%s\"] = str(%s)\n", param.Name, paramName)
				}
			}

			// Add request code
			fmt.Fprintf(&toolRegistrations, "\n    try:\n")
			if method == "GET" {
				fmt.Fprintf(&toolRegistrations, "        response = requests.get(url, headers=headers)\n")
			} else {
				if op.RequestBody != nil && op.RequestBody.Value != nil {
					fmt.Fprintf(&toolRegistrations, "        # Handle request body\n")
					fmt.Fprintf(&toolRegistrations, "        if isinstance(body, str):\n")
					fmt.Fprintf(&toolRegistrations, "            try:\n")
					fmt.Fprintf(&toolRegistrations, "                # Try to parse as JSON\n")
					fmt.Fprintf(&toolRegistrations, "                json_body = json.loads(body)\n")
					fmt.Fprintf(&toolRegistrations, "                response = requests.%s(url, headers=headers, json=json_body)\n", strings.ToLower(method))
					fmt.Fprintf(&toolRegistrations, "            except json.JSONDecodeError:\n")
					fmt.Fprintf(&toolRegistrations, "                # If not JSON, send as raw string\n")
					fmt.Fprintf(&toolRegistrations, "                response = requests.%s(url, headers=headers, data=body)\n", strings.ToLower(method))
					fmt.Fprintf(&toolRegistrations, "        else:\n")
					fmt.Fprintf(&toolRegistrations, "            response = requests.%s(url, headers=headers, json=body)\n", strings.ToLower(method))
				} else {
					fmt.Fprintf(&toolRegistrations, "        response = requests.%s(url, headers=headers)\n", strings.ToLower(method))
				}
			}
			fmt.Fprintf(&toolRegistrations, "        response.raise_for_status()\n")
			fmt.Fprintf(&toolRegistrations, "        return response.text\n")
			fmt.Fprintf(&toolRegistrations, "    except requests.RequestException as e:\n")
			fmt.Fprintf(&toolRegistrations, "        error_msg = str(e)\n")
			fmt.Fprintf(&toolRegistrations, "        if e.response is not None:\n")
			fmt.Fprintf(&toolRegistrations, "            error_msg = f\"{error_msg} - Response: {e.response.text}\"\n")
			fmt.Fprintf(&toolRegistrations, "        logger.error(f\"%s request failed: {error_msg}\")\n", toolID)
			fmt.Fprintf(&toolRegistrations, "        raise\n")
		}
	}

	// Add main block
	fmt.Fprintf(&toolRegistrations, "\nif __name__ == \"__main__\":\n")
	fmt.Fprintf(&toolRegistrations, "    # Get server port from environment or use default\n")
	fmt.Fprintf(&toolRegistrations, "    port = int(os.getenv(\"PORT\", \"8000\"))\n")
	fmt.Fprintf(&toolRegistrations, "    logger.info(f\"Starting MCP server on port {port}\")\n")
	fmt.Fprintf(&toolRegistrations, "    # Run the server\n")
	fmt.Fprintf(&toolRegistrations, "    mcp.run(port=port)\n")

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for server code: %w", err)
	}

	// Write the code to file
	return os.WriteFile(filePath, []byte(toolRegistrations.String()), 0755)
}

// sanitizePathForToolID converts an OpenAPI path to a valid tool ID
func sanitizePathForToolID(path, method string) string {
	// Replace path parameters with camelCase names
	sanitized := strings.ReplaceAll(path, "{", "")
	sanitized = strings.ReplaceAll(sanitized, "}", "")
	sanitized = strings.ReplaceAll(sanitized, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, "-", "_")

	// Remove leading underscore if present
	sanitized = strings.TrimPrefix(sanitized, "_")

	// Add method prefix
	return fmt.Sprintf("%s%s", strings.ToLower(method), strings.Title(sanitized))
}

// sanitizeParamName converts an OpenAPI parameter name to a valid Python variable name
func sanitizeParamName(name string) string {
	// Replace hyphens with underscores
	name = strings.ReplaceAll(name, "-", "_")
	// Replace any other invalid characters
	name = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return r
		}
		return '_'
	}, name)
	return name
}

// sanitizeForPackageName sanitizes a string to be used as a package name
func sanitizeForPackageName(name string) string {
	// Convert to lowercase and replace spaces with underscores
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")

	// Replace invalid characters with underscores
	name = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return r
		}
		return '_'
	}, name)

	// Ensure it starts with a letter
	if len(name) > 0 && !unicode.IsLetter(rune(name[0])) {
		name = "mcp_" + name
	}

	return name
}

// generateRequirements writes the Python package requirements to a file
func (g *Generator) generateRequirements(filePath string) error {
	requirements := `mcp-sdk>=0.1.0
requests>=2.28.0
`
	return os.WriteFile(filePath, []byte(requirements), 0644)
}
