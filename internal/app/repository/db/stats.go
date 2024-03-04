package db

import (
	"context"
)

// Stats returns the number shortened URLs and users.
func (db *dbStore) Stats(ctx context.Context) (int, int, error) {
	var urls, users int

	tx, err := db.store.BeginTx(ctx, nil)
	if err != nil {
		return urls, users, err
	}

	row := tx.QueryRowContext(
		ctx,
		`
			SELECT COUNT(*)
			FROM links
			WHERE is_deleted = false
		`,
	)
	if err := row.Scan(&urls); err != nil {
		tx.Rollback()
		return urls, users, err
	}

	row = tx.QueryRowContext(
		ctx,
		`
			SELECT COUNT(DISTINCT user_id)
			FROM links
			WHERE user_id IS NOT NULL;
		`,
	)

	if err := row.Scan(&users); err != nil {
		tx.Rollback()
		return urls, users, err
	}

	return urls, users, tx.Commit()
}
