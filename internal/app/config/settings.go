package config

import (
	"net"

	"github.com/Galish/url-shortener/pkg/logger"
)

type settings struct {
	BaseURL       string `json:"base_url"`
	DBAddr        string `json:"database_dsn"`
	FilePath      string `json:"file_storage_path"`
	IsTLSEnabled  *bool  `json:"enable_https"`
	LogLevel      string `json:"log_level"`
	ServAddr      string `json:"server_address"`
	TrustedSubnet string `json:"trusted_subnet"`
}

var defaultSettings = &settings{
	BaseURL:  "http://localhost:8080",
	FilePath: "/tmp/short-url-db.json",
	LogLevel: "info",
	ServAddr: ":8080",
}

func withSettings(c *settings) func(*Config) {
	return func(cfg *Config) {
		if c.ServAddr != "" {
			cfg.ServAddr = c.ServAddr
		}

		if c.BaseURL != "" {
			cfg.BaseURL = c.BaseURL
		}

		if c.LogLevel != "" {
			cfg.LogLevel = c.LogLevel
		}

		if c.FilePath != "" {
			cfg.FilePath = c.FilePath
		}

		if c.DBAddr != "" {
			cfg.DBAddr = c.DBAddr
		}

		if c.TrustedSubnet != "" {
			_, ipv4Net, err := net.ParseCIDR(c.TrustedSubnet)
			if err != nil {
				logger.WithError(err).Debug("failed to parse CIDR")
			} else {
				cfg.TrustedSubnet = ipv4Net
			}
		}

		if c.IsTLSEnabled != nil {
			cfg.IsTLSEnabled = *c.IsTLSEnabled
		}
	}
}
