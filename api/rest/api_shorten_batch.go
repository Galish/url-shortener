package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/middleware"
	"github.com/Galish/url-shortener/pkg/logger"
)

// APIShortenBatch is an API handler for creating short URLs in batches.
//
//	POST /api/shorten/batch
func (h *HTTPHandler) APIShortenBatch(w http.ResponseWriter, r *http.Request) {
	var req []APIBatchEntity
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "cannot decode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot decode request JSON body")
		return
	}

	if len(req) == 0 {
		http.Error(w, "empty request body", http.StatusBadRequest)
		return
	}

	user := r.Header.Get(middleware.AuthHeaderName)
	urls := make([]*entity.URL, len(req))

	for i, item := range req {
		urls[i] = &entity.URL{
			Original: item.OriginalURL,
			User:     user,
		}
	}

	err := h.usecase.Shorten(r.Context(), urls...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := make([]APIBatchEntity, len(req))
	for i, row := range req {
		resp[i] = APIBatchEntity{
			CorrelationID: row.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", h.cfg.BaseURL, urls[i].Short),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot encode request JSON body")
	}
}
