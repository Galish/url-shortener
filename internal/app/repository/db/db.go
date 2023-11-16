package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Galish/url-shortener/internal/app/logger"
	repoErr "github.com/Galish/url-shortener/internal/app/repository/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type dbStore struct {
	store *sql.DB
}

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
				original_url varchar(250) NOT NULL
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

func (db *dbStore) Get(ctx context.Context, key string) (string, error) {
	row := db.store.QueryRowContext(
		ctx,
		"SELECT original_url FROM links WHERE short_url = $1;", key,
	)

	var originalLink string
	if err := row.Scan(&originalLink); err != nil {
		return "", err
	}

	return originalLink, nil
}

func (db *dbStore) Set(ctx context.Context, key, value string) error {
	row := db.store.QueryRowContext(
		ctx,
		`
			INSERT INTO links (short_url, original_url)
			VALUES ($1, $2)
			ON CONFLICT (original_url)
			DO UPDATE SET original_url=excluded.original_url
			RETURNING short_url
		`,
		key,
		value,
	)

	var shortURL string
	if err := row.Scan(&shortURL); err != nil {
		return err
	}

	if shortURL != key {
		return repoErr.New(
			repoErr.ErrConflict,
			shortURL,
			value,
		)
	}

	return nil
}

func (db *dbStore) SetBatch(ctx context.Context, rows ...[]interface{}) error {
	conn, err := db.store.Conn(context.Background())
	if err != nil {
		return err
	}

	err = conn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		_, err := conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"links"},
			[]string{"short_url", "original_url"},
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

func (db *dbStore) Ping(ctx context.Context) (bool, error) {
	if err := db.store.PingContext(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (db *dbStore) Close() error {
	return db.store.Close()
}
