package config

import "os"

func parseEnvVars() {
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		cfg.Addr = envAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseURL = envBaseURL
	}
}
