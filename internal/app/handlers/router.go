package handlers

import (
	"github.com/Galish/url-shortener/internal/app/compress"
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/middleware"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/go-chi/chi/v5"
)

func NewRouter(cfg *config.Config, repo repository.Repository) *chi.Mux {
	handler := NewHandler(cfg, repo)
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithRequestLogger)

		r.Get("/ping", handler.ping)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithAuthToken)
		r.Use(middleware.WithRequestLogger)

		r.Get("/{id}", handler.getFullLink)

		r.Group(func(r chi.Router) {
			r.Use(middleware.WithCompressor(compress.NewGzipCompressor()))

			r.Post("/", handler.shorten)

			r.Route("/api/user/urls", func(r chi.Router) {
				r.Use(middleware.WithAuthChecker)

				r.Get("/", handler.apiGetUserLinks)
				r.Delete("/", handler.apiDeleteUserLinks)
			})

			r.Route("/api/shorten", func(r chi.Router) {
				r.Post("/", handler.apiShorten)
				r.Post("/batch", handler.apiShortenBatch)
			})
		})
	})

	return router
}
