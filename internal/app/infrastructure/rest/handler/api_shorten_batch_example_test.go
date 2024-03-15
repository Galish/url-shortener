package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/infrastructure/rest/handler"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func ExampleHandler_APIShortenBatch() {
	bodyRaw, err := json.Marshal([]handler.APIBatchEntity{
		{
			CorrelationID: "#12345",
			OriginalURL:   "https://practicum.yandex.ru/",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	r, _ := http.NewRequest(
		http.MethodPost,
		"/api/shorten/batch",
		bytes.NewBuffer(bodyRaw),
	)

	w := httptest.NewRecorder()

	uc := usecase.New(
		&config.Config{BaseURL: "http://www.shortener.io"},
		memstore.New(),
	)
	defer uc.Close()

	apiHandler := handler.New(uc, nil)

	apiHandler.APIShortenBatch(w, r)

	resp := w.Result()

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	re := regexp.MustCompile("/[A-Za-z0-9]{8}")

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(re.ReplaceAllString(string(body), "/xxxxxx"))

	// Output:
	// 201
	// application/json
	// [{"correlation_id":"#12345","short_url":"http://www.shortener.io/xxxxxx"}]
}
