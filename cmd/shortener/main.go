package main

import (
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository/db"
	"github.com/Galish/url-shortener/internal/app/repository/filestore"
	"github.com/Galish/url-shortener/internal/app/server"
)

func main() {
	cfg := config.New()

	logger.Initialize(cfg.LogLevel)

	store, err := filestore.New(cfg.FilePath)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	db, err := db.New(cfg.DBAddr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := handlers.NewRouter(cfg, store, db)
	httpServer := server.NewHTTPServer(cfg.Addr, router)

	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
