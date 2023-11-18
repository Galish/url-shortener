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
			middleware.WithAuthToken,
			middleware.WithRequestLogger,
		),
	)

	router.Get(
		"/api/user/urls",
		middleware.Apply(
			handler.apiUserLinks,
			middleware.WithAuthToken,
			middleware.WithAuthChecker,
			middleware.WithCompressor(compressor),
			middleware.WithRequestLogger,
		),
	)

	router.Post(
		"/",
		middleware.Apply(
			handler.shorten,
			middleware.WithAuthToken,
			middleware.WithCompressor(compressor),
			middleware.WithRequestLogger,
		),
	)

	router.Post(
		"/api/shorten",
		middleware.Apply(
			handler.apiShorten,
			middleware.WithAuthToken,
			middleware.WithCompressor(compressor),
			middleware.WithRequestLogger,
		),
	)

	router.Post(
		"/api/shorten/batch",
		middleware.Apply(
			handler.apiShortenBatch,
			middleware.WithAuthToken,
			middleware.WithCompressor(compressor),
			middleware.WithRequestLogger,
		),
	)

	return router
}
