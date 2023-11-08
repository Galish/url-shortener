package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
)

type apiBatchEntity struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}

func (h *httpHandler) apiShortenBatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req []apiBatchEntity
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "cannot decode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot decode request JSON body")
		return
	}

	if len(req) == 0 {
		http.Error(w, "empty request body", http.StatusBadRequest)
		return
	}

	resp := make([]apiBatchEntity, 0, len(req))
	entries := make([][2]string, 0, len(req))

	for _, entity := range req {
		if entity.OriginalURL == "" {
			http.Error(w, "link not provided", http.StatusBadRequest)
			return
		}

		id := h.generateUniqueID(ctx, idLength)

		resp = append(
			resp,
			apiBatchEntity{
				CorrelationID: entity.CorrelationID,
				ShortURL:      fmt.Sprintf("%s/%s", h.cfg.BaseURL, id),
			},
		)

		entries = append(entries, [2]string{id, entity.OriginalURL})
	}

	if err := h.repo.SetBatch(ctx, entries...); err != nil {
		http.Error(w, "unable to write to repository", http.StatusInternalServerError)
		logger.WithError(err).Debug("unable to write to repository")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot encode request JSON body")
	}
}
