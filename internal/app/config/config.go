// Package config represents the application configuration.
package config

import "net"

// Config stores the application configuration.
type Config struct {
	BaseURL       string
	DBAddr        string
	FilePath      string
	IsTLSEnabled  bool
	LogLevel      string
	ServAddr      string
	TrustedSubnet *net.IPNet
}

type option func(*Config)

// New generates a configuration by parsing flags and environment variables.
func New() *Config {
	var flags = new(settings)
	var envVars = new(settings)
	var file = new(settings)
	var configFile string

	parseFlags(flags, &configFile)
	parseEnvVars(envVars, &configFile)
	parseFile(configFile, file)

	return newConfig(
		withSettings(defaultSettings),
		withSettings(file),
		withSettings(flags),
		withSettings(envVars),
	)
}

func newConfig(opt ...option) *Config {
	var cfg = new(Config)

	for _, o := range opt {
		o(cfg)
	}

	return cfg
}
