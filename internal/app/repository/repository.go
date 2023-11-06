package repository

import (
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/repository/db"
	"github.com/Galish/url-shortener/internal/app/repository/filestore"
	"github.com/Galish/url-shortener/internal/app/repository/kvstore"
)

type Repository interface {
	Get(string) (string, error)
	Set(string, string) error
	SetBatch(...[2]string) error
	Has(string) bool
	Ping() (bool, error)
	Close() error
}

func New(cfg *config.Config) (Repository, error) {
	if cfg.DBAddr != "" {
		return db.New(cfg.DBAddr)
	}

	if cfg.FilePath != "" {
		return filestore.New(cfg.FilePath)
	}

	return kvstore.New(), nil
}
