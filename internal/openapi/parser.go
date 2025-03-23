package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/berkantay/mcprox/internal/config"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

// Parser handles fetching and parsing OpenAPI documentation
type Parser struct {
	logger        *zap.Logger
	clientTimeout time.Duration
}

// NewParser creates a new OpenAPI parser
func NewParser(logger *zap.Logger) *Parser {
	timeout := time.Duration(config.GetInt("client.timeout")) * time.Second
	return &Parser{
		logger:        logger,
		clientTimeout: timeout,
	}
}

// FetchAndParse retrieves OpenAPI documentation from a URL and parses it
func (p *Parser) FetchAndParse(ctx context.Context, swaggerURL string) (*openapi3.T, error) {
	p.logger.Info("Fetching OpenAPI documentation", zap.String("url", swaggerURL))

	// Validate URL
	_, err := url.ParseRequestURI(swaggerURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: p.clientTimeout,
	}

	// Make HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, swaggerURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OpenAPI documentation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Pre-process body for OpenAPI 3.1.0 compatibility
	body, err = preprocessOpenAPISpec(body, p.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to preprocess OpenAPI spec: %w", err)
	}

	// Parse OpenAPI document
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI documentation: %w", err)
	}

	// Validate the document
	err = doc.Validate(ctx)
	if err != nil {
		return nil, fmt.Errorf("OpenAPI documentation validation failed: %w", err)
	}

	// Count paths and schemas
	pathCount := 0
	if doc.Paths != nil {
		pathCount = len(doc.Paths.Map())
	}

	schemaCount := 0
	if doc.Components.Schemas != nil {
		schemaCount = len(doc.Components.Schemas)
	}

	p.logger.Info("Successfully parsed OpenAPI documentation",
		zap.Int("paths", pathCount),
		zap.Int("components", schemaCount))

	return doc, nil
}

// preprocessOpenAPISpec adapts OpenAPI 3.1.0 to be compatible with OpenAPI 3.0.x
func preprocessOpenAPISpec(data []byte, logger *zap.Logger) ([]byte, error) {
	// Parse the JSON into a generic map
	var spec map[string]interface{}
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("error unmarshaling OpenAPI spec: %w", err)
	}

	// Check OpenAPI version
	if version, ok := spec["openapi"].(string); ok {
		if strings.HasPrefix(version, "3.1") {
			logger.Info("Converting OpenAPI 3.1.x spec to 3.0.0 for compatibility")
			// Change version to 3.0.0
			spec["openapi"] = "3.0.0"
		}
	}

	// Process components.schemas to handle null types
	if componentsInterface, exists := spec["components"]; exists {
		if components, ok := componentsInterface.(map[string]interface{}); ok {
			if schemasInterface, exists := components["schemas"]; exists {
				if schemas, ok := schemasInterface.(map[string]interface{}); ok {
					// Process each schema
					for _, schemaValue := range schemas {
						if schema, ok := schemaValue.(map[string]interface{}); ok {
							fixNullTypes(schema, logger)
							// Remove non-standard fields
							removeNonStandardFields(schema, logger)
						}
					}
				}
			}
		}
	}

	// Fix anyOf in parameters
	if pathsInterface, exists := spec["paths"]; exists {
		if paths, ok := pathsInterface.(map[string]interface{}); ok {
			for _, pathValue := range paths {
				if path, ok := pathValue.(map[string]interface{}); ok {
					for _, methodValue := range path {
						if method, ok := methodValue.(map[string]interface{}); ok {
							if paramsInterface, exists := method["parameters"]; exists {
								if params, ok := paramsInterface.([]interface{}); ok {
									for _, paramInterface := range params {
										if param, ok := paramInterface.(map[string]interface{}); ok {
											if schemaInterface, exists := param["schema"]; exists {
												if schema, ok := schemaInterface.(map[string]interface{}); ok {
													fixAnyOf(schema, logger)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Re-marshal the modified spec
	return json.Marshal(spec)
}

// removeNonStandardFields removes fields that are not standard in OpenAPI 3.0
func removeNonStandardFields(schema map[string]interface{}, logger *zap.Logger) {
	// List of non-standard fields to remove
	nonStandardFields := []string{
		"error_messages",
		"hide_error_details",
	}

	// Remove non-standard fields at top level
	for _, field := range nonStandardFields {
		if _, exists := schema[field]; exists {
			delete(schema, field)
			logger.Debug("Removed non-standard field from schema", zap.String("field", field))
		}
	}

	// Recursively process properties
	if properties, exists := schema["properties"].(map[string]interface{}); exists {
		for _, propValue := range properties {
			if propObj, ok := propValue.(map[string]interface{}); ok {
				removeNonStandardFields(propObj, logger)
			}
		}
	}

	// Process items in arrays
	if items, exists := schema["items"].(map[string]interface{}); exists {
		removeNonStandardFields(items, logger)
	}
}

// fixNullTypes recursively finds and fixes null types in schemas
func fixNullTypes(schema map[string]interface{}, logger *zap.Logger) {
	if properties, exists := schema["properties"].(map[string]interface{}); exists {
		for _, propValue := range properties {
			if propObj, ok := propValue.(map[string]interface{}); ok {
				// Check for anyOf that includes null
				if _, exists := propObj["anyOf"]; exists {
					fixAnyOf(propObj, logger)
				}

				// Recursively process nested properties
				fixNullTypes(propObj, logger)
			}
		}
	}
}

// fixAnyOf handles anyOf with null type by converting to a standard type with nullable
func fixAnyOf(schema map[string]interface{}, logger *zap.Logger) {
	if anyOfInterface, exists := schema["anyOf"]; exists {
		if anyOf, ok := anyOfInterface.([]interface{}); ok {
			// Look for [{"type": "something"}, {"type": "null"}] pattern
			if len(anyOf) == 2 {
				var mainType map[string]interface{}
				hasNull := false

				for _, typeObj := range anyOf {
					if obj, ok := typeObj.(map[string]interface{}); ok {
						if typeVal, exists := obj["type"]; exists {
							if typeVal == "null" {
								hasNull = true
							} else if _, ok := typeVal.(string); ok {
								mainType = obj
							}
						}
					}
				}

				if hasNull && mainType != nil {
					// Convert to standard type with nullable
					for k, v := range mainType {
						schema[k] = v
					}
					schema["nullable"] = true
					delete(schema, "anyOf")
					logger.Debug("Converted anyOf with null to nullable type")
				}
			}
		}
	}
}
