package handlers

import (
	"github.com/Galish/url-shortener/internal/app/compress"
	"github.com/Galish/url-shortener/internal/app/config"
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
	compressor := compress.NewGzipCompressor()

	router.Get(
		"/ping",
		middleware.WithRequestLogger(handler.ping),
	)

	router.Get(
		"/{id}",
		middleware.WithRequestLogger(handler.getFullLink),
	)

	router.Post(
		"/",
		middleware.WithRequestLogger(
			middleware.WithCompression(handler.shorten, compressor),
		),
	)

	router.Post(
		"/api/shorten",
		middleware.WithRequestLogger(
			middleware.WithCompression(handler.apiShorten, compressor),
		),
	)

	router.Post(
		"/api/shorten/batch",
		middleware.WithRequestLogger(
			middleware.WithCompression(handler.apiShortenBatch, compressor),
		),
	)

	return router
}
