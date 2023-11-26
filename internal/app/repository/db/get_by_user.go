package db

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/repository/model"
)

func (db *dbStore) GetByUser(ctx context.Context, userID string) ([]*model.ShortLink, error) {
	rows, err := db.store.QueryContext(
		ctx,
		"SELECT * FROM links WHERE user_id = $1;",
		userID,
	)
	if err != nil {
		return []*model.ShortLink{}, err
	}

	defer rows.Close()

	var list []*model.ShortLink

	for rows.Next() {
		var shortLink model.ShortLink

		if err := rows.Scan(
			&shortLink.ID,
			&shortLink.Short,
			&shortLink.Original,
			&shortLink.User,
			&shortLink.IsDeleted,
		); err != nil {
			return []*model.ShortLink{}, err
		}

		list = append(list, &shortLink)
	}

	if err := rows.Err(); err != nil {
		return []*model.ShortLink{}, err
	}

	return list, nil
}
