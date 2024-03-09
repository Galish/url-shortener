package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Galish/url-shortener/internal/app/entity"
)

// Get returns the URL for the given identifier.
func (uc *ShortenerUseCase) Get(ctx context.Context, id string) (*entity.URL, error) {
	if len(id) < 8 {
		return nil, errors.New("invalid identifier")
	}

	url, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to read from repository: %w", err)
	}

	return url, nil
}

// GetByUser returns a list of URLs created by the user.
func (uc *ShortenerUseCase) GetByUser(ctx context.Context, user string) ([]*entity.URL, error) {
	userURLs, err := uc.repo.GetByUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("unable to read from repository: %w", err)
	}

	urls := make([]*entity.URL, 0, len(userURLs))
	for _, u := range userURLs {
		if u.IsDeleted {
			continue
		}

		urls = append(urls, u)
	}

	return urls, nil
}
