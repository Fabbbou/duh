package tomll

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml"
)

func LoadToml[T any](filePath string) (*T, error) {
	if len(filePath) <= 0 {
		return nil, errors.New("could not load storage, file path is empty")
	}

	tree, err := toml.LoadFile(filePath)
	if err != nil {
		return nil, err
	}

	var typeToLoad T
	if err := tree.Unmarshal(&typeToLoad); err != nil {
		return nil, err
	}
	return &typeToLoad, nil
}

func SaveToml[T any](filePath string, newVersion T) error {
	//save a new version of the storage in the file
	bytes, err := toml.Marshal(newVersion)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, bytes, 0755)
}
