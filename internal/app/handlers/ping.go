package handlers

import (
	"net/http"

	"github.com/Galish/url-shortener/pkg/logger"
)

// ping verifies the server is running.
func (h *HTTPHandler) ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if ping, err := h.repo.Ping(ctx); !ping {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WithError(err).Debug("unable to ping repository")
		return
	}

	w.WriteHeader(http.StatusOK)
}
