package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
)

func TestAPIShorten(t *testing.T) {
	baseURL := "http://localhost:8080"

	handler := NewHandler(
		&config.Config{BaseURL: baseURL},
		memstore.New(),
	)
	defer handler.Close()

	ts := httptest.NewServer(
		NewRouter(handler),
	)
	defer ts.Close()

	type want struct {
		statusCode  int
		body        string
		contentType string
	}

	tests := []struct {
		name   string
		method string
		path   string
		req    APIRequest
		want   want
	}{
		{
			"invalid API endpoint",
			http.MethodPost,
			"/api/shortener",
			APIRequest{
				URL: "https://practicum.yandex.ru/",
			},
			want{
				http.StatusNotFound,
				"404 page not found\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"invalid request method",
			http.MethodGet,
			"/api/shorten",
			APIRequest{
				URL: "https://practicum.yandex.ru/",
			},
			want{
				http.StatusMethodNotAllowed,
				"",
				"",
			},
		},
		{
			"empty request body",
			http.MethodPost,
			"/api/shorten",
			APIRequest{},
			want{
				http.StatusBadRequest,
				"link not provided\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"valid URL",
			http.MethodPost,
			"/api/shorten",
			APIRequest{
				URL: "https://practicum.yandex.ru/",
			},
			want{
				http.StatusCreated,
				"",
				"application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.req)
			require.NoError(t, err)

			req, err := http.NewRequest(
				tt.method,
				ts.URL+tt.path,
				bytes.NewBuffer(reqBody),
			)
			require.NoError(t, err)

			// Disable compression
			req.Header.Set("Accept-Encoding", "identity")

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			if resp.StatusCode < 300 {
				var respBody APIResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)

				assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)

				assert.Regexp(
					t,
					regexp.MustCompile(baseURL+"/[0-9A-Za-z]{8}"),
					respBody.Result,
				)
			} else {
				var raw []byte
				raw, err = io.ReadAll(resp.Body)
				require.NoError(t, err)

				assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)
				assert.Equal(t, tt.want.body, string(raw))
			}

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func BenchmarkAPIShorten(b *testing.B) {
	bodyRaw, _ := json.Marshal(APIRequest{
		URL: "https://practicum.yandex.ru/",
	})

	r, _ := http.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBuffer(bodyRaw),
	)

	bodyEmptyRaw, _ := json.Marshal(APIRequest{
		URL: "",
	})

	rEmpty, _ := http.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBuffer(bodyEmptyRaw),
	)

	w := httptest.NewRecorder()

	handler := NewHandler(&config.Config{}, memstore.New())
	defer handler.Close()

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			handler.APIShorten(w, rEmpty)
		}
	})

	b.Run("valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			handler.APIShorten(w, r)
		}
	})
}
