package config

import (
	"encoding/json"
	"io"
	"os"
)

func parseFile(fileName string, c *settings) error {
	if fileName == "" {
		return nil
	}

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	return nil
}
