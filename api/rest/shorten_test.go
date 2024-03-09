package restapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func TestShorten(t *testing.T) {
	baseURL := "http://localhost:8080"

	uc := usecase.New(memstore.New())
	defer uc.Close()

	handler := NewHandler(
		&config.Config{BaseURL: baseURL},
		uc,
		nil,
	)

	ts := httptest.NewServer(
		NewRouter(handler),
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
				"no URL provided\n",
			},
		},
		{
			"invalid request method",
			http.MethodGet,
			"/",
			"https://practicum.yandex.ru/",
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

	usecase := usecase.New(memstore.New())
	defer usecase.Close()

	handler := NewHandler(&config.Config{}, usecase, nil)

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
