package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Galish/url-shortener/pkg/logger"
)

// APIStats is an API handler that returns the number shortened URLs and users.
//
//	GET /api/internal/stats
func (h *HTTPHandler) APIStats(w http.ResponseWriter, r *http.Request) {
	urls, users, err := h.repo.Stats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WithError(err).Error("unable to read from repository")
	}

	stats := APIStatsResponse{
		Urls:  urls,
		Users: users,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot encode request JSON body")
	}
}
