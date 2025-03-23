package generator

import (
	"fmt"
	"strings"

	"github.com/berkantay/mcprox/internal/mcp/utils"
	"github.com/getkin/kin-openapi/openapi3"
)

// ToolBuilder handles the generation of Python code for MCP tools
type ToolBuilder struct {
	builder strings.Builder
}

// NewToolBuilder creates a new ToolBuilder instance
func NewToolBuilder() *ToolBuilder {
	return &ToolBuilder{
		builder: strings.Builder{},
	}
}

// String returns the built string
func (tb *ToolBuilder) String() string {
	return tb.builder.String()
}

// WriteImports writes the Python imports
func (tb *ToolBuilder) WriteImports() {
	fmt.Fprintf(&tb.builder, `
#!/usr/bin/env python3
"""
MCP Server generated from OpenAPI specification.
"""
import os
import httpx
import logging
import json
from urllib.parse import urlencode
from typing import Dict, Any, Optional, Union

# Import MCP framework
from mcp.server.fastmcp import FastMCP
`)
}

// WriteSetupLogger writes the logger setup code
func (tb *ToolBuilder) WriteSetupLogger() {
	fmt.Fprintf(&tb.builder, `
# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)
`)
}

// WriteCreateMCPServer writes the code to create an MCP server
func (tb *ToolBuilder) WriteCreateMCPServer(serverName string) {
	fmt.Fprintf(&tb.builder, `
# Create MCP server
mcp = FastMCP("%s", description="MCP Server for %s API")
`, serverName, serverName)
}

// WriteGetServiceURL writes the code to get the service URL from environment
func (tb *ToolBuilder) WriteGetServiceURL() {
	fmt.Fprintf(&tb.builder, `
# Get service URL from environment
service_url = os.getenv("SERVICE_URL", "http://localhost:8080")
logger.info(f"Using service URL: {service_url}")
`)
}

// WriteBuildURL writes the function to build URLs
func (tb *ToolBuilder) WriteBuildURL() {
	fmt.Fprintf(&tb.builder, `
def build_url(base_url: str, path: str, params: Dict[str, Any] = None) -> str:
    """Build URL with path parameters and query parameters."""
    # Handle path parameters
    url = base_url
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

    # Return the URL
    return url
`)
}

// WriteToolDefinition writes the code for a tool definition
func (tb *ToolBuilder) WriteToolDefinition(path, method string, op *openapi3.Operation) {
	toolID := utils.SanitizePathForToolID(path, method)
	description := op.Summary
	if description == "" {
		description = op.Description
	}
	if description == "" {
		description = fmt.Sprintf("%s %s", method, path)
	}

	// Start building tool registration code
	fmt.Fprintf(&tb.builder, "\n@mcp.tool()\ndef %s(", toolID)

	// Add parameters
	var params []string
	var requiredParams []string
	var optionalParams []string

	tb.buildParameterLists(op, &requiredParams, &optionalParams)

	// Combine parameters with required ones first, then optional ones
	params = append(requiredParams, optionalParams...)

	fmt.Fprintf(&tb.builder, "%s) -> str:\n", strings.Join(params, ", "))
	fmt.Fprintf(&tb.builder, "    \"\"\"%s\"\"\"\n", description)

	tb.writeParametersDictionary(op)
	tb.writeBuildURLCall(path)
	tb.writeHeadersSetup(op)
	tb.writeRequestCode(method, op)
}

// buildParameterLists builds the lists of required and optional parameters
func (tb *ToolBuilder) buildParameterLists(op *openapi3.Operation, requiredParams, optionalParams *[]string) {
	// Process path/query parameters
	for _, paramRef := range op.Parameters {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}

		param := paramRef.Value
		paramName := utils.SanitizeParamName(param.Name)
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
			*requiredParams = append(*requiredParams, fmt.Sprintf("%s: %s", paramName, paramType))
		} else {
			*optionalParams = append(*optionalParams, fmt.Sprintf("%s: Optional[%s] = None", paramName, paramType))
		}
	}

	// Add body parameter if needed
	if op.RequestBody != nil && op.RequestBody.Value != nil {
		if op.RequestBody.Value.Required {
			*requiredParams = append(*requiredParams, "body: Union[str, Dict[str, Any]]")
		} else {
			*optionalParams = append(*optionalParams, "body: Optional[Union[str, Dict[str, Any]]] = None")
		}
	}
}

// writeParametersDictionary writes the code to build the parameters dictionary
func (tb *ToolBuilder) writeParametersDictionary(op *openapi3.Operation) {
	fmt.Fprintf(&tb.builder, "    params: Dict[str, Any] = {}\n")
	for _, paramRef := range op.Parameters {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}

		param := paramRef.Value
		paramName := utils.SanitizeParamName(param.Name)
		fmt.Fprintf(&tb.builder, "    if %s is not None:\n", paramName)
		fmt.Fprintf(&tb.builder, "        params[\"%s\"] = %s\n", param.Name, paramName)
	}
}

