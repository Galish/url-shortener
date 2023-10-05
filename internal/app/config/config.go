package config

type Config struct {
	Addr    string
	BaseURL string
}

var cfg Config

func Get() Config {
	parseFlags()
	parseEnvVars()

	return cfg
}
