package handlers

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
)

func (h *httpHandler) ping(w http.ResponseWriter, r *http.Request) {
	if ping, err := h.repo.Ping(); !ping {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WithError(err).Debug("unable to ping repository")
		return
	}

	w.WriteHeader(http.StatusOK)
}
