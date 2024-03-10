package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/infrastructure/rest/middleware"
	"github.com/Galish/url-shortener/pkg/logger"
)

// APIDeleteUserURLs is an API handler for deleting user URLs.
//
//	DELETE /api/user/urls
func (h *Handler) APIDeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	var ids []string
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "cannot decode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot decode request JSON body")
		return
	}

	user := r.Header.Get(middleware.AuthHeaderName)
	urls := make([]*entity.URL, len(ids))

	for i, id := range ids {
		urls[i] = &entity.URL{
			Short: id,
			User:  user,
		}
	}

	h.usecase.Delete(r.Context(), urls)

	w.WriteHeader(http.StatusAccepted)
}
