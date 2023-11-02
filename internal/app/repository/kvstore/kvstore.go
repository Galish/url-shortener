package kvstore

import "errors"

type kvStore struct {
	store map[string]string
}

func New() *kvStore {
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

func (s *kvStore) Set(key, value string) error {
	s.store[key] = value
	return nil
}

func (s *kvStore) Has(key string) bool {
	_, ok := s.store[key]

	return ok
}
