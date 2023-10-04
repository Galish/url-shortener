package config

import "flag"

type Config struct {
	Addr    string
	BaseURL string
}

var cfg Config

func init() {
	flag.StringVar(&cfg.Addr, "a", ":8080", "Server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080/", "Base link URL")
}

func Get() Config {
	flag.Parse()
	return cfg
}
