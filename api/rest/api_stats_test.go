package restapi

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func TestAPIStats(t *testing.T) {
	store := memstore.New()
	defer store.Close()

	store.Set(context.Background(), &entity.URL{
		ID:       "#123111",
		Short:    "qw21dfasf",
		Original: "https://practicum.yandex.ru/",
		User:     "e44d9088-1bd6-44dc-af86-f1a551b02db3",
	})

	store.Set(context.Background(), &entity.URL{
		ID:        "#222222",
		Short:     "asd343dgs",
		Original:  "https://www.yandex.ru/",
		User:      "e44d9088-1bd6-44dc-af86-f1a551b02db3",
		IsDeleted: true,
	})

	store.Set(context.Background(), &entity.URL{
		ID:       "#3333333",
		Short:    "lkjsdfu43",
		Original: "https://kinopoisk.ru/",
		User:     "e44d9088-1bd6-44dc-af86-f1a551b04444",
	})

	baseURL := "http://localhost:8080"

	_, ipNet, _ := net.ParseCIDR("192.168.1.0/24")

	cfg := &config.Config{
		BaseURL:       baseURL,
		TrustedSubnet: ipNet,
	}

	uc := usecase.New(cfg, store)
	defer uc.Close()

	handler := NewHandler(uc, nil)

	ts := httptest.NewServer(
		NewRouter(cfg, handler),
	)
	defer ts.Close()

	type want struct {
		statusCode  int
		body        string
		contentType string
	}

	tests := []struct {
		name     string
		method   string
		path     string
		headerIP string
		want     want
	}{
		{
			"invalid API endpoint",
			http.MethodGet,
			"/api/internal-stats",
			"192.168.1.10",
			want{
				http.StatusNotFound,
				"404 page not found\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"invalid request method",
			http.MethodPost,
			"/api/internal/stats",
			"192.168.1.10",
			want{
				http.StatusMethodNotAllowed,
				"",
				"",
			},
		},
		{
			"unauthorized",
			http.MethodGet,
			"/api/internal/stats",
			"",
			want{
				http.StatusUnauthorized,
				"",
				"",
			},
		},
		{
			"unauthorized",
			http.MethodGet,
			"/api/internal/stats",
			"192.168.123.10",
			want{
				http.StatusUnauthorized,
				"",
				"",
			},
		},
		{
			"service stats",
			http.MethodGet,
			"/api/internal/stats",
			"192.168.1.10",
			want{
				http.StatusOK,
				"{\"urls\":2,\"users\":2}\n",
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

			req.Header.Set("X-Real-IP", tt.headerIP)

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
