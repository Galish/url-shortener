package db

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/repository/model"
)

func (db *dbStore) Get(ctx context.Context, key string) (*model.ShortLink, error) {
	row := db.store.QueryRowContext(
		ctx,
		"SELECT * FROM links WHERE short_url = $1;", key,
	)

	var shortLink model.ShortLink

	if err := row.Scan(
		&shortLink.ID,
		&shortLink.Short,
		&shortLink.Original,
		&shortLink.User,
		&shortLink.IsDeleted,
	); err != nil {
		return nil, err
	}

	return &shortLink, nil
}
