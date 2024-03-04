package config

import (
	"os"
	"strconv"
)

func parseEnvVars(c *settings, configFile *string) {
	if servAddr := os.Getenv("SERVER_ADDRESS"); servAddr != "" {
		c.ServAddr = servAddr
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		c.BaseURL = baseURL
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		c.LogLevel = logLevel
	}

	if filePath := os.Getenv("FILE_STORAGE_PATH"); filePath != "" {
		c.FilePath = filePath
	}

	if dbAddr := os.Getenv("DATABASE_DSN"); dbAddr != "" {
		c.DBAddr = dbAddr
	}

	if trustedSubnet := os.Getenv("TRUSTED_SUBNET"); trustedSubnet != "" {
		c.TrustedSubnet = trustedSubnet
	}

	if tlsEnabled := os.Getenv("ENABLE_HTTPS"); tlsEnabled != "" {
		if isEnabled, err := strconv.ParseBool(tlsEnabled); err == nil {
			c.IsTLSEnabled = &isEnabled
		}
	}

	if config := os.Getenv("CONFIG"); config != "" {
		*configFile = config
	}
}
