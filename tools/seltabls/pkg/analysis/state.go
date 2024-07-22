package analysis

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/conneroisu/seltabl/tools/seltabls/data"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/mitchellh/go-homedir"
)

// State is the state of the document analysis.
type State struct {
	// Map of file names to contents.
	Documents map[string]string
	// Selectors is the map of file names to selectors.
	Selectors map[string][]master.Selector
	// URLs is the map of file names to urls.
	URLs map[string][]string
	// Database is the database for the state.
	Database data.Database[master.Queries]
}

// NewState returns a new state with no documents
func NewState() (state State, err error) {
	ctx := context.Background()
	configPath, err := CreateConfigDir("~/.config/seltabls/")
	if err != nil {
		return state, fmt.Errorf("failed to create config directory: %w", err)
	}
	db, err := data.NewDb(
		ctx,
		master.New,
		&data.Config{
			Schema:   master.MasterSchema,
			URI:      "sqlite://uri.sqlite",
			FileName: path.Join(configPath, "uri.sqlite"),
		},
	)
	if err != nil {
		return state, fmt.Errorf("failed to create database: %w", err)
	}
	state = State{
		Documents: make(map[string]string),
		Selectors: make(map[string][]master.Selector),
		Database:  *db,
		URLs:      make(map[string][]string),
	}
	return state, nil
}

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
