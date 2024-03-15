package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/Galish/url-shortener/internal/app/entity"
)

// Delete marks the entity as deleted.
func (db *dbStore) Delete(ctx context.Context, urls ...*entity.URL) error {
	updateQuery := sq.Update("links").
		Set("is_deleted", true).
		PlaceholderFormat(sq.Dollar)

	where := sq.Or{}

	for _, url := range urls {
		where = append(
			where,
			sq.Eq{
				"links.short_url":  url.Short,
				"links.user_id":    url.User,
				"links.is_deleted": false,
			})
	}

	sqlStr, params, err := updateQuery.Where(where).ToSql()
	if err != nil {
		return err
	}

	_, err = db.store.ExecContext(ctx, sqlStr, params...)
	if err != nil {
		return err
	}

	return nil
}
