package config

import (
	"os"
	"strconv"
)

func parseEnvVars(s *settings) {
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		s.ServAddr = envAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		s.BaseURL = envBaseURL
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		s.logLevel = logLevel
	}

	if filePath := os.Getenv("FILE_STORAGE_PATH"); filePath != "" {
		s.FilePath = filePath
	}

	if dbAddr := os.Getenv("DATABASE_DSN"); dbAddr != "" {
		s.DBAddr = dbAddr
	}

	if tlsEnabled := os.Getenv("ENABLE_HTTPS"); tlsEnabled == "" {
		if isEnabled, _ := strconv.ParseBool(tlsEnabled); isEnabled {
			s.IsTLSEnabled = true
		}
	}

	if configPath := os.Getenv("CONFIG"); configPath != "" {
		s.fileConfigPath = configPath
	}
}
