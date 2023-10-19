package logger

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type loggerResponseWriter struct {
	rw     http.ResponseWriter
	status int
	size   int
}

func (l *loggerResponseWriter) Write(b []byte) (int, error) {
	size, error := l.rw.Write(b)
	l.size = size

	return size, error
}

func (l *loggerResponseWriter) WriteHeader(statusCode int) {
	l.status = statusCode
	l.rw.WriteHeader(statusCode)
}

func (l *loggerResponseWriter) Header() http.Header {
	return l.rw.Header()
}

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respWriter := loggerResponseWriter{rw: w}

		h(&respWriter, r)

		logger.WithFields(log.Fields{
			"size":     respWriter.size,
			"status":   respWriter.status,
			"duration": time.Since(start),
			"method":   r.Method,
			"uri":      r.RequestURI,
		}).Info("incoming request")
	}
}
