package handlers

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/storage"
)

type httpHandler struct {
	store storage.KeyValueStorage
}

func NewHandler(store storage.KeyValueStorage) http.Handler {
	handler := httpHandler{store}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.makeShortLink(w, r)

		case http.MethodGet:
			handler.getFullLink(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	})
}
