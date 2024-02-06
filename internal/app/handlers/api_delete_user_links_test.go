package handlers

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestAPIDeleteUserLinks(t *testing.T) {
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
		req    []string
		token  string
		want   want
	}{
		{
			"invalid API endpoint",
			http.MethodDelete,
			"/api/user-urls",
			[]string{"qw21dfasf"},
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
			[]string{"qw21dfasf"},
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJjYTUxM2ZmNy0yMDEwLTQzOTctYWExYS0wNjU4MjhiNGJhMGUifQ.BHuk4u8KXMSEKSXTdI3_DOorpDKaapZzuSZQCSkFX9o",
			want{
				http.StatusMethodNotAllowed,
				"",
				"",
			},
		},
		{
			"non-existent link",
			http.MethodDelete,
			"/api/user/urls",
			[]string{"qw21df123"},
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJjYTUxM2ZmNy0yMDEwLTQzOTctYWExYS0wNjU4MjhiNGJhMGUifQ.BHuk4u8KXMSEKSXTdI3_DOorpDKaapZzuSZQCSkFX9o",
			want{
				http.StatusAccepted,
				"",
				"",
			},
		},
		{
			"existing link",
			http.MethodDelete,
			"/api/user/urls",
			[]string{"qw21dfasf"},
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJjYTUxM2ZmNy0yMDEwLTQzOTctYWExYS0wNjU4MjhiNGJhMGUifQ.BHuk4u8KXMSEKSXTdI3_DOorpDKaapZzuSZQCSkFX9o",
			want{
				http.StatusAccepted,
				"",
				"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.req)
			require.NoError(t, err)

			req, err := http.NewRequest(
				tt.method,
				ts.URL+tt.path,
				bytes.NewBuffer(body),
			)
			require.NoError(t, err)

			req.AddCookie(&http.Cookie{
				Name:  middleware.AuthCookieName,
				Value: tt.token,
			})

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

func BenchmarkAPIDeleteUserLinks(b *testing.B) {
	bodyRaw, _ := json.Marshal([]string{"qw21df123"})

	r, _ := http.NewRequest(
		http.MethodDelete,
		"/api/user/urls",
		bytes.NewBuffer(bodyRaw),
	)

	bodyEmptyRaw, _ := json.Marshal([]string{})

	rEmpty, _ := http.NewRequest(
		http.MethodDelete,
		"/api/user/urls",
		bytes.NewBuffer(bodyEmptyRaw),
	)

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
			handler.APIDeleteUserLinks(w, rEmpty)
		}
	})

	b.Run("valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			handler.APIDeleteUserLinks(w, r)
		}
	})
}
