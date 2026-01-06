package githubb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
)

func GetLatestRelease() (*GitHubRelease, error) {
	resp, err := http.Get("https://api.github.com/repos/Fabbbou/duh/releases/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release GitHubRelease
	err = json.Unmarshal(body, &release)
	if err != nil {
		return nil, err
	}

	return &release, nil
}

func GetAssetName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	if os == "windows" {
		return fmt.Sprintf("duh-%s-%s.exe", os, arch)
	}
	return fmt.Sprintf("duh-%s-%s", os, arch)
}
