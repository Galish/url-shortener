package config

type Config struct {
	ServAddr string
	BaseURL  string
	LogLevel string
	FilePath string
	DBAddr   string
}

var cfg Config

func New() *Config {
	parseFlags()
	parseEnvVars()

	return &cfg
}
