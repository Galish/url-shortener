package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Galish/url-shortener/pkg/logger"
)

// APIStats is an API handler that returns the number of users and shortened URLs.
//
//	GET /api/internal/stats
func (h *Handler) APIStats(w http.ResponseWriter, r *http.Request) {
	urls, users, err := h.usecase.GetStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WithError(err).Error("unable to get stats")
		return
	}

	resp := APIStatsResponse{
		Urls:  urls,
		Users: users,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot encode request JSON body")
	}
}
