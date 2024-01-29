package middleware

import (
	"net/http"
	"time"

	"github.com/Galish/url-shortener/internal/logger"
)

type loggerResponseWriter struct {
	rw     http.ResponseWriter
	status int
	size   int
}

// WithRequestLogger provides request logging.
func WithRequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respWriter := loggerResponseWriter{rw: w}

		logger.WithFields(logger.Fields{
			"size":     respWriter.size,
			"status":   respWriter.status,
			"duration": time.Since(start),
			"method":   r.Method,
			"uri":      r.RequestURI,
			"user":     r.Header.Get(AuthHeaderName),
		}).Info("incoming request")

		h.ServeHTTP(&respWriter, r)
	})
}

// Write overrides response Write method.
func (l *loggerResponseWriter) Write(b []byte) (int, error) {
	size, error := l.rw.Write(b)
	l.size = size

	return size, error
}

// WriteHeader overrides response WriteHeader method.
func (l *loggerResponseWriter) WriteHeader(statusCode int) {
	l.status = statusCode
	l.rw.WriteHeader(statusCode)
}

// Header overrides response Header method.
func (l *loggerResponseWriter) Header() http.Header {
	return l.rw.Header()
}
