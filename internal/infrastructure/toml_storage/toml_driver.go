package toml_storage

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml"
)

type Storage struct {
	Aliases map[string]string `toml:"aliases"`
	Exports map[string]string `toml:"exports"`
}

type TomlDriver struct {
	filePath string
}

func (d *TomlDriver) Load() (*Storage, error) {
	if len(d.filePath) <= 0 {
		return nil, errors.New("could not load storage, file path is empty")
	}

	tree, err := toml.LoadFile(d.filePath)
	if err != nil {
		return nil, err
	}
	var cfg Storage
	if err := tree.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (d *TomlDriver) Save(newVersion *Storage) error {
	//save a new version of the storage in the file
	bytes, err := toml.Marshal(*newVersion)
	if err != nil {
		return err
	}
	return os.WriteFile(d.filePath, bytes, 0755)
}
