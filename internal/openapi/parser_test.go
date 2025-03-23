package openapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestFetchAndParse(t *testing.T) {
	// Setup a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"openapi": "3.0.0",
			"info": {
				"title": "Test API",
				"version": "1.0.0",
				"description": "A test API"
			},
			"paths": {
				"/users": {
					"get": {
						"summary": "Get all users",
						"responses": {
							"200": {
								"description": "Successful response",
								"content": {
									"application/json": {
										"schema": {
											"type": "array",
											"items": {
												"$ref": "#/components/schemas/User"
											}
										}
									}
								}
							}
						}
					}
				}
			},
			"components": {
				"schemas": {
					"User": {
						"type": "object",
						"properties": {
							"id": {
								"type": "integer"
							},
							"name": {
								"type": "string"
							}
						}
					}
				}
			}
		}`))
	}))
	defer server.Close()

	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Create parser
	parser := NewParser(logger)

	// Test valid URL
	ctx := context.Background()
	doc, err := parser.FetchAndParse(ctx, server.URL)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Verify parsed content
	if doc.Info.Title != "Test API" {
		t.Errorf("Expected title 'Test API', got '%s'", doc.Info.Title)
	}

	pathCount := 0
	if doc.Paths != nil {
		pathCount = len(doc.Paths.Map())
	}
	if pathCount != 1 {
		t.Errorf("Expected 1 path, got %d", pathCount)
	}

	schemaCount := 0
	if doc.Components.Schemas != nil {
		schemaCount = len(doc.Components.Schemas)
	}
	if schemaCount != 1 {
		t.Errorf("Expected 1 schema, got %d", schemaCount)
	}

	// Test invalid URL
	_, err = parser.FetchAndParse(ctx, "invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL but got none")
	}
}
