package filestore

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository/models"
)

func (fs *fileStore) restore() error {
	if fs.filepath == "" {
		return nil
	}

	file, err := os.OpenFile(fs.filepath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := scanner.Bytes()
		var shortLink models.ShortLink
		if err := json.Unmarshal(data, &shortLink); err != nil {
			return err
		}

		fs.store.Set(context.Background(), &shortLink)
		fs.size++
	}

	logger.WithFields(logger.Fields{
		"recordCount": fs.size,
	}).Info("recover data from file")

	return nil
}
