package config

import (
	"encoding/json"
	"io"
	"os"
)

func parseFile(s *settings, fileName ...string) error {
	var name string

	for _, n := range fileName {
		if n != "" {
			name = n
			break
		}
	}

	if name == "" {
		return nil
	}

	f, err := os.Open(name)
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	return nil
}
