package config

import "flag"

func init() {
	flag.StringVar(&cfg.Addr, "a", ":8080", "Server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base link URL")
}

func parseFlags() {
	flag.Parse()
}
