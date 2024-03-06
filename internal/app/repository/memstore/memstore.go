// Package implements in-memory storage.
package memstore

import (
	"context"
	"errors"

	"github.com/Galish/url-shortener/internal/app/entity"
)

// MemStore represents in-memory storage.
type MemStore struct {
	store map[string]*entity.ShortLink
}

// New returns a new in-memory storage instance.
func New() *MemStore {
	return &MemStore{
		store: make(map[string]*entity.ShortLink),
	}
}

// Get returns the entity for a given short URL.
func (ms *MemStore) Get(ctx context.Context, shortURL string) (*entity.ShortLink, error) {
	shortLink, ok := ms.store[shortURL]
	if ok {
		return shortLink, nil
	}

	return nil, errors.New("record doesn't not exist")
}

// GetByUser returns all entities created by the user.
func (ms *MemStore) GetByUser(ctx context.Context, userID string) ([]*entity.ShortLink, error) {
	var userShortLinks []*entity.ShortLink

	for _, shortLink := range ms.store {
		if shortLink.User == userID {
			userShortLinks = append(userShortLinks, shortLink)
		}
	}

	return userShortLinks, nil
}

// Set adds a new entity to the store.
func (ms *MemStore) Set(ctx context.Context, shortLink *entity.ShortLink) error {
	ms.store[shortLink.Short] = shortLink
	return nil
}

// SetBatch inserts new entities into the store in batches.
func (ms *MemStore) SetBatch(ctx context.Context, shortLinks ...*entity.ShortLink) error {
	for _, shortLink := range shortLinks {
		ms.Set(ctx, shortLink)
	}

	return nil
}

// Delete marks the entity as deleted.
func (ms *MemStore) Delete(ctx context.Context, shortLinks ...*entity.ShortLink) error {
	for _, shortLink := range shortLinks {
		deleteLink, ok := ms.store[shortLink.Short]
		if !ok {
			continue
		}

		if shortLink.User == deleteLink.User {
			ms.store[shortLink.Short].IsDeleted = true
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
