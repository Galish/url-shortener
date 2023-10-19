package config

type Config struct {
	Addr     string
	BaseURL  string
	LogLevel string
}

var cfg Config

func New() *Config {
	parseFlags()
	parseEnvVars()

	return &cfg
}
