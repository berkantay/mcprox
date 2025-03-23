// This file provides backward compatibility for the generator package refactoring

package mcp

import (
	"context"

	"github.com/berkantay/mcprox/internal/mcp/generator"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

// NewGenerator creates a new MCP generator
// This maintains backward compatibility with existing code
func NewGenerator(logger *zap.Logger, outputDir ...string) *Generator {
	return &Generator{
		gen: generator.New(logger, outputDir...),
	}
}

// Generator handles the creation of MCP server from OpenAPI specs
// This is a facade that delegates to the new generator package
type Generator struct {
	gen *generator.Generator
}

// Generate generates an MCP server from an OpenAPI spec
func (g *Generator) Generate(ctx context.Context, doc *openapi3.T) error {
	return g.gen.Generate(ctx, doc)
}
