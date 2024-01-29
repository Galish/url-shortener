package handlers_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
)

func ExampleHTTPHandler_Shorten() {
	apiHandler := handlers.NewHandler(
		&config.Config{BaseURL: "http://www.shortener.io"},
		memstore.New(),
	)

	router := handlers.NewRouter(apiHandler)
	server := httptest.NewServer(router)

	resp, err := http.Post(
		server.URL,
		"text/plain",
		strings.NewReader("https://practicum.yandex.ru/"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	re := regexp.MustCompile("[A-Za-z0-9]+$")

	fmt.Println(re.ReplaceAllString(string(body), "xxxxxx"))

	// Output:
	// 201
	// text/html
	// http://www.shortener.io/xxxxxx
}
