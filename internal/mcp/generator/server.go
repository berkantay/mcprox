package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// generateServerCode writes the MCP server code to a file
func (g *Generator) generateServerCode(filePath string) error {
	// Get the OpenAPI document from the Generator context
	doc := g.document

	// Create a new ToolBuilder to handle code generation
	tb := NewToolBuilder()

	// Write Python imports
	tb.WriteImports()

	// Write logger setup
	tb.WriteSetupLogger()

	// Create MCP server
	tb.WriteCreateMCPServer(doc.Info.Title)

	// Get service URL from environment
	tb.WriteGetServiceURL()

	// Write function to build URL with path parameters and query parameters
	tb.WriteBuildURL()

	// Iterate over all paths in the OpenAPI document
	for path, pathItem := range doc.Paths.Map() {
		for method, op := range pathItem.Operations() {
			if op == nil {
				continue
			}

			// Generate the tool definition code
			tb.WriteToolDefinition(path, method, op)
		}
	}

	// Add main block
	tb.WriteMainBlock()

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for server code: %w", err)
	}

	// Write the code to file
	return os.WriteFile(filePath, []byte(tb.String()), 0755)
}
