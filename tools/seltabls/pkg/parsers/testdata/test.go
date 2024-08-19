package testdata

// Structure is a test struct
type Structure struct {
	Fields []Field
}

// Field is a struct for a field within a struct
type Field struct {
	Name string
	Type string
	Tags Tags
	Line int
}

// Tags is a struct for a tag
type Tags struct {
	tags []*Tag
}

// Tag is a struct for a tag
type Tag struct {
	Key     string
	Name    string
	Options []string
	Line    int
	Start   int
	End     int
}

// Tagger is an interface for tagging a struct
type Tagger interface {
	Get(string) (*Tag, error)
	Set(string, string)
	AddOptions(string, string)
}
