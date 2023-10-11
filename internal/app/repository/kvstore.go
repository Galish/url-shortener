package repository

import "errors"

type kvStore struct {
	store map[string]string
}

func newKVStore() *kvStore {
	return &kvStore{
		store: make(map[string]string),
	}
}

func (s *kvStore) Get(key string) (string, error) {
	value, ok := s.store[key]
	if !ok {
		return "", errors.New("record doesn't not exist")
	}

	return value, nil
}

func (s *kvStore) Set(key, value string) {
	s.store[key] = value
}

func (s *kvStore) Has(key string) bool {
	_, ok := s.store[key]

	return ok
}
