package handlers

import "net/http"

func (h *httpHandler) ping(w http.ResponseWriter, r *http.Request) {
	if ping, err := h.repo.Ping(); !ping {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
