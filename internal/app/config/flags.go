package config

import (
	"flag"
)

func parseFlags(c *settings, configFile *string) {
	var (
		isTLSEnabled bool
		config       string
	)

	flag.StringVar(&c.ServAddr, "a", "", "Server address")
	flag.StringVar(&c.BaseURL, "b", "", "Base link URL")
	flag.StringVar(&c.LogLevel, "l", "", "Log level")
	flag.StringVar(&c.FilePath, "f", "", "File storage path")
	flag.StringVar(&c.DBAddr, "d", "", "DB address")
	flag.BoolVar(&isTLSEnabled, "s", false, "Enable HTTPS")
	flag.StringVar(&config, "c", "", "Config file path")
	flag.StringVar(&c.TrustedSubnet, "t", "", "Trusted subnet in CIDR notation")
	flag.StringVar(&config, "config", "", "Config file path")

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "s":
			c.IsTLSEnabled = &isTLSEnabled

		case "c", "config":
			*configFile = f.Value.String()
		}
	})
}
