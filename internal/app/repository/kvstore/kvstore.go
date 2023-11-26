package kvstore

import (
	"context"
	"errors"

	"github.com/Galish/url-shortener/internal/app/repository/models"
)

type KVStore struct {
	store map[string]string
}

func New() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
}

func (s *KVStore) Get(ctx context.Context, key string) (*models.ShortLink, error) {
	value, ok := s.store[key]
	if !ok {
		return nil, errors.New("record doesn't not exist")
	}

	return &models.ShortLink{
		Short:    key,
		Original: value,
	}, nil
}

func (s *KVStore) GetByUser(ctx context.Context, userID string) ([]*models.ShortLink, error) {
	// TODO: implement
	return []*models.ShortLink{}, nil
}

func (s *KVStore) Set(ctx context.Context, shortLink *models.ShortLink) error {
	s.store[shortLink.Short] = shortLink.Original
	return nil
}

func (s *KVStore) SetBatch(ctx context.Context, shortLinks ...*models.ShortLink) error {
	for _, shortLink := range shortLinks {
		s.Set(ctx, shortLink)
	}

	return nil
}

func (s *KVStore) Delete(ctx context.Context, shortLinks ...*models.ShortLink) error {
	return nil
}

func (s *KVStore) Has(ctx context.Context, key string) bool {
	_, ok := s.store[key]

	return ok
}

func (s *KVStore) Ping(ctx context.Context) (bool, error) {
	return true, nil
}

func (s *KVStore) Close() error {
	return nil
}
