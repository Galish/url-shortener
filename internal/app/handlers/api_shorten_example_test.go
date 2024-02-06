package handlers_test

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
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
)

func ExampleHTTPHandler_APIShorten() {
	bodyRaw, err := json.Marshal(handlers.APIRequest{
		URL: "https://practicum.yandex.ru/",
	})
	if err != nil {
		log.Fatal(err)
	}

	r, _ := http.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBuffer(bodyRaw),
	)

	w := httptest.NewRecorder()

	apiHandler := handlers.NewHandler(
		&config.Config{BaseURL: "http://www.shortener.io"},
		memstore.New(),
	)
	defer apiHandler.Close()

	apiHandler.APIShorten(w, r)

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
	// {"result":"http://www.shortener.io/xxxxxx"}
}
