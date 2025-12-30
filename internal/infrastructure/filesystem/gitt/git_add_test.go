package gitt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CloneGitRepository(t *testing.T) {
	outputPath := t.TempDir()
	defer os.RemoveAll(outputPath)
	err := CloneGitRepository("https://github.com/isomorphic-git/test.empty", outputPath)
	assert.NoError(t, err)
	assert.DirExists(t, outputPath+"/.git")
}

func Test_extractGitRepoName(t *testing.T) {
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
		result := ExtractGitRepoName(tt.url)
		assert.Equal(t, tt.expected, result)
	}
}