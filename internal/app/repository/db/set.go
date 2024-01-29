package db

import (
	"context"

	repoErr "github.com/Galish/url-shortener/internal/app/repository/errors"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

// Set inserts a new entity or returns a conflict error if one exists.
func (db *dbStore) Set(ctx context.Context, shortLink *model.ShortLink) error {
	row := db.store.QueryRowContext(
		ctx,
		`
			INSERT INTO links (short_url, original_url, user_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (original_url)
			DO UPDATE SET original_url=excluded.original_url
			RETURNING short_url
		`,
		shortLink.Short,
		shortLink.Original,
		shortLink.User,
	)

	var shortURL string
	if err := row.Scan(&shortURL); err != nil {
		return err
	}

	if shortURL != shortLink.Short {
		return repoErr.New(
			repoErr.ErrConflict,
			shortURL,
			shortLink.Original,
		)
	}

	return nil
}
