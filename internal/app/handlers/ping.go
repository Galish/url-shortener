package handlers

import "net/http"

func (h *httpHandler) ping(w http.ResponseWriter, r *http.Request) {
	ping, err := h.db.Ping()
	if !ping {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
