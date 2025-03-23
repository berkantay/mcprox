package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Default configuration values
const (
	DefaultPort    = 8080
	DefaultTimeout = 30
)

// Init initializes the configuration
func Init(cfgFile string) {
	// Use config file from the flag if provided
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mcprox" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mcprox")
	}

	// Set default values
	SetDefaults()

	// Environment variables override config file
	viper.AutomaticEnv()

	// Read in config file
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// SetDefaults sets the default configuration values
func SetDefaults() {
	viper.SetDefault("server.port", DefaultPort)
	viper.SetDefault("client.timeout", DefaultTimeout)
	viper.SetDefault("debug", false)
	viper.SetDefault("output.dir", filepath.Join(".", "generated"))
	viper.SetDefault("service.url", "")
	viper.SetDefault("service.authorization", "")
}

// GetString retrieves a string configuration value
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt retrieves an integer configuration value
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool retrieves a boolean configuration value
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetStringMap retrieves a map configuration value
func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

// SetBool sets a boolean configuration value
func SetBool(key string, value bool) {
	viper.Set(key, value)
}

// GetDuration gets a duration value from the configuration
func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
