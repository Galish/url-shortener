package handler_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/entity"
	restAPI "github.com/Galish/url-shortener/internal/app/infrastructure/rest"
	"github.com/Galish/url-shortener/internal/app/infrastructure/rest/handler"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func TestGet(t *testing.T) {
	store := memstore.New()

	store.SetBatch(
		context.Background(),
		&entity.URL{
			Short:    "c2WD8F2q",
			Original: "https://practicum.yandex.ru/",
		},
		&entity.URL{
			Short:     "h9h2fhfU",
			Original:  "https://practicum.yandex.ru/",
			IsDeleted: true,
		},
	)

	cfg := &config.Config{}

	uc := usecase.New(cfg, store)
	defer uc.Close()

	h := handler.New(uc, nil)

	ts := httptest.NewServer(
		restAPI.NewRouter(cfg, h),
	)
	defer ts.Close()

	type want struct {
		statusCode int
		location   string
		body       string
	}
	tests := []struct {
		name   string
		method string
		path   string
		want   want
	}{
		{
			"base URL path",
			http.MethodGet,
			"/",
			want{
				http.StatusMethodNotAllowed,
				"",
				"",
			},
		},
		{
			"invalid request method",
			http.MethodPost,
			"/abKs232d",
			want{
				http.StatusMethodNotAllowed,
				"",
				"",
			},
		},
		{
			"missing entry",
			http.MethodGet,
			"/abKs232d",
			want{
				http.StatusBadRequest,
				"",
				"unable to read from repository: record doesn't not exist\n",
			},
		},
		{
			"existing entry",
			http.MethodGet,
			"/c2WD8F2q",
			want{
				http.StatusTemporaryRedirect,
				"https://practicum.yandex.ru/",
				"<a href=\"https://practicum.yandex.ru/\">Temporary Redirect</a>.\n\n",
			},
		},
		{
			"deleted entry",
			http.MethodGet,
			"/h9h2fhfU",
			want{
				http.StatusGone,
				"",
				"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, ts.URL+tt.path, nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.location, resp.Header.Get("Location"))

			raw, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, string(raw), tt.want.body)
		})
	}
}

func BenchmarkGet(b *testing.B) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	rCtxEmpty := chi.NewRouteContext()
	rCtxEmpty.URLParams.Add("id", "Edz0Thb1")
	rEmpty := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rCtxEmpty))

	rCtxFound := chi.NewRouteContext()
	rCtxFound.URLParams.Add("id", "Edz0ThbX")
	rFound := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rCtxFound))

	w := httptest.NewRecorder()

	store := memstore.New()
	store.Set(context.Background(), &entity.URL{
		ID:       "#123111",
		Short:    "Edz0ThbX",
		Original: "https://practicum.yandex.ru/",
		User:     "e44d9088-1bd6-44dc-af86-f1a551b02db3",
	})

	uc := usecase.New(&config.Config{}, store)
	defer uc.Close()

	h := handler.New(uc, nil)

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("wrong", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h.Get(w, r)
		}
	})

	b.Run("empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h.Get(w, rEmpty)
		}
	})

	b.Run("found", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h.Get(w, rFound)
		}
	})
}