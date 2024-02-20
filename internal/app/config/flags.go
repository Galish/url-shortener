package config

import "flag"

func init() {
	flag.StringVar(&cfg.ServAddr, "a", ":8080", "Server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base link URL")
	flag.StringVar(&cfg.LogLevel, "l", "info", "Log level")
	flag.StringVar(&cfg.FilePath, "f", "/tmp/short-url-db.json", "File storage path")
	flag.StringVar(&cfg.DBAddr, "d", "", "DB address")
	flag.BoolVar(&cfg.IsTLSEnabled, "s", false, "Enable HTTPS")
}

func parseFlags() {
	flag.Parse()
}
