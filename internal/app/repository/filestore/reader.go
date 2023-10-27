package filestore

import (
	"bufio"
	"encoding/json"
	"os"
)

func (fs *fileStore) restore() error {
	file, err := os.OpenFile(fs.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := scanner.Bytes()
		var rec record
		if err := json.Unmarshal(data, &rec); err != nil {
			return err
		}

		fs.store.Set(rec.ShortUrl, rec.OriginalUrl)
		fs.size++
	}

	return nil
}
