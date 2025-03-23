package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/berkantay/mcprox/internal/config"
	"github.com/berkantay/mcprox/internal/mcp/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// Generator handles the creation of MCP server from OpenAPI specs
type Generator struct {
	logger    *zap.Logger
	outputDir string
	document  *openapi3.T
}

// New creates a new MCP generator
func New(logger *zap.Logger, outputDir ...string) *Generator {
	// Use provided output directory if specified, otherwise use default from config
	dir := config.GetString("output.dir")
	if len(outputDir) > 0 && outputDir[0] != "" {
		dir = outputDir[0]
	}

	return &Generator{
		logger:    logger,
		outputDir: dir,
	}
}

// Generate generates an MCP server from an OpenAPI spec
func (g *Generator) Generate(ctx context.Context, doc *openapi3.T) error {
	g.logger.Info("Generating MCP server from OpenAPI documentation")

	// Store the document in the generator
	g.document = doc

	folderName := strings.ToLower(strings.ReplaceAll(doc.Info.Title, " ", "_")) + "_mcp_server"

	// Set up project directory
	projectDir := filepath.Join(g.outputDir, folderName)
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
	if err := g.generateServerCode(serverPath); err != nil {
		return fmt.Errorf("failed to generate server code: %w", err)
	}

	// Generate project files
	if err := g.generateProjectFiles(doc); err != nil {
		return fmt.Errorf("failed to generate project files: %w", err)
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

// generateProjectFiles generates all required project files
func (g *Generator) generateProjectFiles(doc *openapi3.T) error {
	// Generate requirements.txt
	// requirementsPath := filepath.Join(g.outputDir, "requirements.txt")
	// if err := utils.GenerateRequirements(requirementsPath); err != nil {
	// 	return fmt.Errorf("failed to generate requirements.txt: %w", err)
	// }

	// Generate pyproject.toml
	pyprojectPath := filepath.Join(g.outputDir, "pyproject.toml")
	if err := utils.GeneratePyprojectToml(pyprojectPath, doc); err != nil {
		return fmt.Errorf("failed to generate pyproject.toml: %w", err)
	}

	// Generate .gitignore
	gitignorePath := filepath.Join(g.outputDir, ".gitignore")
	if err := utils.GenerateGitignore(gitignorePath); err != nil {
		return fmt.Errorf("failed to generate .gitignore: %w", err)
	}

	// Generate README.md
	readmePath := filepath.Join(g.outputDir, "README.md")
	if err := utils.GenerateReadme(readmePath, doc); err != nil {
		return fmt.Errorf("failed to generate README.md: %w", err)
	}

	// Generate setup scripts
	if err := utils.GenerateSetupScripts(g.outputDir); err != nil {
		return fmt.Errorf("failed to generate setup scripts: %w", err)
	}

	// Generate __init__.py files for package structure
	if err := utils.GenerateInitFiles(g.outputDir); err != nil {
		return fmt.Errorf("failed to generate __init__.py files: %w", err)
	}

	return nil
}
