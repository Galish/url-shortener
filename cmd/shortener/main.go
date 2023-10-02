package main

import (
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/server"
	"github.com/Galish/url-shortener/internal/app/storage"
)

func main() {
	store := storage.NewKeyValueStorage()
	httpServer := server.NewHttpServer(`:8080`, handlers.NewHandler(store))

	err := httpServer.Run()
	if err != nil {
		panic(err)
	}
}
