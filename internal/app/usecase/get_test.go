package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/repository/mocks"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	m.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, id string) (*entity.URL, error) {
			switch id {
			case "abcde123":
				return nil, errors.New("error occurred")

			default:
				return &entity.URL{Short: "ase21fd"}, nil
			}
		}).
		AnyTimes()

	uc := usecase.New(&config.Config{}, m)
	defer uc.Close()

	tests := []struct {
		name string
		id   string
		url  *entity.URL
		err  error
	}{
		{
			name: "invalid identifier",
			id:   "a123",
			url:  nil,
			err:  usecase.ErrInvalidID,
		},
		{
			name: "valid URL",
			id:   "abc12345",
			url:  &entity.URL{Short: "ase21fd"},
			err:  nil,
		},
		{
			name: "reading from repo error",
			id:   "abcde123",
			url:  nil,
			err:  fmt.Errorf("unable to read from repository: %w", errors.New("error occurred")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := uc.Get(context.Background(), tt.id)

			assert.Equal(t, tt.url, url)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestGetByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	m.EXPECT().
		GetByUser(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, user string) ([]*entity.URL, error) {
			switch user {
			case "#user12345":
				return []*entity.URL{
					{Short: "ase21fd"},
					{Short: "lkert423"},
					{Short: "34kf759d", IsDeleted: true},
				}, nil

			default:
				return nil, errors.New("error occurred")
			}
		}).
		AnyTimes()

	uc := usecase.New(&config.Config{}, m)
	defer uc.Close()

	tests := []struct {
		name string
		user string
		urls []*entity.URL
		err  error
	}{
		{
			name: "invalid user",
			user: "",
			urls: nil,
			err:  usecase.ErrInvalidUser,
		},
		{
			name: "valid URLs",
			user: "#user12345",
			urls: []*entity.URL{
				{Short: "ase21fd"},
				{Short: "lkert423"},
			},
			err: nil,
		},
		{
			name: "reading from repo error",
			user: "#user123",
			urls: nil,
			err:  fmt.Errorf("unable to read from repository: %w", errors.New("error occurred")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls, err := uc.GetByUser(context.Background(), tt.user)

			assert.Equal(t, tt.urls, urls)
			assert.Equal(t, tt.err, err)
		})
	}
}
