package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *httpHandler) getFullLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if len(id) < 8 {
		http.Error(w, "invalid identifier", http.StatusBadRequest)
		return
	}

	fullLink, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	http.Redirect(w, r, fullLink, http.StatusTemporaryRedirect)
}
