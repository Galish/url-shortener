package filestore

import (
	"bufio"
	"os"

	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/Galish/url-shortener/internal/app/repository/kvstore"
)

type fileStore struct {
	size     int
	store    repository.Repository
	filepath string
	file     *os.File
	writer   *bufio.Writer
}

func New(filepath string) (*fileStore, error) {
	fs := &fileStore{
		store:    kvstore.New(),
		filepath: filepath,
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

func (fs *fileStore) Has(key string) bool {
	return fs.store.Has(key)
}

func (fs *fileStore) Close() error {
	return fs.file.Close()
}
