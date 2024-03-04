package config

type settings struct {
	BaseURL      string `json:"base_url"`
	DBAddr       string `json:"database_dsn"`
	FilePath     string `json:"file_storage_path"`
	IsTLSEnabled *bool  `json:"enable_https"`
	LogLevel     string `json:"log_level"`
	ServAddr     string `json:"server_address"`
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

		if c.IsTLSEnabled != nil {
			cfg.IsTLSEnabled = *c.IsTLSEnabled
		}
	}
}
