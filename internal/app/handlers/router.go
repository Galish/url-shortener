package handlers

import (
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/go-chi/chi/v5"
)

type httpHandler struct {
	cfg  *config.Config
	repo repository.Repository
}

func NewRouter(cfg *config.Config, repo repository.Repository) *chi.Mux {
	router := chi.NewRouter()
	handler := httpHandler{cfg, repo}

	router.Get("/{id}", logger.WithLogging(handler.getFullLink))
	router.Post("/", logger.WithLogging(handler.makeShortLink))

	return router
}
