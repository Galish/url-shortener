package db

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/entity"
)

// GetByUser returns all entities created by the user.
func (db *dbStore) GetByUser(ctx context.Context, userID string) ([]*entity.URL, error) {
	rows, err := db.store.QueryContext(
		ctx,
		"SELECT * FROM links WHERE user_id = $1;",
		userID,
	)
	if err != nil {
		return []*entity.URL{}, err
	}

	defer rows.Close()

	var list []*entity.URL

	for rows.Next() {
		var url entity.URL

		if err := rows.Scan(
			&url.ID,
			&url.Short,
			&url.Original,
			&url.User,
			&url.IsDeleted,
		); err != nil {
			return []*entity.URL{}, err
		}

		list = append(list, &url)
	}

	if err := rows.Err(); err != nil {
		return []*entity.URL{}, err
	}

	return list, nil
}
