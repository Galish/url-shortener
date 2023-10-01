package handler

import (
	"github.com/Galish/url-shortener/cmd/shortener/storage"
)

type shortenerService struct {
	store storage.KeyValueStorage
}
