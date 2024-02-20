package config

import (
	"os"
	"strconv"
)

func parseEnvVars() {
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		cfg.ServAddr = envAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseURL = envBaseURL
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	if filePath := os.Getenv("FILE_STORAGE_PATH"); filePath != "" {
		cfg.FilePath = filePath
	}

	if dbAddr := os.Getenv("DATABASE_DSN"); dbAddr != "" {
		cfg.DBAddr = dbAddr
	}

	if tlsEnabled, ok := os.LookupEnv("ENABLE_HTTPS"); ok {
		isEnabled, _ := strconv.ParseBool(tlsEnabled)
		cfg.IsTLSEnabled = isEnabled
	}
}
