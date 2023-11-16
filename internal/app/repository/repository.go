package repository

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/repository/db"
	"github.com/Galish/url-shortener/internal/app/repository/filestore"
	"github.com/Galish/url-shortener/internal/app/repository/kvstore"
)

type Repository interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string) error
	SetBatch(context.Context, ...[]interface{}) error
	Has(context.Context, string) bool
	Ping(context.Context) (bool, error)
	Close() error
}

func New(cfg *config.Config) (Repository, error) {
	if cfg.DBAddr != "" {
		repo, err := db.New(cfg.DBAddr)
		if err != nil {
			return nil, err
		}

		if err := repo.Bootstrap(context.Background()); err != nil {
			return nil, err
		}

		return repo, nil
	}

	if cfg.FilePath != "" {
		return filestore.New(cfg.FilePath)
	}

	return kvstore.New(), nil
}
