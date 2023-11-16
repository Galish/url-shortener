package filestore

import (
	"bufio"
	"context"
	"os"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository/kvstore"
)

type fileStore struct {
	size     int
	store    *kvstore.KVStore
	filepath string
	file     *os.File
	writer   *bufio.Writer
}

func New(filepath string) (*fileStore, error) {
	fs := &fileStore{
		store:    kvstore.New(),
		filepath: filepath,
	}

	if filepath == "" {
		logger.Debug("writing to file is disabled")
	}

	if err := fs.restore(); err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *fileStore) Get(ctx context.Context, key string) (string, error) {
	return fs.store.Get(ctx, key)
}

func (fs *fileStore) Set(ctx context.Context, key, value string) error {
	if err := fs.write(key, value); err != nil {
		return err
	}

	if err := fs.store.Set(ctx, key, value); err != nil {
		return err
	}

	fs.size++

	return nil
}

func (fs *fileStore) SetBatch(ctx context.Context, rows ...[]interface{}) error {
	for _, row := range rows {
		fs.Set(ctx, row[0].(string), row[1].(string))
	}

	return nil
}

func (fs *fileStore) Has(ctx context.Context, key string) bool {
	return fs.store.Has(ctx, key)
}

func (fs *fileStore) Ping(ctx context.Context) (bool, error) {
	return fs.filepath != "", nil
}

func (fs *fileStore) Close() error {
	if fs.file != nil {
		return fs.file.Close()
	}

	return nil
}
