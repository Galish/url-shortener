package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr    string
	BaseURL string
}

var cfg Config

func init() {
	flag.StringVar(&cfg.Addr, "a", ":8080", "Server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base link URL")
}

func Get() Config {
	flag.Parse()

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		cfg.Addr = envAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseURL = envBaseURL
	}

	return cfg
}
