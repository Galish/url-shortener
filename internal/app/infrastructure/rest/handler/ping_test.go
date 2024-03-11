package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Galish/url-shortener/internal/app/config"
	restAPI "github.com/Galish/url-shortener/internal/app/infrastructure/rest"
	"github.com/Galish/url-shortener/internal/app/infrastructure/rest/handler"
	"github.com/Galish/url-shortener/internal/app/repository/mocks"
)

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	m.EXPECT().
		Ping(gomock.Any()).
		Return(true, nil).
		Times(1)

	m.EXPECT().
		Ping(gomock.Any()).
		Return(false, errors.New("error occurred")).
		Times(1)

	baseURL := "http://localhost:8080"

	cfg := &config.Config{BaseURL: baseURL}

	h := handler.New(nil, m)

	ts := httptest.NewServer(
		restAPI.NewRouter(cfg, h),
	)
	defer ts.Close()

	type want struct {
		statusCode int
		err        error
	}
	tests := []struct {
		name   string
		method string
		path   string
		want   want
	}{
		{
			"invalid request method",
			http.MethodPost,
			"/ping",
			want{
				http.StatusMethodNotAllowed,
				nil,
			},
		},
		{
			"positive response",
			http.MethodGet,
			"/ping",
			want{
				http.StatusOK,
				nil,
			},
		},
		{
			"negative response",
			http.MethodGet,
			"/ping",
			want{
				http.StatusInternalServerError,
				nil,
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

			client := &http.Client{}
			resp, err := client.Do(req)

			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		})
	}
}
