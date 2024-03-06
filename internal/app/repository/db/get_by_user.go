package db

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/entity"
)

// GetByUser returns all entities created by the user.
func (db *dbStore) GetByUser(ctx context.Context, userID string) ([]*entity.ShortLink, error) {
	rows, err := db.store.QueryContext(
		ctx,
		"SELECT * FROM links WHERE user_id = $1;",
		userID,
	)
	if err != nil {
		return []*entity.ShortLink{}, err
	}

	defer rows.Close()

	var list []*entity.ShortLink

	for rows.Next() {
		var shortLink entity.ShortLink

		if err := rows.Scan(
			&shortLink.ID,
			&shortLink.Short,
			&shortLink.Original,
			&shortLink.User,
			&shortLink.IsDeleted,
		); err != nil {
			return []*entity.ShortLink{}, err
		}

		list = append(list, &shortLink)
	}

	if err := rows.Err(); err != nil {
		return []*entity.ShortLink{}, err
	}

	return list, nil
}
