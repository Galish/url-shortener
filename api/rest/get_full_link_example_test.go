package restapi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	restapi "github.com/Galish/url-shortener/api/rest"
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
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
		&entity.ShortLink{
			Short:    "Edz0Thb1",
			Original: "https://practicum.yandex.ru/",
		},
	)

	apiHandler := restapi.NewHandler(
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
