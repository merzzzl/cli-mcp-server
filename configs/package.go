package configs

import (
	"embed"
	"fmt"
)

//go:embed *.json
var files embed.FS

func Read(name string) ([]byte, error) {
	data, err := files.ReadFile(name + ".json")
	if err != nil {
		return nil, fmt.Errorf("read embedded config %q: %w", name, err)
	}

	return data, nil
}
