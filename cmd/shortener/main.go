package main

import (
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/Galish/url-shortener/internal/app/repository/db"
	"github.com/Galish/url-shortener/internal/app/repository/filestore"
	"github.com/Galish/url-shortener/internal/app/repository/kvstore"
	"github.com/Galish/url-shortener/internal/app/server"
)

func main() {
	cfg := config.New()

	logger.Initialize(cfg.LogLevel)

	store, err := newRepo(cfg)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	router := handlers.NewRouter(cfg, store)
	httpServer := server.NewHTTPServer(cfg.ServAddr, router)

	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}

func newRepo(cfg *config.Config) (repository.Repository, error) {
	if cfg.DBAddr != "" {
		repo, err := db.New(cfg.DBAddr)
		if err != nil {
			return nil, err
		}

		if err := repo.Bootstrap(); err != nil {
			return nil, err
		}

		return repo, nil
	}

	if cfg.FilePath != "" {
		return filestore.New(cfg.FilePath)
	}

	return kvstore.New(), nil
}
