package usecase

import (
	"context"
	"fmt"

	"github.com/Galish/url-shortener/internal/app/entity"
	repoErr "github.com/Galish/url-shortener/internal/app/repository/errors"
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

	if len(urls) != 1 {
		return uc.repo.SetBatch(ctx, urls...)
	}

	err := uc.repo.Set(ctx, urls[0])

	errConflict := repoErr.AsErrConflict(err)
	if errConflict != nil {
		urls[0].Short = errConflict.ShortURL
		return ErrConflict
	}

	return err
}

func (uc *ShortenerUseCase) ShortURL(url *entity.URL) string {
	return fmt.Sprintf("%s/%s", uc.cfg.BaseURL, url.Short)
}
