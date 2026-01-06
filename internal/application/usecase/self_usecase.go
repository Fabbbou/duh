package usecase

import (
	"duh/internal/domain/port"
	"duh/internal/domain/utils/version"
	"duh/internal/infrastructure/githubb"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type SelfUsecase struct {
	dbPort port.DbPort
}

func NewSelfUsecase(dbPort port.DbPort) *SelfUsecase {
	return &SelfUsecase{
		dbPort: dbPort,
	}
}

func (p *SelfUsecase) GetBasePath() (string, error) {
	return p.dbPort.GetBasePath()
}

func (p *SelfUsecase) GetAllPaths() ([]string, error) {
	return p.dbPort.ListRepoPath()
}

func (p *SelfUsecase) RepositoriesPath() (string, error) {
	path, err := p.dbPort.GetBasePath()
	if err != nil {
		return "", err
	}
	repoPath := filepath.Join(path, "repositories")
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return "", nil
	}
	return repoPath, nil
}

func (s *SelfUsecase) GetVersion() string {
	return version.BuildInfo()
}

func (s *SelfUsecase) UpdateSelf() error {
	currentVersion := strings.TrimPrefix(version.BuildInfo(), "v")
	latestRelease, err := githubb.GetLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")
	if currentVersion == latestVersion {
		return fmt.Errorf("already running the latest version (%s)", currentVersion)
	}

	// Get the current executable path
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Find the appropriate asset for current OS and architecture
	assetName := githubb.GetAssetName()
	fmt.Println(assetName)
	var downloadURL string
	// fmt.Println(latestRelease.Assets)
	for _, asset := range latestRelease.Assets {
		fmt.Println(asset.Name)
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no compatible binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// Download the new binary
	tempFile := executable + ".tmp"
	err = githubb.DownloadFile(downloadURL, tempFile)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}

	// Replace the current executable
	backupFile := executable + ".backup"
	if err := os.Rename(executable, backupFile); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to backup current executable: %w", err)
	}

	if err := os.Rename(tempFile, executable); err != nil {
		os.Rename(backupFile, executable) // Restore backup
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	// Make the new file executable (Unix systems)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(executable, 0755); err != nil {
			return fmt.Errorf("failed to set executable permissions: %w", err)
		}
	}

	// Clean up backup file
	os.Remove(backupFile)

	return nil
}
