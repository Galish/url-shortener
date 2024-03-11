package usecase

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/entity"
)

// Delete deletes URLs based on the given identifiers.
func (uc *ShortenerUseCase) Delete(ctx context.Context, urls ...*entity.URL) error {
	if len(urls) == 0 {
		return ErrMissingURL
	}

	go func() {
		for _, url := range urls {
			uc.messageCh <- &Message{
				action: "delete",
				url:    url,
			}
		}
	}()

	return nil
}
