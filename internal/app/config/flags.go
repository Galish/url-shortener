package config

import (
	"flag"
)

func parseFlags(s *settings) {
	flag.StringVar(&s.ServAddr, "a", "", "Server address")
	flag.StringVar(&s.BaseURL, "b", "", "Base link URL")
	flag.StringVar(&s.logLevel, "l", "", "Log level")
	flag.StringVar(&s.FilePath, "f", "", "File storage path")
	flag.StringVar(&s.DBAddr, "d", "", "DB address")
	flag.BoolVar(&s.IsTLSEnabled, "s", false, "Enable HTTPS")

	flag.Func("c", "Config file path", func(v string) error {
		s.fileConfigPath = v
		return nil
	})

	flag.Func("config", "Config file path", func(v string) error {
		s.fileConfigPath = v
		return nil
	})

	flag.Parse()
}
