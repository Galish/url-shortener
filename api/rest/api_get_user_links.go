package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/middleware"
	"github.com/Galish/url-shortener/pkg/logger"
)

// APIGetByUser is an API handler that returns a list of URLs created by the user.
//
//	POST /api/user/urls
func (h *HTTPHandler) APIGetByUser(w http.ResponseWriter, r *http.Request) {
	urls, err := h.usecase.GetByUser(r.Context(), r.Header.Get(middleware.AuthHeaderName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WithError(err).Error("unable to fetch user urls")
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]APIBatchEntity, 0, len(urls))

	for _, url := range urls {
		resp = append(
			resp,
			APIBatchEntity{
				ShortURL:    fmt.Sprintf("%s/%s", h.cfg.BaseURL, url.Short),
				OriginalURL: url.Original,
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
