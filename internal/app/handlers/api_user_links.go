package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
)

func (h *httpHandler) apiUserLinks(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("UserID") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	shortLinks, err := h.repo.GetByUser(r.Context(), r.Header.Get("UserID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WithError(err).Error("unable to read from repository")
	}

	if len(shortLinks) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]apiBatchEntity, 0, len(shortLinks))

	for _, link := range shortLinks {
		resp = append(
			resp,
			apiBatchEntity{
				ShortURL:    fmt.Sprintf("%s/%s", h.cfg.BaseURL, link.Short),
				OriginalURL: link.Original,
			},
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot encode request JSON body")
	}
}
