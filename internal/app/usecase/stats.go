package usecase

import (
	"context"
	"fmt"
)

// GetStats returns the number of total URLs and users.
func (uc *ShortenerUseCase) Stats(ctx context.Context) (int, int, error) {
	urls, users, err := uc.repo.Stats(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to read from repository: %w", err)
	}

	return urls, users, nil
}
