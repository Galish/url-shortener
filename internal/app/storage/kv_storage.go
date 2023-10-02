package storage

import "errors"

type KeyValueStorage map[string]string

func NewKeyValueStorage() KeyValueStorage {
	return make(KeyValueStorage)
}

func (kvStorage KeyValueStorage) Get(key string) (string, error) {
	value, ok := kvStorage[key]

	if !ok {
		return "", errors.New("record doesn't not exist")
	}

	return value, nil
}

func (kvStorage KeyValueStorage) Set(key, value string) {
	kvStorage[key] = value
}

func (kvStorage KeyValueStorage) Has(key string) bool {
	_, ok := kvStorage[key]

	return ok
}
