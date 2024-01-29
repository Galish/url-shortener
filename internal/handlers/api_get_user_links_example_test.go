package handlers_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/Galish/url-shortener/internal/config"
	"github.com/Galish/url-shortener/internal/handlers"
	"github.com/Galish/url-shortener/internal/repository/memstore"
	"github.com/Galish/url-shortener/internal/repository/model"
)

func ExampleHTTPHandler_APIGetUserLinks() {
	r, _ := http.NewRequest(
		http.MethodGet,
		"/api/user/urls",
		nil,
	)

	r.Header.Add("X-User", "e44d9088-1bd6-44dc-af86-f1a551b02db3")

	w := httptest.NewRecorder()

	store := memstore.New()

	store.Set(
		context.Background(),
		&model.ShortLink{
			Short:    "Edz0Thb1",
			Original: "https://practicum.yandex.ru/",
			User:     "e44d9088-1bd6-44dc-af86-f1a551b02db3",
		},
	)

	apiHandler := handlers.NewHandler(
		&config.Config{BaseURL: "http://www.shortener.io"},
		store,
	)

	apiHandler.APIGetUserLinks(w, r)

	resp := w.Result()

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Output:
	// 200
	// application/json
	// [{"original_url":"https://practicum.yandex.ru/","short_url":"http://www.shortener.io/Edz0Thb1"}]
}
