// Package implements PostgreSQL storage.
package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Galish/url-shortener/internal/logger"
)

type dbStore struct {
	store *sql.DB
}

// New returns a new database connection instance.
func New(addr string) (*dbStore, error) {
	if addr == "" {
		return nil, errors.New("database address missing")
	}

	logger.Info("database connection")

	db, err := sql.Open("pgx", addr)
	if err != nil {
		return nil, err
	}

	return &dbStore{db}, nil
}

// Bootstrap initializes default database state.
func (db *dbStore) Bootstrap(ctx context.Context) error {
	tx, err := db.store.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	logger.Info("database initialization")

	_, err = tx.ExecContext(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS links (
				id serial PRIMARY KEY,
				short_url char(8) NOT NULL,
				original_url varchar(250) NOT NULL,
				user_id char(36),
				is_deleted boolean DEFAULT false
			)
		`,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON links (
			original_url
		)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Ping is used to make sure the database connection is alive.
func (db *dbStore) Ping(ctx context.Context) (bool, error) {
	if err := db.store.PingContext(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// Close closes a connection.
func (db *dbStore) Close() error {
	return db.store.Close()
}
