package filestore

import (
	"bufio"
	"encoding/json"
	"os"
)

func (fs *fileStore) restore() error {
	file, err := os.OpenFile(fs.filepath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := scanner.Bytes()
		var lnk link
		if err := json.Unmarshal(data, &lnk); err != nil {
			return err
		}

		fs.store.Set(lnk.Short, lnk.Original)
		fs.size++
	}

	return nil
}
