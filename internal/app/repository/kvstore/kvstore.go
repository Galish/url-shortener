package kvstore

import "errors"

type KVStore struct {
	store map[string]string
}

func New() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
}

func (s *KVStore) Get(key string) (string, error) {
	value, ok := s.store[key]
	if !ok {
		return "", errors.New("record doesn't not exist")
	}

	return value, nil
}

func (s *KVStore) Set(key, value string) error {
	s.store[key] = value
	return nil
}

func (s *KVStore) Has(key string) bool {
	_, ok := s.store[key]

	return ok
}

func (s *KVStore) Ping() (bool, error) {
	return true, nil
}

func (s *KVStore) Close() error {
	return nil
}
