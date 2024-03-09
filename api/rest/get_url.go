package restapi

import (
	"net/http"

	"github.com/Galish/url-shortener/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// Get redirects to the original page URL for the given short URL.
//
//	GET /%s
func (h *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		id = r.URL.Path[1:]
	}

	url, err := h.usecase.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.WithError(err).Debug("unable to fetch url")
		return
	}

	if url.IsDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	http.Redirect(w, r, url.Original, http.StatusTemporaryRedirect)
}
