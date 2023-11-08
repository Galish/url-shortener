package config

type Config struct {
	Addr     string
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
