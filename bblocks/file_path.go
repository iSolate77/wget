package bblocks

import (
	"os"
	"path/filepath"
)

func DetermineFilePath(outputFileName string) string {
	if *New_file_path != "" {
		cleanedPath := filepath.Clean(*New_file_path)
		homeDir, _ := os.UserHomeDir()
		filePath := filepath.Join(homeDir, cleanedPath[1:], outputFileName)
		*New_file_path = filePath
		return filePath
	}
	return outputFileName
}
