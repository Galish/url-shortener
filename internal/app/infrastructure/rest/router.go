// Package implements the HTTP router and handlers.
package restapi

import (
	"github.com/go-chi/chi/v5"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/infrastructure/rest/handler"
	"github.com/Galish/url-shortener/internal/app/infrastructure/rest/middleware"
	"github.com/Galish/url-shortener/pkg/compress"
)

// NewRouter returns a new Mux object that implements the Router interface.
func NewRouter(cfg *config.Config, h *handler.Handler) *chi.Mux {
	router := chi.NewRouter()

	var withCompression = middleware.WithCompressor(compress.NewGzipCompressor())

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithRequestLogger)

		r.Get("/ping", h.Ping)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithAuthToken)
		r.Use(middleware.WithRequestLogger)

		r.Get("/{id}", h.Get)

		r.Group(func(r chi.Router) {
			r.Use(withCompression)

			r.Post("/", h.Shorten)

			r.Route("/api/user/urls", func(r chi.Router) {
				r.Use(middleware.WithAuthChecker)

				r.Get("/", h.APIGetByUser)
				r.Delete("/", h.APIDeleteUserURLs)
			})

			r.Route("/api/shorten", func(r chi.Router) {
				r.Post("/", h.APIShorten)
				r.Post("/batch", h.APIShortenBatch)
			})
		})
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithTrustedSubnet(cfg.TrustedSubnet))
		r.Use(middleware.WithRequestLogger)
		r.Use(withCompression)

		r.Get("/api/internal/stats", h.APIStats)
	})

	return router
}
