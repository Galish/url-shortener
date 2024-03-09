package restapi

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/usecase"
)

type Pinger interface {
	Ping(context.Context) (bool, error)
}

// HTTPHandler represents API handler.
type HTTPHandler struct {
	repo    Pinger
	usecase usecase.Shortener
}

// NewHandler implements HTTP handlers.
func NewHandler(usecase usecase.Shortener, repo Pinger) *HTTPHandler {
	return &HTTPHandler{
		usecase: usecase,
		repo:    repo,
	}
}
