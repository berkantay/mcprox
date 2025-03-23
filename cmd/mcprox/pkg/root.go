package pkg

import (
	"fmt"
	"os"

	"github.com/berkantay/mcprox/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgFile string
	debug   bool
	logger  *zap.Logger
	rootCmd = &cobra.Command{
		Use:   "mcprox",
		Short: "Generate MCP proxy from OpenAPI documentation",
		Long: `A robust tool that retrieves and parses OpenAPI/Swagger documentation from a URL and
generates a fully functional Model Context Protocol (MCP) proxy using the mark3labs/mcp-go library.`,
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mcprox.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")

	// Add service configuration flags
	rootCmd.PersistentFlags().String("service-url", "", "base URL of the target API service")
	rootCmd.PersistentFlags().String("service-auth", "", "authorization header value for the target API")

	// Bind flags to viper
	viper.BindPFlag("service.url", rootCmd.PersistentFlags().Lookup("service-url"))
	viper.BindPFlag("service.authorization", rootCmd.PersistentFlags().Lookup("service-auth"))
}

func initConfig() {
	config.Init(cfgFile)

	// Override config with command line flags
	if debug {
		config.SetBool("debug", true)
	}
}

func initLogger() {
	var err error
	if config.GetBool("debug") {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
}
