package db

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/entity"
)

// Get returns the entity for a given short URL.
func (db *dbStore) Get(ctx context.Context, key string) (*entity.URL, error) {
	row := db.store.QueryRowContext(
		ctx,
		"SELECT * FROM links WHERE short_url = $1;", key,
	)

	var url entity.URL

	if err := row.Scan(
		&url.ID,
		&url.Short,
		&url.Original,
		&url.User,
		&url.IsDeleted,
	); err != nil {
		return nil, err
	}

	return &url, nil
}
