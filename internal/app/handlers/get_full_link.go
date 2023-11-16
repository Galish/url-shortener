package handlers

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/go-chi/chi/v5"
)

func (h *httpHandler) getFullLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if len(id) < 8 {
		http.Error(w, "invalid identifier", http.StatusBadRequest)
		return
	}

	fullLink, err := h.repo.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.WithError(err).Debug("unable to read from repository")
		return
	}

	http.Redirect(w, r, fullLink, http.StatusTemporaryRedirect)
}
