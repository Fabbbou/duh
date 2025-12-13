package toml_repo

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml"
)

type TomlDriver[T any] struct {
	filePath string
}

func (d *TomlDriver[T]) Load() (*T, error) {
	if len(d.filePath) <= 0 {
		return nil, errors.New("could not load storage, file path is empty")
	}

	tree, err := toml.LoadFile(d.filePath)
	if err != nil {
		return nil, err
	}

	var typeToLoad T
	if err := tree.Unmarshal(&typeToLoad); err != nil {
		return nil, err
	}
	return &typeToLoad, nil
}

func (d *TomlDriver[T]) Save(newVersion T) error {
	//save a new version of the storage in the file
	bytes, err := toml.Marshal(newVersion)
	if err != nil {
		return err
	}
	return os.WriteFile(d.filePath, bytes, 0755)
}
