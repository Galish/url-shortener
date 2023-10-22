package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/models"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIShorten(t *testing.T) {
	baseURL := "http://localhost:8080"
	ts := httptest.NewServer(
		NewRouter(
			&config.Config{BaseURL: baseURL},
			repository.New(),
		),
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
		req    models.ApiRequest
		want   want
	}{
		{
			"invalid API endpoint",
			http.MethodPost,
			"/api/shortener",
			models.ApiRequest{
				Url: "https://practicum.yandex.ru/",
			},
			want{
				404,
				"404 page not found\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"invalid request method",
			http.MethodGet,
			"/api/shorten",
			models.ApiRequest{
				Url: "https://practicum.yandex.ru/",
			},
			want{
				405,
				"",
				"",
			},
		},
		{
			"empty request body",
			http.MethodPost,
			"/api/shorten",
			models.ApiRequest{},
			want{
				400,
				"no URL provided\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"valid URL",
			http.MethodPost,
			"/api/shorten",
			models.ApiRequest{
				Url: "https://practicum.yandex.ru/",
			},
			want{
				201,
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

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			if resp.StatusCode < 300 {
				var respBody models.ApiResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)

				assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)

				assert.Regexp(
					t,
					regexp.MustCompile(baseURL+"/[0-9A-Za-z]{8}"),
					respBody.Result,
				)
			} else {
				raw, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)
				assert.Equal(t, tt.want.body, string(raw))
			}

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}
