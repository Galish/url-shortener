package db

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/repository/model"
	sq "github.com/Masterminds/squirrel"
)

func (db *dbStore) Delete(ctx context.Context, shortLinks ...*model.ShortLink) error {
	updateQuery := sq.Update("links").
		Set("is_deleted", true).
		PlaceholderFormat(sq.Dollar)

	where := sq.Or{}

	for _, link := range shortLinks {
		where = append(
			where,
			sq.Eq{
				"links.short_url":  link.Short,
				"links.user_id":    link.User,
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
