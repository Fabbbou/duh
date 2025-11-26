package toml_storage

import (
	"github.com/pelletier/go-toml"
)

type Storage struct {
	Aliases map[string]string `toml:"aliases"`
	Exports map[string]string `toml:"exports"`
}

func loadFile(path string) (*Storage, error) {
	tree, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Storage
	if err := tree.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
