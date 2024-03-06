package restapi

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/middleware"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

func TestAPIGetUserLinks(t *testing.T) {
	store := memstore.New()
	defer store.Close()

	store.Set(context.Background(), &model.ShortLink{
		ID:       "#123111",
		Short:    "qw21dfasf",
		Original: "https://practicum.yandex.ru/",
		User:     "e44d9088-1bd6-44dc-af86-f1a551b02db3",
	})

	baseURL := "http://localhost:8080"

	handler := NewHandler(
		&config.Config{BaseURL: baseURL},
		store,
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
		token  string
		want   want
	}{
		{
			"invalid API endpoint",
			http.MethodGet,
			"/api/user-urls",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJjYTUxM2ZmNy0yMDEwLTQzOTctYWExYS0wNjU4MjhiNGJhMGUifQ.BHuk4u8KXMSEKSXTdI3_DOorpDKaapZzuSZQCSkFX9o",
			want{
				http.StatusNotFound,
				"404 page not found\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"invalid request method",
			http.MethodPost,
			"/api/user/urls",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJjYTUxM2ZmNy0yMDEwLTQzOTctYWExYS0wNjU4MjhiNGJhMGUifQ.BHuk4u8KXMSEKSXTdI3_DOorpDKaapZzuSZQCSkFX9o",
			want{
				http.StatusMethodNotAllowed,
				"",
				"",
			},
		},
		{
			"unauthorized",
			http.MethodGet,
			"/api/user/urls",
			"",
			want{
				http.StatusUnauthorized,
				"",
				"",
			},
		},
		{
			"user has no links",
			http.MethodGet,
			"/api/user/urls",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJjYTUxM2ZmNy0yMDEwLTQzOTctYWExYS0wNjU4MjhiNGJhMGUifQ.BHuk4u8KXMSEKSXTdI3_DOorpDKaapZzuSZQCSkFX9o",
			want{
				http.StatusNoContent,
				"",
				"",
			},
		},
		{
			"user has links",
			http.MethodGet,
			"/api/user/urls",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJlNDRkOTA4OC0xYmQ2LTQ0ZGMtYWY4Ni1mMWE1NTFiMDJkYjMifQ.e8r2pJHwwLWRKXyWR6Rk3MYpNkyV2LIGAthqGIheyUU",
			want{
				http.StatusOK,
				"[{\"original_url\":\"https://practicum.yandex.ru/\",\"short_url\":\"http://localhost:8080/qw21dfasf\"}]\n",
				"application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(
				tt.method,
				ts.URL+tt.path,
				nil,
			)
			require.NoError(t, err)

			if tt.token != "" {
				req.AddCookie(&http.Cookie{
					Name:  middleware.AuthCookieName,
					Value: tt.token,
				})
			}

			// Disable compression
			req.Header.Set("Accept-Encoding", "identity")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)

			raw, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.body, string(raw))

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func BenchmarkAPIGetUserLinks(b *testing.B) {
	r, _ := http.NewRequest(
		http.MethodGet,
		"/api/user/urls",
		nil,
	)

	rEmpty := r.Clone(context.Background())

	r.Header.Add(middleware.AuthHeaderName, "e44d9088-1bd6-44dc-af86-f1a551b02db3")

	w := httptest.NewRecorder()

	store := memstore.New()
	defer store.Close()

	store.Set(context.Background(), &model.ShortLink{
		ID:       "#123111",
		Short:    "qw21dfasf",
		Original: "https://practicum.yandex.ru/",
		User:     "e44d9088-1bd6-44dc-af86-f1a551b02db3",
	})

	handler := NewHandler(&config.Config{}, store)
	defer handler.Close()

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			handler.APIGetUserLinks(w, rEmpty)
		}
	})

	b.Run("valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			handler.APIGetUserLinks(w, r)
		}
	})
}
