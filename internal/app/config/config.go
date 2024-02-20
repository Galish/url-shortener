// Package config represents the application configuration.
package config

// Config stores the application configuration.
type Config struct {
	ServAddr     string
	BaseURL      string
	LogLevel     string
	FilePath     string
	DBAddr       string
	IsTLSEnabled bool
}

var cfg Config

// New generates a configuration by parsing flags and environment variables.
func New() *Config {
	parseFlags()
	parseEnvVars()

	return &cfg
}
