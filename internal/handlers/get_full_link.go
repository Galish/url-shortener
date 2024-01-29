package handlers

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/logger"

	"github.com/go-chi/chi/v5"
)

// GetFullLink redirects to the original page URL for the given short link.
//
//	GET /%s
func (h *HTTPHandler) GetFullLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		id = r.URL.Path[1:]
	}

	if len(id) < 8 {
		http.Error(w, "invalid identifier", http.StatusBadRequest)
		return
	}

	shortLink, err := h.repo.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.WithError(err).Debug("unable to read from repository")
		return
	}

	if shortLink.IsDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	http.Redirect(w, r, shortLink.Original, http.StatusTemporaryRedirect)
}
