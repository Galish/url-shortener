// Package  implements the file storage.
package filestore

import (
	"bufio"
	"context"
	"os"

	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/repository/memstore"
	"github.com/Galish/url-shortener/pkg/logger"
)

type fileStore struct {
	size     int
	store    *memstore.MemStore
	filepath string
	file     *os.File
	writer   *bufio.Writer
}

// New returns a new file store instance.
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

// Get returns the entity for a given short URL.
func (fs *fileStore) Get(ctx context.Context, key string) (*entity.URL, error) {
	return fs.store.Get(ctx, key)
}

// GetByUser returns all entities created by the user.
func (fs *fileStore) GetByUser(ctx context.Context, userID string) ([]*entity.URL, error) {
	return fs.store.GetByUser(ctx, userID)
}

// Set adds a new entity to the store.
func (fs *fileStore) Set(ctx context.Context, shortLink *entity.URL) error {
	if err := fs.write(shortLink); err != nil {
		return err
	}

	if err := fs.store.Set(ctx, shortLink); err != nil {
		return err
	}

	fs.size++

	return nil
}

// SetBatch inserts new entities into the store in batches.
func (fs *fileStore) SetBatch(ctx context.Context, urls ...*entity.URL) error {
	for _, shortLink := range urls {
		fs.Set(ctx, shortLink)
	}

	return nil
}

// Delete marks the entity as deleted.
func (fs *fileStore) Delete(ctx context.Context, urls ...*entity.URL) error {
	if err := fs.store.Delete(ctx, urls...); err != nil {
		return err
	}

	for _, url := range urls {
		deleteLink, err := fs.store.Get(ctx, url.Short)
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

// Stats returns the number shortened URLs and users.
func (fs *fileStore) Stats(ctx context.Context) (int, int, error) {
	return fs.store.Stats(ctx)
}

// Has checks whether an entity with a given short URL exists.
func (fs *fileStore) Has(ctx context.Context, key string) bool {
	return fs.store.Has(ctx, key)
}

// Ping is used to make sure the store is up.
func (fs *fileStore) Ping(ctx context.Context) (bool, error) {
	return fs.filepath != "", nil
}

// Close closes the File.
func (fs *fileStore) Close() error {
	if fs.file != nil {
		return fs.file.Close()
	}

	return nil
}
