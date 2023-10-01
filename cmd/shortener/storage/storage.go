package storage

import "errors"

type KeyValueStorage interface {
	Get(string) (string, error)
	Set(string, string)
}

type Storage map[string]string

func NewKeyValueStorage() KeyValueStorage {
	return make(Storage)
}

func (s Storage) Get(key string) (string, error) {
	value, ok := s[key]
	if !ok {
		return "", errors.New("record doesn't not exist")
	}

	return value, nil
}

func (s Storage) Set(key, value string) {
	s[key] = value
}
