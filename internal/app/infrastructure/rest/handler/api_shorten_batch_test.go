package handler_test

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
	restAPI "github.com/Galish/url-shortener/internal/app/infrastructure/rest"
	"github.com/Galish/url-shortener/internal/app/infrastructure/rest/handler"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func TestAPIShortenBatch(t *testing.T) {
	baseURL := "http://localhost:8080"

	cfg := &config.Config{BaseURL: baseURL}

	uc := usecase.New(cfg, memstore.New())
	defer uc.Close()

	h := handler.New(uc, nil)

	ts := httptest.NewServer(
		restAPI.NewRouter(cfg, h),
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
		req    []handler.APIBatchEntity
		want   want
	}{
		{
			"invalid API endpoint",
			http.MethodPost,
			"/api/shorten/batches",
			[]handler.APIBatchEntity{
				{
					CorrelationID: "#12345",
					OriginalURL:   "https://practicum.yandex.ru/",
				},
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
			"/api/shorten/batch",
			[]handler.APIBatchEntity{
				{
					CorrelationID: "#12345",
					OriginalURL:   "https://practicum.yandex.ru/",
				},
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
			"/api/shorten/batch",
			[]handler.APIBatchEntity{},
			want{
				http.StatusBadRequest,
				"empty request body\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"no URL provided",
			http.MethodPost,
			"/api/shorten/batch",
			[]handler.APIBatchEntity{
				{
					CorrelationID: "#12345",
				},
			},
			want{
				http.StatusBadRequest,
				"no URL provided\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"valid URL list",
			http.MethodPost,
			"/api/shorten/batch",
			[]handler.APIBatchEntity{
				{
					CorrelationID: "#12345",
					OriginalURL:   "https://practicum.yandex.ru/",
				},
				{
					CorrelationID: "#23456",
					OriginalURL:   "https://www.yandex.ru/",
				},
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
				var respBody []handler.APIBatchEntity
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)

				assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)

				for i, v := range respBody {
					assert.Equal(
						t,
						v.CorrelationID,
						tt.req[i].CorrelationID,
					)

					assert.Regexp(
						t,
						regexp.MustCompile(baseURL+"/[0-9A-Za-z]{8}"),
						v.ShortURL,
					)
				}
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

func BenchmarkAPIShortenBatch(b *testing.B) {
	bodyRaw, _ := json.Marshal([]handler.APIBatchEntity{
		{
			CorrelationID: "#11111",
			OriginalURL:   "https://practicum.yandex.ru/",
		},
		{
			CorrelationID: "#22222",
			OriginalURL:   "https://www.google.com/",
		},
		{
			CorrelationID: "#33333",
			OriginalURL:   "https://www.ozon.ru/",
		},
	})

	r, _ := http.NewRequest(
		http.MethodPost,
		"/api/shorten/batch",
		bytes.NewBuffer(bodyRaw),
	)

	bodyEmptyRaw, _ := json.Marshal([]handler.APIBatchEntity{})

	rEmpty, _ := http.NewRequest(
		http.MethodPost,
		"/api/shorten/batch",
		bytes.NewBuffer(bodyEmptyRaw),
	)

	w := httptest.NewRecorder()

	uc := usecase.New(&config.Config{}, memstore.New())
	defer uc.Close()

	h := handler.New(uc, nil)

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h.APIShortenBatch(w, rEmpty)
		}
	})

	b.Run("valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h.APIShortenBatch(w, r)
		}
	})
}