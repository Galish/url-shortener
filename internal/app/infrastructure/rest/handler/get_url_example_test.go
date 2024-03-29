package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/infrastructure/rest/handler"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func ExampleHandler_Get() {
	r, _ := http.NewRequest(
		http.MethodGet,
		"/Edz0Thb1",
		nil,
	)

	w := httptest.NewRecorder()

	store := memstore.New()

	store.Set(
		context.Background(),
		&entity.URL{
			Short:    "Edz0Thb1",
			Original: "https://practicum.yandex.ru/",
		},
	)

	uc := usecase.New(
		&config.Config{BaseURL: "http://www.shortener.io"},
		store,
	)
	defer uc.Close()

	apiHandler := handler.New(uc, nil)
	apiHandler.Get(w, r)

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
