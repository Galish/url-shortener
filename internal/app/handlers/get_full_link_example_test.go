package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

func ExampleHTTPHandler_GetFullLink() {
	r, _ := http.NewRequest(
		http.MethodGet,
		"/Edz0Thb1",
		nil,
	)

	w := httptest.NewRecorder()

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
	defer apiHandler.Close()

	apiHandler.GetFullLink(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(resp.Header.Get("Location"))

	// Output:
	// 307
	// text/html; charset=utf-8
	// https://practicum.yandex.ru/
}
