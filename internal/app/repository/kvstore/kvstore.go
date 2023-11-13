package kvstore

import (
	"context"
	"errors"
)

type KVStore struct {
	store map[string]string
}

func New() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
}

func (s *KVStore) Get(ctx context.Context, key string) (string, error) {
	value, ok := s.store[key]
	if !ok {
		return "", errors.New("record doesn't not exist")
	}

	return value, nil
}

func (s *KVStore) Set(ctx context.Context, key, value string) error {
	s.store[key] = value
	return nil
}

func (s *KVStore) SetBatch(ctx context.Context, rows ...[]interface{}) error {
	for _, row := range rows {
		s.Set(ctx, row[0].(string), row[1].(string))
	}

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
