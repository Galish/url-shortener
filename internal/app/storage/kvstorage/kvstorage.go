package kvstorage

import "errors"

type KeyValueStorage struct {
	store map[string]string
}

func New() *KeyValueStorage {
	return &KeyValueStorage{
		store: make(map[string]string),
	}
}

func (s *KeyValueStorage) Get(key string) (string, error) {
	value, ok := s.store[key]

	if !ok {
		return "", errors.New("record doesn't not exist")
	}

	return value, nil
}

func (s *KeyValueStorage) Set(key, value string) {
	s.store[key] = value
}

func (s *KeyValueStorage) Has(key string) bool {
	_, ok := s.store[key]

	return ok
}
