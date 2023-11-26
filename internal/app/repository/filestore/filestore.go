package filestore

import (
	"bufio"
	"context"
	"os"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

type fileStore struct {
	size     int
	store    *memstore.MemStore
	filepath string
	file     *os.File
	writer   *bufio.Writer
}

func New(filepath string) (*fileStore, error) {
	fs := &fileStore{
		store:    memstore.New(),
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

func (fs *fileStore) Get(ctx context.Context, key string) (*model.ShortLink, error) {
	return fs.store.Get(ctx, key)
}

func (fs *fileStore) GetByUser(ctx context.Context, userID string) ([]*model.ShortLink, error) {
	return fs.store.GetByUser(ctx, userID)
}

func (fs *fileStore) Set(ctx context.Context, shortLink *model.ShortLink) error {
	if err := fs.write(shortLink); err != nil {
		return err
	}

	if err := fs.store.Set(ctx, shortLink); err != nil {
		return err
	}

	fs.size++

	return nil
}

func (fs *fileStore) SetBatch(ctx context.Context, shortLinks ...*model.ShortLink) error {
	for _, shortLink := range shortLinks {
		fs.Set(ctx, shortLink)
	}

	return nil
}

func (fs *fileStore) Delete(ctx context.Context, shortLinks ...*model.ShortLink) error {
	if err := fs.store.Delete(ctx, shortLinks...); err != nil {
		return err
	}

	for _, shortLink := range shortLinks {
		deleteLink, err := fs.store.Get(ctx, shortLink.Short)
		if err != nil {
			logger.WithError(err).Debug("unable to read from store")
			continue
		}

		if !deleteLink.IsDeleted {
			continue
		}

		if err := fs.write(deleteLink); err != nil {
			logger.WithError(err).Debug("unable to write to repository")
			continue
		}
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
