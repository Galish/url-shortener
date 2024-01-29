package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/middleware"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

// apiShortenBatch is an API handler for creating short links in batches.
func (h *HttpHandler) apiShortenBatch(w http.ResponseWriter, r *http.Request) {
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

	resp := make([]apiBatchEntity, len(req))
	rows := make([]*model.ShortLink, len(req))

	for i, entity := range req {
		if entity.OriginalURL == "" {
			http.Error(w, "link not provided", http.StatusBadRequest)
			return
		}

		id := h.generateUniqueID(ctx, idLength)

		resp[i] = apiBatchEntity{
			CorrelationID: entity.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", h.cfg.BaseURL, id),
		}

		rows[i] = &model.ShortLink{
			Short:    id,
			Original: entity.OriginalURL,
			User:     r.Header.Get(middleware.AuthHeaderName),
		}
	}

	if err := h.repo.SetBatch(ctx, rows...); err != nil {
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
