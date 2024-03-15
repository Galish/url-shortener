package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/repository/mocks"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	m.EXPECT().
		Delete(
			gomock.Any(),
			&entity.URL{
				Short: "acnq339s",
				User:  "#12345",
			},
			&entity.URL{
				Short: "7asmss6",
				User:  "#11111",
			},
			&entity.URL{
				Short: "su73hd",
				User:  "#22222",
			},
			&entity.URL{
				Short: "jjs730",
				User:  "#33333",
			},
		).
		Return(nil).
		Times(1)

	uc := usecase.New(&config.Config{}, m)
	defer uc.Close()

	tests := []struct {
		name string
		urls []*entity.URL
		err  error
	}{
		{
			name: "no URL provided",
			urls: nil,
			err:  usecase.ErrMissingURL,
		},
		{
			name: "single URL",
			urls: []*entity.URL{
				{
					Short: "acnq339s",
					User:  "#12345",
				},
			},
			err: nil,
		},
		{
			name: "multiple URLs",
			urls: []*entity.URL{
				{
					Short: "7asmss6",
					User:  "#11111",
				},
				{
					Short: "su73hd",
					User:  "#22222",
				},
				{
					Short: "jjs730",
					User:  "#33333",
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.Delete(context.Background(), tt.urls...)

			assert.Equal(t, tt.err, err)
		})
	}
}
