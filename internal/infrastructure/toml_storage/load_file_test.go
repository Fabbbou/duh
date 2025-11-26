package toml_storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFile(t *testing.T) {
	cfg, err := loadFile("test_file.toml")

	assert.NoError(t, err)

	assert.Len(t, cfg.Aliases, 2)
	assert.Equal(t, "ls -al", cfg.Aliases["ll"])
	assert.Equal(t, "echo 'Hello, World!'", cfg.Aliases["h-w"])

	assert.Len(t, cfg.Exports, 2)
	assert.Equal(t, "bar", cfg.Exports["FOO"])
	assert.Equal(t, "qux", cfg.Exports["BAZ"])
}
