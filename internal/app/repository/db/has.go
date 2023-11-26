package db

import "context"

func (db *dbStore) Has(ctx context.Context, key string) bool {
	row := db.store.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM links WHERE short_url = $1);", key,
	)

	var value bool
	if err := row.Scan(&value); err != nil {
		return false
	}

	return value
}
