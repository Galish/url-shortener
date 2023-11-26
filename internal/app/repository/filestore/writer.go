package filestore

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

func (fs *fileStore) initWriter() error {
	if fs.writer != nil {
		return nil
	}

	file, err := os.OpenFile(fs.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	fs.file = file
	fs.writer = bufio.NewWriter(file)

	return nil
}

func (fs *fileStore) write(shortLink *model.ShortLink) error {
	if fs.filepath == "" {
		return nil
	}

	shortLink.ID = strconv.Itoa(fs.size)

	data, err := json.Marshal(shortLink)
	if err != nil {
		return err
	}

	if err := fs.initWriter(); err != nil {
		return err
	}

	if _, err := fs.writer.Write(data); err != nil {
		return err
	}

	if err := fs.writer.WriteByte('\n'); err != nil {
		return err
	}

	logger.Info("writing a record to a file")

	return fs.writer.Flush()
}
