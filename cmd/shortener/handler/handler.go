package handler

import (
	"net/http"

	"github.com/url-shortener/cmd/shortener/storage"
)

func NewHandler(store storage.KeyValueStorage) http.Handler {
	service := shortenerService{store}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			service.makeShortLink(w, r)

		case http.MethodGet:
			service.getFullLink(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	})
}
