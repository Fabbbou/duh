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
