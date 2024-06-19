package parsers

import (
	"context"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

//go:embed struct.test.txt
var testContent string

// TestParseStruct tests the PositionStatusInStructTag function
func TestParseStruct(t *testing.T) {
	src := testContent
	ctx := context.Background()
	got, err := ParseStruct(ctx, []byte(src))
	expected := &Structure{
		Fields: []Field{
			{
				Name: "Name",
				Type: "string",
				Tags: Tags{
					tags: []*Tag{
						{
							Key:     "json",
							Name:    "name",
							Options: []string{},
						},
					},
				},
				Line: 5,
			},
			{
				Name: "Age",
				Type: "int",
				Tags: Tags{
					tags: []*Tag{
						{
							Key:     "json",
							Name:    "age",
							Options: []string{},
						},
					},
				},
				Line: 6,
			},
			{
				Name: "Address",
				Type: "string",
				Tags: Tags{
					tags: []*Tag{
						{
							Key:     "json",
							Name:    "address",
							Options: []string{"omitempty"},
						},
					},
				},
				Line: 7,
			},
		},
	}
	assert.NoError(t, err)
	for i, field := range expected.Fields {
		gotField := got.Fields[i]
		if gotField.Name == field.Name {
			assert.Equal(
				t,
				field.Type,
				gotField.Type,
				"field %s type not found",
				field.Name,
			)
		}
		// check the line
		if gotField.Line != field.Line {
			t.Errorf("field line: %d, expected: %d", gotField.Line, field.Line)
		}
	}
}

func TestParseStructee(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{
			name: "valid struct",
			src: `package main
				type User struct {
					Name string ` + "`json:\"name\"`" + `
				}`,
			wantErr: false,
		},
		{
			name: "invalid struct",
			src: `package main
				type User struct {
					Name string ` + "`json:\"name\"`" + `
					Age int ` + "`json:\"age\"`" + `
				`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := ParseStruct(ctx, []byte(tt.src))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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

	err := tags.Set(&Tag{Key: "json", Name: "name", Options: []string{"omitempty"}})
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
		t.Errorf("Tags.AddOptions() = %v, want %v", tag.Options, []string{"omitempty", "readonly"})
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
