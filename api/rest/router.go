// Package implements the HTTP router and handlers.
package restapi

import (
	"github.com/go-chi/chi/v5"

	"github.com/Galish/url-shortener/internal/app/middleware"
	"github.com/Galish/url-shortener/pkg/compress"
)

// NewRouter returns a new Mux object that implements the Router interface.
func NewRouter(handler *HTTPHandler) *chi.Mux {
	router := chi.NewRouter()

	var withCompression = middleware.WithCompressor(compress.NewGzipCompressor())

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithRequestLogger)

		r.Get("/ping", handler.ping)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithAuthToken)
		r.Use(middleware.WithRequestLogger)

		r.Get("/{id}", handler.GetFullLink)

		r.Group(func(r chi.Router) {
			r.Use(withCompression)

			r.Post("/", handler.Shorten)

			r.Route("/api/user/urls", func(r chi.Router) {
				r.Use(middleware.WithAuthChecker)

				r.Get("/", handler.APIGetUserLinks)
				r.Delete("/", handler.APIDeleteUserLinks)
			})

			r.Route("/api/shorten", func(r chi.Router) {
				r.Post("/", handler.APIShorten)
				r.Post("/batch", handler.APIShortenBatch)
			})
		})
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithTrustedSubnet(handler.cfg.TrustedSubnet))
		r.Use(middleware.WithRequestLogger)
		r.Use(withCompression)

		r.Get("/api/internal/stats", handler.APIStats)
	})

	return router
}
