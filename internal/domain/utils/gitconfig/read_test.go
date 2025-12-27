package gitconfig

import (
	"duh/internal/domain/utils"
	"os"
	"slices"
	"testing"

	cc "github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/stretchr/testify/assert"
)

func Test_ReadConfigFile(t *testing.T) {
	cfg, err := readConfigFile("gitconfig.ini")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func Test_AddNewInclude(t *testing.T) {
	err := utils.CopyFile("gitconfig.ini", "gitconfig_test.ini")
	defer os.Remove("gitconfig_test.ini")
	assert.NoError(t, err)
	err = AddNewIncludeIfNotExists("path/to/include", "gitconfig_test.ini")
	assert.NoError(t, err)
	cfg, err := readConfigFile("gitconfig_test.ini")
	assert.NoError(t, err)
	sectionInclude := cfg.Raw.Section("include")
	assert.NotNil(t, sectionInclude)

	containsNewPath := slices.ContainsFunc(sectionInclude.Options, func(opt *cc.Option) bool {
		return opt.Key == "path" && opt.Value == "path/to/include"
	})
	assert.True(t, containsNewPath)
}
