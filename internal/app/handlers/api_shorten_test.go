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
	"github.com/Galish/url-shortener/internal/app/repository/kvstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIShorten(t *testing.T) {
	baseURL := "http://localhost:8080"
	ts := httptest.NewServer(
		NewRouter(
			&config.Config{BaseURL: baseURL},
			kvstore.New(),
			nil,
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
		req    apiRequest
		want   want
	}{
		{
			"invalid API endpoint",
			http.MethodPost,
			"/api/shortener",
			apiRequest{
				URL: "https://practicum.yandex.ru/",
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
			apiRequest{
				URL: "https://practicum.yandex.ru/",
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
			apiRequest{},
			want{
				400,
				"link not provided\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"valid URL",
			http.MethodPost,
			"/api/shorten",
			apiRequest{
				URL: "https://practicum.yandex.ru/",
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
				var respBody apiResponse
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
