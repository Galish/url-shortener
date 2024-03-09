package restapi_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	restapi "github.com/Galish/url-shortener/api/rest"
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func ExampleHTTPHandler_APIGetByUser() {
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
		&entity.URL{
			Short:    "Edz0Thb1",
			Original: "https://practicum.yandex.ru/",
			User:     "e44d9088-1bd6-44dc-af86-f1a551b02db3",
		},
	)

	uc := usecase.New(store)
	defer uc.Close()

	apiHandler := restapi.NewHandler(
		&config.Config{BaseURL: "http://www.shortener.io"},
		uc,
		nil,
	)

	apiHandler.APIGetByUser(w, r)

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
