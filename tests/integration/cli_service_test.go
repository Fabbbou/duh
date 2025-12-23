package integration

import (
	"duh/internal/domain/service"
	"duh/internal/infrastructure/file_db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) service.CliService {
	tempdir := t.TempDir()
	pathProvider := file_db.NewCustomPathProvider(tempdir)
	initService := file_db.NewInitDbService(pathProvider)
	hasChanged, err := initService.Run()
	assert.NoError(t, err)
	assert.Truef(t, hasChanged, "initialization should have made changes")
	fileDbRepo := file_db.NewFileDbRepository(pathProvider)

	return service.NewCliService(fileDbRepo)
}

