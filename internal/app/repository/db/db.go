package db

import (
	"database/sql"
	"errors"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository"
	_ "github.com/jackc/pgx/v5/stdlib"
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

func (db *dbStore) Bootstrap() error {
	tx, err := db.store.Begin()
	if err != nil {
		return err
	}

	logger.Info("database initialization")

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS links (
			id serial PRIMARY KEY,
			short_url char(8) NOT NULL,
			original_url varchar(250) NOT NULL
		)
	`)
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

func (db *dbStore) Get(key string) (string, error) {
	row := db.store.QueryRow(
		"SELECT original_url FROM links WHERE short_url = $1;", key,
	)

	var originalLink string
	if err := row.Scan(&originalLink); err != nil {
		return "", err
	}

	return originalLink, nil
}

func (db *dbStore) Set(key, value string) error {
	row := db.store.QueryRow(
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
		return repository.NewRepoError(
			repository.ErrConflict,
			shortURL,
			value,
		)
	}

	return nil
}

func (db *dbStore) SetBatch(entries ...[2]string) error {
	tx, err := db.store.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		"INSERT INTO links (short_url, original_url) VALUES ($1, $2)",
	)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		_, err := stmt.Exec(
			entry[0],
			entry[1],
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (db *dbStore) Has(key string) bool {
	row := db.store.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM links WHERE short_url = $1);", key,
	)

	var value bool
	if err := row.Scan(&value); err != nil {
		return false
	}

	return value
}

func (db *dbStore) Ping() (bool, error) {
	if err := db.store.Ping(); err != nil {
		return false, err
	}

	return true, nil
}

func (db *dbStore) Close() error {
	return db.store.Close()
}
