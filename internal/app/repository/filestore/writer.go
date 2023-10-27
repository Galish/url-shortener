package filestore

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
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

func (fs *fileStore) write(key, value string) error {
	data, err := json.Marshal(record{
		ID:          strconv.Itoa(fs.size),
		ShortUrl:    key,
		OriginalUrl: value,
	})
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

	return fs.writer.Flush()
}
