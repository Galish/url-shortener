package usecase

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/pkg/generator"
)

// Shorten generates short URLs and saves them to the repo.
func (uc *ShortenerUseCase) Shorten(ctx context.Context, urls ...*entity.URL) error {
	if len(urls) == 0 {
		return ErrMissingURL
	}

	for _, url := range urls {
		if url.Original == "" {
			return ErrMissingURL
		}

		var short string

		for {
			short = generator.NewID(8)

			if !uc.repo.Has(ctx, short) {
				break
			}
		}

		url.Short = short
	}

	if len(urls) == 1 {
		return uc.repo.Set(ctx, urls[0])
	}

	return uc.repo.SetBatch(ctx, urls...)
}
