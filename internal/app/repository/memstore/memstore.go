// Package implements in-memory storage.
package memstore

import (
	"context"
	"errors"

	"github.com/Galish/url-shortener/internal/app/entity"
	repoErr "github.com/Galish/url-shortener/internal/app/repository/errors"
)

// MemStore represents in-memory storage.
type MemStore struct {
	store map[string]*entity.URL
}

// New returns a new in-memory storage instance.
func New() *MemStore {
	return &MemStore{
		store: make(map[string]*entity.URL),
	}
}

// Get returns the entity for a given short URL.
func (ms *MemStore) Get(ctx context.Context, shortURL string) (*entity.URL, error) {
	shortLink, ok := ms.store[shortURL]
	if ok {
		return shortLink, nil
	}

	return nil, errors.New("record doesn't not exist")
}

// GetByUser returns all entities created by the user.
func (ms *MemStore) GetByUser(ctx context.Context, userID string) ([]*entity.URL, error) {
	var userURLs []*entity.URL

	for _, shortLink := range ms.store {
		if shortLink.User == userID {
			userURLs = append(userURLs, shortLink)
		}
	}

	return userURLs, nil
}

// Set adds a new entity to the store.
func (ms *MemStore) Set(ctx context.Context, url *entity.URL) error {
	for _, u := range ms.store {
		if u.Original == url.Original {
			return repoErr.New(
				repoErr.ErrConflict,
				u.Short,
				u.Original,
			)
		}
	}

	ms.store[url.Short] = url

	return nil
}

// SetBatch inserts new entities into the store in batches.
func (ms *MemStore) SetBatch(ctx context.Context, urls ...*entity.URL) error {
	for _, url := range urls {
		ms.Set(ctx, url)
	}

	return nil
}

// Delete marks the entity as deleted.
func (ms *MemStore) Delete(ctx context.Context, urls ...*entity.URL) error {
	for _, url := range urls {
		delete, ok := ms.store[url.Short]
		if !ok {
			continue
		}

		if url.User == delete.User {
			ms.store[url.Short].IsDeleted = true
		}
	}

	return nil
}

// Stats returns the number shortened URLs and users.
func (ms *MemStore) Stats(ctx context.Context) (int, int, error) {
	var urls int

	users := make(map[string]bool)

	for _, v := range ms.store {
		if v.IsDeleted {
			continue
		}

		if v.User != "" {
			users[v.User] = true
		}

		urls++
	}

	return urls, len(users), nil
}

// Has checks whether an entity with a given short URL exists.
func (ms *MemStore) Has(ctx context.Context, shortURL string) bool {
	_, ok := ms.store[shortURL]

	return ok
}

// Ping is used to make sure the store is up.
func (ms *MemStore) Ping(ctx context.Context) (bool, error) {
	return true, nil
}

// Close does nothing in the context of in-memory storage.
func (ms *MemStore) Close() error {
	return nil
}
