package handlers

import (
	"net/http"
)

func (h *httpHandler) getFullLink(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]

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
