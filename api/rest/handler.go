package restapi

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/entity"
)

type Usecase interface {
	Shorten(context.Context, ...*entity.URL) error
	Get(context.Context, string) (*entity.URL, error)
	GetByUser(context.Context, string) ([]*entity.URL, error)
	GetStats(context.Context) (int, int, error)
	Delete(context.Context, []*entity.URL)
}

type Pinger interface {
	Ping(context.Context) (bool, error)
}

// HTTPHandler represents API handler.
type HTTPHandler struct {
	cfg     *config.Config
	repo    Pinger
	usecase Usecase
}

// NewHandler implements HTTP handlers.
func NewHandler(cfg *config.Config, usecase Usecase, repo Pinger) *HTTPHandler {
	return &HTTPHandler{
		cfg:     cfg,
		usecase: usecase,
		repo:    repo,
	}
}
