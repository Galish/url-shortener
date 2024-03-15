package usecase_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/entity"
	repoErr "github.com/Galish/url-shortener/internal/app/repository/errors"
	"github.com/Galish/url-shortener/internal/app/repository/mocks"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func TestShorten(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	m.EXPECT().
		Has(
			gomock.Any(),
			gomock.Any(),
		).
		Return(false).
		AnyTimes()

	m.EXPECT().
		Set(
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(func(_ context.Context, url *entity.URL) error {
			switch url.Original {
			case "https://yandex.ru/":
				return repoErr.New(
					repoErr.ErrConflict,
					url.Short,
					url.Original,
				)

			default:
				return nil
			}
		}).
		AnyTimes()

	m.EXPECT().
		SetBatch(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil).
		AnyTimes()

	uc := usecase.New(&config.Config{BaseURL: "http://shortener.io"}, m)
	defer uc.Close()

	tests := []struct {
		name string
		urls []*entity.URL
		err  error
	}{
		{
			name: "no URL provided",
			urls: []*entity.URL{},
			err:  usecase.ErrMissingURL,
		},
		{
			name: "no original URL provided",
			urls: []*entity.URL{
				{
					User: "#12345",
				},
			},
			err: usecase.ErrMissingURL,
		},
		{
			name: "single URL",
			urls: []*entity.URL{
				{
					Original: "https://practicum.yandex.ru/",
					User:     "#12345",
				},
			},
			err: nil,
		},
		{
			name: "existing URL",
			urls: []*entity.URL{
				{
					Original: "https://yandex.ru/",
					User:     "#12345",
				},
			},
			err: usecase.ErrConflict,
		},
		{
			name: "multiple URLs",
			urls: []*entity.URL{
				{
					Original: "https://practicum.yandex.ru/#11111",
					User:     "#12345",
				},
				{
					Original: "https://practicum.yandex.ru/#22222",
					User:     "#12345",
				},
				{
					Original: "https://practicum.yandex.ru/#333333",
					User:     "#12345",
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls := tt.urls

			err := uc.Shorten(context.Background(), urls...)
			assert.Equal(t, tt.err, err)

			if tt.err != nil {
				return
			}

			for _, u := range urls {
				assert.Regexp(
					t,
					regexp.MustCompile("http://shortener.io/[0-9A-Za-z]{8}"),
					uc.ShortURL(u),
				)
			}
		})
	}
}
