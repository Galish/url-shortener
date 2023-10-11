package handlers

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

type httpHandler struct {
	cfg   *config.Config
	store storage.Repository
}

func NewRouter(cfg *config.Config, store storage.Repository) http.Handler {
	router := chi.NewRouter()
	handler := httpHandler{cfg, store}

	router.Get("/{id}", handler.getFullLink)
	router.Post("/", handler.makeShortLink)

	return router
}
