package filestore

import (
	"bufio"
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

func (fs *fileStore) Get(key string) (string, error) {
	return fs.store.Get(key)
}

func (fs *fileStore) Set(key, value string) error {
	if err := fs.write(key, value); err != nil {
		return err
	}

	if err := fs.store.Set(key, value); err != nil {
		return err
	}

	fs.size++

	return nil
}

func (fs *fileStore) SetBatch(entries ...[2]string) error {
	for _, entry := range entries {
		fs.Set(entry[0], entry[1])
	}

	return nil
}

func (fs *fileStore) Has(key string) bool {
	return fs.store.Has(key)
}

func (fs *fileStore) Ping() (bool, error) {
	return fs.filepath != "", nil
}

func (fs *fileStore) Close() error {
	if fs.file != nil {
		return fs.file.Close()
	}

	return nil
}
