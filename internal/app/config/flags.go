package config

import "flag"

func init() {
	flag.StringVar(&cfg.Addr, "a", ":8080", "Server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base link URL")
	flag.StringVar(&cfg.LogLevel, "l", "info", "Log level")
	flag.StringVar(&cfg.FilePath, "f", "/tmp/short-url-db.json", "File storage path")
}

func parseFlags() {
	flag.Parse()
}
