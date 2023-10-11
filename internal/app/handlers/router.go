package handlers

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/go-chi/chi/v5"
)

type httpHandler struct {
	cfg  *config.Config
	repo repository.Repository
}

func NewRouter(cfg *config.Config, repo repository.Repository) http.Handler {
	router := chi.NewRouter()
	handler := httpHandler{cfg, repo}

	router.Get("/{id}", handler.getFullLink)
	router.Post("/", handler.makeShortLink)

	return router
}
