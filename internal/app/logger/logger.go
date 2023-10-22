package logger

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func init() {
	logger = log.New()
	logger.SetOutput(io.Discard)
}

func Initialize(level string) {
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stderr)

	logLevel, err := log.ParseLevel(level)
	if err != nil {
		logger.Error("invalid log level", err)
		return
	}

	logger.SetLevel(logLevel)
}

func WithError(err error) *log.Entry {
	return log.WithError(err)
}
