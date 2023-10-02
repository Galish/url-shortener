package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Galish/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFullLink(t *testing.T) {
	type want struct {
		statusCode int
		location   string
		body       string
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			"base URL path",
			"/",
			want{
				400,
				"",
				"invalid identifier\n",
			},
		},
		{
			"missing entry",
			"/abKs232d",
			want{
				400,
				"",
				"record doesn't not exist\n",
			},
		},
		{
			"existing entry",
			"/c2WD8F2q",
			want{
				307,
				"https://practicum.yandex.ru/",
				"<a href=\"https://practicum.yandex.ru/\">Temporary Redirect</a>.\n\n",
			},
		},
	}

	store := storage.NewKeyValueStorage()
	store.Set("c2WD8F2q", "https://practicum.yandex.ru/")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			h := &httpHandler{store}
			w := httptest.NewRecorder()

			h.getFullLink(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))

			raw, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, string(raw), tt.want.body)
		})
	}
}
