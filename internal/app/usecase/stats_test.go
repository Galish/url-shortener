package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/repository/mocks"
	"github.com/Galish/url-shortener/internal/app/usecase"
)

func TestStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	m.EXPECT().
		Stats(gomock.Any()).
		Return(10, 2, nil)

	m.EXPECT().
		Stats(gomock.Any()).
		Return(0, 0, errors.New("error occurred"))

	uc := usecase.New(&config.Config{}, m)
	defer uc.Close()

	tests := []struct {
		name  string
		urls  int
		users int
		err   error
	}{
		{
			name:  "successful execution",
			urls:  10,
			users: 2,
			err:   nil,
		},
		{
			name:  "reading from repo error",
			urls:  0,
			users: 0,
			err:   fmt.Errorf("unable to read from repository: %w", errors.New("error occurred")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls, users, err := uc.Stats(context.Background())

			assert.Equal(t, tt.urls, urls)
			assert.Equal(t, tt.users, users)
			assert.Equal(t, tt.err, err)
		})
	}
}
