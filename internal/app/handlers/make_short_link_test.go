package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeShortLink(t *testing.T) {
	baseUrl := "http://localhost:8080"
	ts := httptest.NewServer(
		NewRouter(
			config.Config{BaseURL: baseUrl},
			storage.NewKeyValueStorage(),
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
					regexp.MustCompile(baseUrl+"/[0-9A-Za-z]{8}"),
					string(raw),
				)
			} else {
				assert.Equal(t, tt.want.body, string(raw))
			}
		})
	}
}
