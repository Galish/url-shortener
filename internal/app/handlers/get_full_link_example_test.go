package handlers_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

func ExampleHTTPHandler_GetFullLink() {
	store := memstore.New()
	store.Set(
		context.Background(),
		&model.ShortLink{
			Short:    "Edz0Thb1",
			Original: "https://practicum.yandex.ru/",
		},
	)

	apiHandler := handlers.NewHandler(
		&config.Config{BaseURL: "http://www.shortener.io"},
		store,
	)

	router := handlers.NewRouter(apiHandler)
	server := httptest.NewServer(router)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/Edz0Thb1", nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(resp.Header.Get("Location"))

	// Output:
	// 307
	// text/html; charset=utf-8
	// https://practicum.yandex.ru/
}
