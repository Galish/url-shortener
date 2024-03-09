package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/Galish/url-shortener/internal/app/entity"
)

// SetBatch inserts new entities in batches.
func (db *dbStore) SetBatch(ctx context.Context, urls ...*entity.URL) error {
	conn, err := db.store.Conn(context.Background())
	if err != nil {
		return err
	}

	err = conn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		var rows [][]interface{}
		for _, url := range urls {
			rows = append(
				rows,
				[]interface{}{url.Short, url.Original, url.User},
			)
		}

		_, err = conn.CopyFrom(
			ctx,
			pgx.Identifier{"links"},
			[]string{"short_url", "original_url", "user_id"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
