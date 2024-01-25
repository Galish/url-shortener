package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/repository/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFullLink(t *testing.T) {
	repo := memstore.New()

	repo.SetBatch(
		context.Background(),
		&model.ShortLink{
			Short:    "c2WD8F2q",
			Original: "https://practicum.yandex.ru/",
		},
		&model.ShortLink{
			Short:     "h9h2fhfU",
			Original:  "https://practicum.yandex.ru/",
			IsDeleted: true,
		},
	)

	ts := httptest.NewServer(
		NewRouter(
			NewHandler(
				&config.Config{},
				repo,
			),
		),
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
				"record doesn't not exist\n",
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
