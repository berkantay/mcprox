package pkg

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	helpCmd := &cobra.Command{
		Use:   "instructions",
		Short: "Display detailed instructions for mcprox",
		Long:  `Displays detailed information about how to use the mcprox tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			showHelp()
		},
	}

	rootCmd.AddCommand(helpCmd)
}

func showHelp() {
	fmt.Println("MCProx - OpenAPI to MCP Bridge Generator")
	fmt.Println("=========================================")

	fmt.Println("DESCRIPTION:")
	fmt.Println("    MCProx is a tool that analyzes OpenAPI/Swagger documentation and")
	fmt.Println("    generates a fully functional Model Context Protocol (MCP) proxy.")
	fmt.Println("    It uses the high-performance Fiber web framework as the HTTP server.")
	fmt.Println("    The generated proxy acts as a bridge between LLMs and your API.")

	fmt.Println("USAGE:")
	fmt.Println("    mcprox generate --url <swagger-url> [options]")

	fmt.Println("EXAMPLES:")
	fmt.Println("    # Generate MCP proxy from local Swagger")
	fmt.Println("    mcprox generate --url http://localhost:8080/swagger/doc.json")

	fmt.Println("    # Generate with service URL to call the actual API")
	fmt.Println("    mcprox generate --url https://api.example.com/swagger --service-url https://api.example.com")

	fmt.Println("    # Generate with increased timeout")
	fmt.Println("    mcprox generate --url https://api.example.com/swagger --timeout 60")

	fmt.Println("    # Use a custom configuration file")
	fmt.Println("    mcprox --config /path/to/config.yaml generate --url http://localhost:8080/swagger/doc.json")

	fmt.Println("CONFIGURATION:")
	fmt.Println("    Configuration can be specified through a YAML file (~/.mcprox.yaml by default)")
	fmt.Println("    with the following structure:")
	fmt.Println("    ```yaml")
	fmt.Println("    debug: false")
	fmt.Println("    client:")
	fmt.Println("      timeout: 30")
	fmt.Println("    server:")
	fmt.Println("      port: 8080")
	fmt.Println("    output:")
	fmt.Println("      dir: ./generated")
	fmt.Println("    service:")
	fmt.Println("      url: https://api.example.com")
	fmt.Println("      authorization: Bearer your-token")
	fmt.Println("    ```")

	fmt.Println("SERVER DETAILS:")
	fmt.Println("    The generated MCP proxy uses Fiber (github.com/gofiber/fiber) for high performance.")
	fmt.Println("    Endpoints:")
	fmt.Println("    - POST /api/mcp   : MCP protocol endpoint")
	fmt.Println("    - GET  /health    : Health check endpoint")

	fmt.Println("REPOSITORY:")
	fmt.Println("    https://github.com/berkantay/mcprox")
}
