package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/handlers"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
)

func TestShorten(t *testing.T) {
	baseURL := "http://localhost:8080"
	ts := httptest.NewServer(
		handlers.NewRouter(
			handlers.NewHandler(
				&config.Config{BaseURL: baseURL},
				memstore.New(),
			),
		),
	)
	defer ts.Close()

	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name   string
		method string
		path   string
		body   string
		want   want
	}{
		{
			"empty request body",
			http.MethodPost,
			"/",
			"",
			want{
				400,
				"link not provided\n",
			},
		},
		{
			"invalid request method",
			http.MethodGet,
			"/",
			"",
			want{
				405,
				"",
			},
		},
		{
			"valid URL",
			http.MethodPost,
			"/",
			"https://practicum.yandex.ru/",
			want{
				201,
				"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(
				tt.method,
				ts.URL+tt.path,
				strings.NewReader(tt.body),
			)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			raw, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			err = resp.Body.Close()
			require.NoError(t, err)

			if resp.StatusCode < 300 {
				assert.Regexp(
					t,
					regexp.MustCompile(baseURL+"/[0-9A-Za-z]{8}"),
					string(raw),
				)
			} else {
				assert.Equal(t, tt.want.body, string(raw))
			}
		})
	}
}

func BenchmarkShorten(b *testing.B) {
	r, _ := http.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader("qwewqewqe"),
	)

	rEmpty, _ := http.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(""),
	)

	w := httptest.NewRecorder()

	handler := handlers.NewHandler(&config.Config{}, memstore.New())

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			handler.Shorten(w, rEmpty)
		}
	})

	b.Run("valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			handler.Shorten(w, r)
		}
	})
}

func ExampleHttpHandler_Shorten() {
	apiHandler := handlers.NewHandler(
		&config.Config{BaseURL: "http://www.shortener.io"},
		memstore.New(),
	)

	router := handlers.NewRouter(apiHandler)
	server := httptest.NewServer(router)

	resp, _ := http.Post(
		server.URL,
		"text/plain",
		strings.NewReader("https://practicum.yandex.ru/"),
	)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	re := regexp.MustCompile("[A-Za-z0-9]+$")

	fmt.Println(re.ReplaceAllString(string(body), "xxxxxx"))

	// Output:
	// 201
	// text/html
	// http://www.shortener.io/xxxxxx
}
