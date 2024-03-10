package handler_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/infrastructure/rest/handler"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func ExampleHandler_Shorten() {
	r, _ := http.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader("https://practicum.yandex.ru/"),
	)

	w := httptest.NewRecorder()

	uc := usecase.New(
		&config.Config{BaseURL: "http://www.shortener.io"},
		memstore.New(),
	)
	defer uc.Close()

	apiHandler := handler.New(uc, nil)

	apiHandler.Shorten(w, r)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	re := regexp.MustCompile("[A-Za-z0-9]+$")

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(re.ReplaceAllString(string(body), "xxxxxx"))

	// Output:
	// 201
	// text/html
	// http://www.shortener.io/xxxxxx
}
