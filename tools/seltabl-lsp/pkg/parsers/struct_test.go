package parsers

import (
	"context"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

//go:embed struct.test.txt
var testContent string

// TestParseTags tests the ParseTags function
func TestParseTags(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		tag     string
		wantErr bool
	}{
		{
			name:    "valid tag",
			tag:     `json:"name,omitempty"`,
			wantErr: false,
		},
		{
			name:    "invalid tag",
			tag:     `json:name,omitempty"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTags(tt.tag, 0, 0, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTags_Get tests the Get function to get
func TestTags_Get(t *testing.T) {
	t.Parallel()
	tags := &Tags{
		tags: []*Tag{
			{Key: "json", Name: "name", Options: []string{"omitempty"}},
		},
	}

	tag, err := tags.Get("json")
	if err != nil {
		t.Errorf("Tags.Get() error = %v", err)
	}
	if tag.Key != "json" {
		t.Errorf("Tags.Get() = %v, want %v", tag.Key, "json")
	}
}

// TestTags_Set tests the Set function to set a tag
func TestTags_Set(t *testing.T) {
	t.Parallel()
	tags := &Tags{}

	err := tags.Set(
		&Tag{Key: "json", Name: "name", Options: []string{"omitempty"}},
	)
	if err != nil {
		t.Errorf("Tags.Set() error = %v", err)
	}

	tag, _ := tags.Get("json")
	if tag.Key != "json" {
		t.Errorf("Tags.Set() = %v, want %v", tag.Key, "json")
	}
}

// TestTags_AddOptions tests the AddOptions function for the Tags struct type
func TestTags_AddOptions(t *testing.T) {
	t.Parallel()
	tags := &Tags{
		tags: []*Tag{
			{Key: "json", Name: "name", Options: []string{"omitempty"}},
		},
	}

	tags.AddOptions("json", "readonly")
	tag, _ := tags.Get("json")
	if !tag.HasOption("readonly") {
		t.Errorf(
			"Tags.AddOptions() = %v, want %v",
			tag.Options,
			[]string{"omitempty", "readonly"},
		)
	}
}

// TestTags_Tags tests the Tags function
func TestTags_Tags(t *testing.T) {
	t.Parallel()
	tags := &Tags{
		tags: []*Tag{
			{Key: "json", Name: "name", Options: []string{"omitempty"}},
		},
	}

	got := tags.Tags()
	if len(got) != 1 || got[0].Key != "json" {
		t.Errorf("Tags.Tags() = %v, want %v", got, tags.tags)
	}
}

// TestTags_Keys tests the Keys function
func TestTags_Keys(t *testing.T) {
	t.Parallel()
	tags := &Tags{
		tags: []*Tag{
			{Key: "json", Name: "name", Options: []string{"omitempty"}},
		},
	}
	got := tags.Keys()
	if len(got) != 1 || got[0] != "json" {
		t.Errorf("Tags.Keys() = %v, want %v", got, []string{"json"})
	}
}

// TestTags_String tests the String function
func TestTags_String(t *testing.T) {
	t.Parallel()
	tags := &Tags{
		tags: []*Tag{
			{Key: "json", Name: "name", Options: []string{"omitempty"}},
		},
	}
	got := tags.String()
	want := `json:"name,omitempty"`
	if got != want {
		t.Errorf("Tags.String() = %v, want %v", got, want)
	}
}

// TestTag_HasOption tests the HasOption function
func TestTag_HasOption(t *testing.T) {
	t.Parallel()
	tag := &Tag{Key: "json", Name: "name", Options: []string{"omitempty"}}
	if !tag.HasOption("omitempty") {
		t.Errorf("Tag.HasOption() = %v, want %v", false, true)
	}
}

// TestTag_Value tests the Value function
func TestTag_Value(t *testing.T) {
	tag := &Tag{Key: "json", Name: "name", Options: []string{"omitempty"}}
	got := tag.Value()
	want := "name,omitempty"
	if got != want {
		t.Errorf("Tag.Value() = %v, want %v", got, want)
	}
}

// TestTag_String tests the String function
func TestTag_String(t *testing.T) {
	t.Parallel()
	tag := &Tag{Key: "json", Name: "name", Options: []string{"omitempty"}}
	got := tag.String()
	want := `json:"name,omitempty"`
	if got != want {
		t.Errorf("Tag.String() = %v, want %v", got, want)
	}
}

// TestParseStructs tests the ParseStructs function
func TestParseStructs(t *testing.T) {
	src := `package main

// @url: https://example.com/one
// @ignore-elements: img, link
type MyStruct struct {
	Field1 string ` + "`json:\"field1\"`" + `
	Field2 int    ` + "`json:\"field2\"`" + `
}
// @url: https://example.com/one
// @ignore-elements: img, link
type MyStruct2 struct {
	Field1 string ` + "`json:\"field1\"`" + `
	Field2 int    ` + "`json:\"field2\"`" + `
}
	`
	expected := []Structure{
		{
			Fields: []Field{
				{
					Name: "Field1",
					Type: "string",
					Tags: Tags{
						tags: []*Tag{
							{
								Key:     "json",
								Name:    "field1",
								Options: []string{},
							},
						},
					},
					Line: 5,
				},
				{
					Name: "Field2",
					Type: "int",
					Tags: Tags{
						tags: []*Tag{
							{
								Key:     "json",
								Name:    "field2",
								Options: []string{},
							},
						},
					},
					Line: 6,
				},
			},
		},
		{
			Fields: []Field{
				{
					Name: "Field1",
					Type: "string",
					Tags: Tags{
						tags: []*Tag{
							{
								Key:     "json",
								Name:    "field1",
								Options: []string{},
							},
						},
					},
					Line: 12,
				},
				{
					Name: "Field2",
					Type: "int",
					Tags: Tags{
						tags: []*Tag{
							{
								Key:     "json",
								Name:    "field2",
								Options: []string{},
							},
						},
					},
					Line: 13,
				},
			},
		},
	}
	structs, err := ParseStructs(context.Background(), []byte(src))
	if err != nil {
		t.Fatalf("Failed to parse struct: %v", err)
	}

	assert.Equal(t, 2, len(structs), "expected 2 structs")

	for i, structDef := range structs {
		assert.Equal(t, len(expected[i].Fields), len(structDef.Fields), "expected %d fields", i)
		t.Logf("struct %d: %v", i, structDef)
		// expect the field names to match
		assert.Equal(t, expected[i].Fields[0].Name, structDef.Fields[0].Name)
		assert.Equal(t, expected[i].Fields[1].Name, structDef.Fields[1].Name)
	}

}
