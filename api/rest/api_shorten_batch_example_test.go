package restapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"

	restapi "github.com/Galish/url-shortener/api/rest"
	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func ExampleHTTPHandler_APIShortenBatch() {
	bodyRaw, err := json.Marshal([]restapi.APIBatchEntity{
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

	apiHandler := restapi.NewHandler(uc, nil)

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
