package domain

import (
	"os"

	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
)

// ConfigFile is a struct for a config file.
type ConfigFile struct {
	// Name is the name of the config file.
	Name string `yaml:"name"`
	// Description is the description of the config file.
	Description *string `yaml:"description,omitempty"`
	URLs        []struct {
		// URL is the url for the config file.
		URL string `yaml:"url"`
		// Sections is a list of sections in the html.
		Sections []Section `json:"sections" yaml:"sections"`
	} `yaml:"urls"`
}

// IdentifyResponse is a struct for the respond of an identify prompt.
//
// The identify prompt is used to describe the structure of a given
// html returning this struct in the form of json.
type IdentifyResponse struct {
	// Sections is a list of sections in the html.
	Sections []Section `json:"sections" yaml:"sections"`
	// Name is the name of the package.
	Name string `json:"name"     yaml:"name"`
}

// Section is a struct for a section in the html.
type Section struct {
	Name        string  `json:"name"        yaml:"name"`        // Name is the name of the section.
	Description string  `json:"description" yaml:"description"` // Description is a description of the section.
	CSS         string  `json:"css"         yaml:"css"`         // CSS is the css selector for the section.
	Fields      []Field `json:"fields"      yaml:"-"`           // Fields is a list of fields in the section.
}

// Field is a struct for a field
type Field struct {
	Name            string `json:"name"`             // Name is the name of the field.
	Type            string `json:"type"`             // Type is the type of the field.
	Description     string `json:"description"`      // Description is a description of the field.
	HeaderSelector  string `json:"header-selector"`  // HeaderSelector is the header selector for the field.
	DataSelector    string `json:"data-selector"`    // DataSelector is the data selector for the field.
	ControlSelector string `json:"control-selector"` // ControlSelector is the control selector for the field.
	QuerySelector   string `json:"query-selector"`   // QuerySelector is the query selector for the field.
	MustBePresent   string `json:"must-be-present"`  // MustBePresent is the must be present selector for the field.
}

// TestFile is a struct for a test file
type TestFile struct {
	// Name is the name of the test file
	Name string `json:"name" yaml:"name"`
	// URL is the url for the test file
	URL string `json:"url"  yaml:"url"`
	// PackageName is the package name for the test file
	PackageName string `json:"-"    yaml:"package-name"`
}

// WriteFile writes the test file to the file system
func (t *TestFile) WriteFile(p []byte) (n int, err error) {
	err = os.WriteFile(
		fmt.Sprintf("%s_test.go", t.Name),
		p,
		0644,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to write test file: %w", err)
	}
	return len(p), nil
}

// StructFile is a struct for a struct file.
//
// It contains attributes relating to the name, url, and ignore elements of the
// struct file.
type StructFile struct {
	// Name is the name of the struct file.
	Name string `json:"-"               yaml:"name"`
	// URL is the url for the struct file.
	URL string `json:"-"               yaml:"url"`
	// IgnoreElements is a list of elements to ignore when generating the
	// struct.
	IgnoreElements []string `json:"ignore-elements" yaml:"ignore-elements"`
	// Fields is a list of fields for the struct.
	Fields []Field `json:"fields"          yaml:"fields"`

	// TreeWidth is the width of the tree when generating the struct.
	TreeWidth int `json:"-" yaml:"tree-width"`

	// JSONValue is the json value for the struct yaml file.
	JSONValue string `json:"-" yaml:"json-value"`
	// HTMLContent is the html content for the struct file.
	HTMLContent string `json:"-" yaml:"html-content"`

	// Db is the database for the struct file.
	Db *data.Database[master.Queries] `json:"-" yaml:"-"`

	// Section is the section of the struct file.
	Section Section `json:"-" yaml:"section"`
}

// EmbedFile is a struct for the embed file.
// The embed file contains the html to use in the generated test file.
type EmbedFile struct {
	HTMLContent string `json:"html-content" yaml:"html-content"`
}
