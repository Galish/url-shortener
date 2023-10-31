package handlers

import (
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/gzip"
	"github.com/Galish/url-shortener/internal/app/middleware"
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

	router.Get(
		"/{id}",
		middleware.RequestLogger(handler.getFullLink),
	)

	router.Post(
		"/",
		middleware.RequestLogger(gzip.WithCompression(handler.shorten)),
	)

	router.Post(
		"/api/shorten",
		middleware.RequestLogger(gzip.WithCompression(handler.apiShorten)),
	)

	return router
}
