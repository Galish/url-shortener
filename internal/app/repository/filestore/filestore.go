package filestore

import (
	"bufio"
	"context"
	"os"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository/kvstore"
	"github.com/Galish/url-shortener/internal/app/repository/models"
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

func (fs *fileStore) Get(ctx context.Context, key string) (*models.ShortLink, error) {
	return fs.store.Get(ctx, key)
}

func (fs *fileStore) GetByUser(ctx context.Context, userID string) ([]*models.ShortLink, error) {
	// TODO: implement
	return []*models.ShortLink{}, nil
}

func (fs *fileStore) Set(ctx context.Context, shortLink *models.ShortLink) error {
	if err := fs.write(shortLink); err != nil {
		return err
	}

	if err := fs.store.Set(ctx, shortLink); err != nil {
		return err
	}

	fs.size++

	return nil
}

func (fs *fileStore) SetBatch(ctx context.Context, shortLinks ...*models.ShortLink) error {
	for _, shortLink := range shortLinks {
		fs.Set(ctx, shortLink)
	}

	return nil
}

func (fs *fileStore) Delete(ctx context.Context, shortLinks ...*models.ShortLink) error {
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
