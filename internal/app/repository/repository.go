// Package repository implements  the persistence layer of the application.
package repository

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/repository/db"
	"github.com/Galish/url-shortener/internal/app/repository/filestore"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

// Repository is an abstraction over persistent storage.
// It provides a common set of methods for interacting with data sources.
type Repository interface {
	Get(context.Context, string) (*model.ShortLink, error)
	GetByUser(context.Context, string) ([]*model.ShortLink, error)
	Set(context.Context, *model.ShortLink) error
	SetBatch(context.Context, ...*model.ShortLink) error
	Delete(context.Context, ...*model.ShortLink) error
	Has(context.Context, string) bool
	Ping(context.Context) (bool, error)
	Close() error
}

// New creates a store based on the configuration.
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

	return memstore.New(), nil
}
