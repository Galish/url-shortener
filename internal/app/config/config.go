// Package config represents the application configuration.
package config

import (
	"fmt"
)

// Config stores the application configuration.
type Config struct {
	ServAddr     string
	BaseURL      string
	LogLevel     string
	FilePath     string
	DBAddr       string
	IsTLSEnabled bool
}

type option func(*Config)

// New generates a configuration by parsing flags and environment variables.
func New() *Config {
	var flags = new(settings)
	var envVars = new(settings)
	var file = new(settings)

	parseFlags(flags)
	parseEnvVars(envVars)
	if err := parseFile(
		file,
		flags.fileConfigPath,
		envVars.fileConfigPath,
	); err != nil {
		fmt.Println(fmt.Errorf("unable to read config file: %s", err))
	}

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
