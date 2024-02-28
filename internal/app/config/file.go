package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/Galish/url-shortener/pkg/logger"
)

func parseFile(fileName string, c *settings) {
	if fileName == "" {
		return
	}

	f, err := os.Open(fileName)
	if err != nil {
		logger.WithError(err).Debug("failed to open configuration file")
		return
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		logger.WithError(err).Debug("failed to read configuration file")
		return
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		logger.WithError(err).Debug("failed to decode configuration file")
		return
	}
}
