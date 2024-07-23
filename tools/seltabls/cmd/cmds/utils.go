package cmds

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
)

// CreateConfigDir creates a new config directory and returns the path.
func CreateConfigDir(dirPath string) (string, error) {
	path, err := homedir.Expand(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to expand home directory: %w", err)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		if os.IsExist(err) {
			return path, nil
		}
		return "", fmt.Errorf(
			"failed to create or find config directory: %w",
			err,
		)
	}
	return path, nil
}
