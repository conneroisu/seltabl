package domain

import (
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
)

// PromptFile is a struct for a prompt file that is used to keep track of
// the state of the generation process.
//
// It allows resumable generation processes by keeping track of the state of
// the generation process, and by storing information regarding the generation
// process.
type PromptFile struct {
	TreeWidth      int              `yaml:"tree-width"`               // required
	PackageName    string           `yaml:"name"`                     // required
	Description    string           `yaml:"description"`              // required
	URL            string           `yaml:"url"`                      // required
	IgnoreElements []string         `yaml:"ignore-elements"`          // required
	Selectors      master.Selectors `yaml:"selectors"`                // required
	Sections       []Section        `json:"sections" yaml:"sections"` // required
	SmartModel     string           `yaml:"model"`
	FastModel      string           `yaml:"fast-model"`
	RuledHTMLBody  string           `yaml:"ruled-html-body"`
}

// Section is a struct for a section in the html.
type Section struct {
	Name        string           `json:"name"        yaml:"name"`
	Description string           `json:"description" yaml:"description"`
	CSS         string           `json:"css"         yaml:"css"`
	Start       int              `json:"start"       yaml:"start"`
	End         int              `json:"end"         yaml:"end"`
	Fields      []Field          `json:"fields"      yaml:"fields"`
	Selectors   master.Selectors `json:"-" yaml:"selectors"`
	HTMLContent string           `json:"-" yaml:"html-content"`
}

// Field is a struct for a field.
type Field struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	Description     string `json:"description"`
	HeaderSelector  string `json:"header-selector"`
	DataSelector    string `json:"data-selector"`
	ControlSelector string `json:"control-selector"`
	QuerySelector   string `json:"query-selector"`
	MustBePresent   string `json:"must-be-present"`
}

// TestFile is a struct for a test file.
type TestFile struct {
	Name        string      `json:"name" yaml:"name"`
	URL         string      `json:"url"  yaml:"url"`
	PackageName string      `json:"-" yaml:"package-name"`
	ConfigFile  *PromptFile `json:"-" yaml:"config-file"`
	StructFile  *StructFile `json:"-" yaml:"struct-file"`
	Section     *Section    `json:"-" yaml:"section"`
}

// StructFile is a struct for a struct file.
//
// It contains attributes relating to the name, url, and ignore elements of the struct file.
type StructFile struct {
	File           os.File                        `json:"-" yaml:"-"`                             // required
	URL            string                         `json:"-" yaml:"url"`                           // required
	IgnoreElements []string                       `json:"ignore-elements" yaml:"ignore-elements"` // required
	TreeWidth      int                            `json:"-" yaml:"tree-width"`                    // required
	TreeDepth      int                            `json:"-" yaml:"tree-depth"`                    // required
	Section        Section                        `json:"-" yaml:"section"`                       // required
	PackageName    string                         `json:"-" yaml:"package-name"`                  // required
	ConfigFile     *PromptFile                    `json:"-" yaml:"config-file"`                   // required
	JSONValue      string                         `json:"-" yaml:"json-value"`                    // required
	HTMLContent    string                         `json:"-" yaml:"html-content"`                  // required
	Db             *data.Database[master.Queries] `json:"-" yaml:"-"`                             // required
}

// SectionsResponse is a struct for the response of a request to generate
// sections from a given url's html response.
//
// The sections generate prompt is used as input to describe the structure of a given
// html returning this struct in the form of json.
type SectionsResponse struct {
	// Sections is a list of sections in the html.
	Sections []Section `json:"sections" yaml:"sections"`
}
