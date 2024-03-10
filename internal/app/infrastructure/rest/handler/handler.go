package handler

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/usecase"
)

type Pinger interface {
	Ping(context.Context) (bool, error)
}

// Handler represents REST API handler.
type Handler struct {
	repo    Pinger
	usecase usecase.Shortener
}

// New implements HTTP handlers.
func New(usecase usecase.Shortener, repo Pinger) *Handler {
	return &Handler{
		usecase: usecase,
		repo:    repo,
	}
}
