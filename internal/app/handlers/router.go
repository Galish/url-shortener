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
		),
	)

	router.Get(
		"/{id}",
		middleware.Apply(
			handler.getFullLink,
			middleware.WithRequestLogger,
		),
	)

	router.Get(
		"/api/user/urls",
		middleware.Apply(
			handler.apiUserLinks,
			middleware.WithRequestLogger,
			middleware.WithCompressor(compressor),
			middleware.WithAuthChecker,
		),
	)

	router.Post(
		"/",
		middleware.Apply(
			handler.shorten,
			middleware.WithCompressor(compressor),
			middleware.WithAuthToken,
			middleware.WithRequestLogger,
		),
	)

	router.Post(
		"/api/shorten",
		middleware.Apply(
			handler.apiShorten,
			middleware.WithCompressor(compressor),
			middleware.WithAuthToken,
			middleware.WithRequestLogger,
		),
	)

	router.Post(
		"/api/shorten/batch",
		middleware.Apply(
			handler.apiShortenBatch,
			middleware.WithCompressor(compressor),
			middleware.WithAuthToken,
			middleware.WithRequestLogger,
		),
	)

	return router
}
