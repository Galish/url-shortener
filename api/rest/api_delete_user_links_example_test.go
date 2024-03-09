package restapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	restapi "github.com/Galish/url-shortener/api/rest"
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func ExampleHTTPHandler_APIDeleteUserURLs() {
	bodyRaw, err := json.Marshal([]string{"Edz0Thb1"})
	if err != nil {
		log.Fatal(err)
	}

	r, _ := http.NewRequest(
		http.MethodDelete,
		"/api/user/urls",
		bytes.NewBuffer(bodyRaw),
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

	uc := usecase.New(memstore.New())
	defer uc.Close()

	apiHandler := restapi.NewHandler(
		&config.Config{BaseURL: "http://www.shortener.io"},
		uc,
		nil,
	)

	apiHandler.APIDeleteUserURLs(w, r)

	resp := w.Result()

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Output:
	// 202
	//
	//
}
