package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/Galish/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeShortLink(t *testing.T) {
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			"empty URL",
			"",
			want{
				400,
				"link not provided\n",
			},
		},
		{
			"valid URL",
			"https://practicum.yandex.ru/",
			want{
				201,
				"",
			},
		},
	}

	store := storage.NewKeyValueStorage()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			h := &httpHandler{store}
			w := httptest.NewRecorder()

			h.makeShortLink(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			raw, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if result.StatusCode < 300 {
				assert.Regexp(
					t,
					regexp.MustCompile("http://localhost:8080/[0-9A-Za-z]{8}"),
					string(raw),
				)
			} else {
				assert.Equal(t, tt.want.body, string(raw))
			}
		})
	}
}
