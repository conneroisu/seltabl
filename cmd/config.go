package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// configName is the pattern used to match config files
var configName = "*.go"

// GetConfigs returns a list of files that match the configName pattern
func GetConfigs(rootDir string) ([]fs.FileInfo, error) {
	var files []fs.FileInfo
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the path %v: %w", path, err)
		}
		// Check if the file matches the pattern
		if matched, err := filepath.Match(configName, info.Name()); err != nil {
			return fmt.Errorf("error matching file: %w", err)
		} else if matched {
			files = append(files, info)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the path %v: %w", rootDir, err)
	}
	return files, nil
}
