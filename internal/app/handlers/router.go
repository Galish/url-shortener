package handlers

import (
	"context"
	"time"

	"github.com/Galish/url-shortener/internal/app/compress"
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/middleware"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/Galish/url-shortener/internal/app/repository/models"
	"github.com/go-chi/chi/v5"
)

type httpHandler struct {
	cfg      *config.Config
	repo     repository.Repository
	deleteCh chan *models.ShortLink
}

func NewRouter(cfg *config.Config, repo repository.Repository) *chi.Mux {
	router := chi.NewRouter()
	handler := httpHandler{
		cfg:      cfg,
		repo:     repo,
		deleteCh: make(chan *models.ShortLink, 100),
	}
	go handler.flushMessages()

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
			handler.apiGetUserLinks,
			middleware.WithAuthToken,
			middleware.WithAuthChecker,
			middleware.WithCompressor(compressor),
			middleware.WithRequestLogger,
		),
	)

	router.Delete(
		"/api/user/urls",
		middleware.Apply(
			handler.apiDeleteUserLinks,
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

func (h *httpHandler) flushMessages() {
	ticker := time.NewTicker(2 * time.Second)

	var list []*models.ShortLink

	for {
		select {
		case shortLink := <-h.deleteCh:
			list = append(list, shortLink)
		case <-ticker.C:
			if len(list) == 0 {
				continue
			}

			if err := h.repo.Delete(context.TODO(), list...); err != nil {
				logger.WithError(err).Debug("cannot delete messages")
				continue
			}

			list = nil
		}
	}
}
