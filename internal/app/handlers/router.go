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
		middleware.Apply(
			handler.ping,
			middleware.WithRequestLogger,
			middleware.WithAuthentication,
		),
	)

	router.Get(
		"/{id}",
		middleware.Apply(
			handler.getFullLink,
			middleware.WithRequestLogger,
			middleware.WithAuthentication,
		),
	)

	router.Post(
		"/",
		middleware.Apply(
			handler.shorten,
			middleware.WithRequestLogger,
			middleware.WithCompressor(compressor),
			middleware.WithAuthentication,
		),
	)

	router.Post(
		"/api/shorten",
		middleware.Apply(
			handler.apiShorten,
			middleware.WithRequestLogger,
			middleware.WithCompressor(compressor),
			middleware.WithAuthentication,
		),
	)

	router.Post(
		"/api/shorten/batch",
		middleware.Apply(
			handler.apiShortenBatch,
			middleware.WithRequestLogger,
			middleware.WithCompressor(compressor),
			middleware.WithAuthentication,
		),
	)

	router.Get(
		"/api/user/urls",
		middleware.Apply(
			handler.apiUserLinks,
			middleware.WithRequestLogger,
			middleware.WithCompressor(compressor),
			middleware.WithAuthentication,
		),
	)

	return router
}
