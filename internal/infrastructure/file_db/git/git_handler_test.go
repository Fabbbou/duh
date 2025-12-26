package gitt_test

import (
	gitt "duh/internal/infrastructure/file_db/git"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CloneGitRepository(t *testing.T) {
	outputPath := t.TempDir()
	defer os.RemoveAll(outputPath)
	err := gitt.CloneGitRepository("https://github.com/isomorphic-git/test.empty", outputPath)
	assert.NoError(t, err)
	assert.DirExists(t, outputPath+"/.git")
}

func Test_ExtractGitRepoName(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://github.com/isomorphic-git/test.empty.git", "test.empty"},
		{"https://github.com/isomorphic-git/test.empty/", "test.empty"},
		{"https://github.com/isomorphic-git/test.empty", "test.empty"},
		{"https://gitlab.com/shchelchkov/tstation", "tstation"},
		{"git@gitlab.com:shchelchkov/tstation.git", "tstation"},
		{"", ""},
	}

	for _, tt := range tests {
		result := gitt.ExtractGitRepoName(tt.url)
		assert.Equal(t, tt.expected, result)
	}
}
