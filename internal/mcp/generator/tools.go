package generator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/berkantay/mcprox/internal/config"
	"github.com/berkantay/mcprox/internal/mcp/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

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
			toolID := utils.SanitizePathForToolID(path, method)
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
