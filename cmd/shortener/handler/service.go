package handler

import (
	"github.com/url-shortener/cmd/shortener/storage"
)

type shortenerService struct {
	store storage.KeyValueStorage
}
