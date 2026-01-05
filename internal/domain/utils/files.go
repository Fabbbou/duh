package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Ensure all writes are flushed
	return destFile.Sync()
}

func ExpandUserPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func EnsureEscapeDoubleQuotes(input string) string {
	// Use a placeholder for already-escaped quotes
	temp := strings.ReplaceAll(input, `\"`, "\x00")
	// Escape all remaining quotes
	temp = strings.ReplaceAll(temp, `"`, `\"`)
	// Restore the already-escaped quotes
	return strings.ReplaceAll(temp, "\x00", `\"`)
}

func ParseCommaSeparatedValues(input string) []string {
	parts := strings.Split(input, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func JoinCommaSeparatedValues(values []string) string {
	return strings.Join(values, ",")
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func ListFilesInDirectory(dirPath string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(dirPath, entry.Name()))
		}
	}

	return files, nil
}

func ReadFileAsString(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetFileNameWithoutExtension(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

func SplitStringByNewLine(input string) []string {
	return strings.Split(input, "\n")
}

func JoinStringsWithNewLine(lines []string) string {
	return strings.Join(lines, "\n")
}
