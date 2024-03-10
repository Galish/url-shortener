package usecase

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/entity"
)

// Delete deletes URLs based on the given identifiers.
func (uc *ShortenerUseCase) Delete(ctx context.Context, urls ...*entity.URL) {
	if len(urls) == 0 {
		return
	}

	go func() {
		for _, url := range urls {
			uc.messageCh <- &shortenerMessage{
				action: "delete",
				url:    url,
			}
		}
	}()
}
