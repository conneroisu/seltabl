package parse

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func GetConfigs(rootDir string) ([]fs.FileInfo, error) {
	var files []fs.FileInfo
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the path %v: %w", path, err)
		}
		// Check if the file matches the pattern
		if matched, err := filepath.Match("*_seltabl.yaml", info.Name()); err != nil {
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
