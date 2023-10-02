package handlers

import (
	"net/http"
)

func (h *httpHandler) getFullLink(w http.ResponseWriter, r *http.Request) {
	fullLink, err := h.store.Get(r.URL.Path[1:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	http.Redirect(w, r, fullLink, http.StatusTemporaryRedirect)
}
