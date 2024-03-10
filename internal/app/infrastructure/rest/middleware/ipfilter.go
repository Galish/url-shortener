package middleware

import (
	"net"
	"net/http"

	"github.com/Galish/url-shortener/pkg/logger"
)

// WithTrustedSubnet serves 401 error if request IP verification fails.
func WithTrustedSubnet(ipNet *net.IPNet) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := net.ParseIP(r.Header.Get("X-Real-IP"))

			if ipNet == nil || !ipNet.Contains(ip) {
				logger.Debug("unauthorized access attempt")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
