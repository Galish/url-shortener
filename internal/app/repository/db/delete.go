package db

import (
	"context"

	"github.com/Galish/url-shortener/internal/app/repository/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func (db *dbStore) Delete(ctx context.Context, shortLinks ...*model.ShortLink) error {
	conn, err := db.store.Conn(context.Background())
	if err != nil {
		return err
	}

	err = conn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		query := `
			UPDATE links
			SET
				is_deleted = true
			WHERE (
				short_url = $1
				AND
				user_id = $2
			)
		`

		batch := &pgx.Batch{}
		for _, link := range shortLinks {
			batch.Queue(query, link.Short, link.User)
		}

		batchResults := conn.SendBatch(
			ctx,
			batch,
		)

		return batchResults.Close()
	})
	if err != nil {
		return err
	}

	return nil
}
