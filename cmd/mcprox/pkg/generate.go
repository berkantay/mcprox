package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/berkantay/mcprox/internal/mcp"
	"github.com/berkantay/mcprox/internal/openapi"
	"github.com/spf13/cobra"
)

var (
	swaggerURL string
	timeout    int
	outputDir  string
)

func init() {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate MCP server from OpenAPI documentation",
		Long: `Fetches OpenAPI/Swagger documentation from a URL and generates a fully functional 
Model Context Protocol (MCP) server.

Example:
  godoc-mcp generate --url http://localhost:8080/swagger/doc.json`,
		RunE: generateMCP,
	}

	generateCmd.Flags().StringVarP(&swaggerURL, "url", "u", "", "URL to fetch OpenAPI documentation (required)")
	generateCmd.MarkFlagRequired("url")
	generateCmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "Timeout in seconds for HTTP requests")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for generated server (default is ./generated)")

	rootCmd.AddCommand(generateCmd)
}

func generateMCP(cmd *cobra.Command, args []string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// Create OpenAPI parser
	parser := openapi.NewParser(logger)

	// Fetch and parse OpenAPI documentation
	doc, err := parser.FetchAndParse(ctx, swaggerURL)
	if err != nil {
		return fmt.Errorf("failed to fetch and parse OpenAPI documentation: %w", err)
	}

	// Create MCP generator
	generator := mcp.NewGenerator(logger, outputDir)

	// Generate MCP server
	if err := generator.Generate(ctx, doc); err != nil {
		return fmt.Errorf("failed to generate MCP server: %w", err)
	}

	logger.Info("MCP server generation completed successfully")
	return nil
}
