package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

func ExampleHTTPHandler_APIDeleteUserLinks() {
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
	defer apiHandler.Close()

	apiHandler.APIDeleteUserLinks(w, r)

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