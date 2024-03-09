package restapi

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/entity"
)

type Usecase interface {
	Shorten(context.Context, ...*entity.URL) error
	ShortURL(*entity.URL) string
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
	repo    Pinger
	usecase Usecase
}

// NewHandler implements HTTP handlers.
func NewHandler(usecase Usecase, repo Pinger) *HTTPHandler {
	return &HTTPHandler{
		usecase: usecase,
		repo:    repo,
	}
}
