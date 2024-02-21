package config

type settings struct {
	BaseURL        string `json:"base_url"`
	DBAddr         string `json:"database_dsn"`
	FilePath       string `json:"file_storage_path"`
	IsTLSEnabled   bool   `json:"enable_https"`
	ServAddr       string `json:"server_address"`
	fileConfigPath string
	logLevel       string
}

var defaultSettings = &settings{
	BaseURL:  "http://localhost:8080",
	FilePath: "/tmp/short-url-db.json",
	ServAddr: ":8080",
	logLevel: "info",
}

func withSettings(s *settings) func(*Config) {
	return func(cfg *Config) {
		if s.ServAddr != "" {
			cfg.ServAddr = s.ServAddr
		}

		if s.BaseURL != "" {
			cfg.BaseURL = s.BaseURL
		}

		if s.logLevel != "" {
			cfg.LogLevel = s.logLevel
		}

		if s.FilePath != "" {
			cfg.FilePath = s.FilePath
		}

		if s.DBAddr != "" {
			cfg.DBAddr = s.DBAddr
		}

		if s.IsTLSEnabled {
			cfg.IsTLSEnabled = true
		}
	}
}
