package handlers

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

type httpHandler struct {
	store storage.KeyValueStorage
}

func NewHandler(store storage.KeyValueStorage) http.Handler {
	router := chi.NewRouter()
	handler := httpHandler{store}

	router.Get("/{id}", handler.getFullLink)
	router.Post("/", handler.makeShortLink)

	return router
}
