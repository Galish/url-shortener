package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type dbStore struct {
	store *sql.DB
}

func New(addr string) (*dbStore, error) {
	db, err := sql.Open("pgx", addr)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS links (
		"short_url" CHAR(8) NOT NULL,
		"original_url" VARCHAR(250) NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	return &dbStore{
		store: db,
	}, nil
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
	_, err := db.store.Exec(
		"INSERT INTO links (short_url, original_url) VALUES ($1, $2)", key, value,
	)
	if err != nil {
		return err
	}

	return nil
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
