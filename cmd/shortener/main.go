package main

import (
	"net/http"

	"github.com/url-shortener/cmd/shortener/handler"
	"github.com/url-shortener/cmd/shortener/storage"
)

func main() {
	store := storage.NewKeyValueStorage()
	httpHandler := handler.NewHandler(store)

	err := http.ListenAndServe(`:8080`, httpHandler)
	if err != nil {
		panic(err)
	}
}