// writeBuildURLCall writes the code to build the URL
func (tb *ToolBuilder) writeBuildURLCall(path string) {
	fmt.Fprintf(&tb.builder, "    url = build_url(service_url, \"%s\", params)\n", path)
	fmt.Fprintf(&tb.builder, "    logger.info(f\"Making request to: {url}\")\n\n")
}

// writeHeadersSetup writes the code to set up headers
func (tb *ToolBuilder) writeHeadersSetup(op *openapi3.Operation) {
	fmt.Fprintf(&tb.builder, "    headers = {\"Content-Type\": \"application/json\"}\n")
	for _, paramRef := range op.Parameters {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}

		param := paramRef.Value
		if param.In == "header" {
			paramName := utils.SanitizeParamName(param.Name)
			fmt.Fprintf(&tb.builder, "    if %s is not None:\n", paramName)
			fmt.Fprintf(&tb.builder, "        headers[\"%s\"] = str(%s)\n", param.Name, paramName)
		}
	}
}

// writeRequestCode writes the code to make the HTTP request
func (tb *ToolBuilder) writeRequestCode(method string, op *openapi3.Operation) {
	toolID := utils.SanitizePathForToolID("", method) // Only need method for error message

	fmt.Fprintf(&tb.builder, "\n    try:\n")
	if method == "GET" {
		fmt.Fprintf(&tb.builder, "        response = httpx.get(url, headers=headers)\n")
	} else {
		if op.RequestBody != nil && op.RequestBody.Value != nil {
			fmt.Fprintf(&tb.builder, "        # Handle request body\n")
			fmt.Fprintf(&tb.builder, "        if isinstance(body, str):\n")
			fmt.Fprintf(&tb.builder, "            try:\n")
			fmt.Fprintf(&tb.builder, "                # Try to parse as JSON\n")
			fmt.Fprintf(&tb.builder, "                json_body = json.loads(body)\n")
			fmt.Fprintf(&tb.builder, "                response = httpx.%s(url, headers=headers, json=json_body)\n", strings.ToLower(method))
			fmt.Fprintf(&tb.builder, "            except json.JSONDecodeError:\n")
			fmt.Fprintf(&tb.builder, "                # If not JSON, send as raw string\n")
			fmt.Fprintf(&tb.builder, "                response = httpx.%s(url, headers=headers, content=body)\n", strings.ToLower(method))
			fmt.Fprintf(&tb.builder, "        else:\n")
			fmt.Fprintf(&tb.builder, "            response = httpx.%s(url, headers=headers, json=body)\n", strings.ToLower(method))
		} else {
			fmt.Fprintf(&tb.builder, "        response = httpx.%s(url, headers=headers)\n", strings.ToLower(method))
		}
	}
	fmt.Fprintf(&tb.builder, "        response.raise_for_status()\n")
	fmt.Fprintf(&tb.builder, "        return response.text\n")
	fmt.Fprintf(&tb.builder, "    except httpx.RequestError as e:\n")
	fmt.Fprintf(&tb.builder, "        error_msg = str(e)\n")
	fmt.Fprintf(&tb.builder, "        logger.error(f\"%s request failed: {error_msg}\")\n", toolID)
	fmt.Fprintf(&tb.builder, "        raise\n")
	fmt.Fprintf(&tb.builder, "    except httpx.HTTPStatusError as e:\n")
	fmt.Fprintf(&tb.builder, "        error_msg = str(e)\n")
	fmt.Fprintf(&tb.builder, "        if e.response is not None:\n")
	fmt.Fprintf(&tb.builder, "            error_msg = f\"{error_msg} - Response: {e.response.text}\"\n")
	fmt.Fprintf(&tb.builder, "        logger.error(f\"%s request failed: {error_msg}\")\n", toolID)
	fmt.Fprintf(&tb.builder, "        raise\n")
}

// WriteMainBlock writes the code for the main block to run the server
func (tb *ToolBuilder) WriteMainBlock() {
	fmt.Fprintf(&tb.builder, "\nif __name__ == \"__main__\":\n")
	fmt.Fprintf(&tb.builder, "    # Get server port from environment or use default\n")
	fmt.Fprintf(&tb.builder, "    port = int(os.getenv(\"PORT\", \"8000\"))\n")
	fmt.Fprintf(&tb.builder, "    logger.info(f\"Starting MCP server on port {port}\")\n")
	fmt.Fprintf(&tb.builder, "    # Run the server\n")
	fmt.Fprintf(&tb.builder, "    mcp.run(port=port)\n")
}
