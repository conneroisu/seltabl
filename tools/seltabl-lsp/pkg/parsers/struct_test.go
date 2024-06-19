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
					Tgs: []*Tag{
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
					Tgs: []*Tag{
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
					Tgs: []*Tag{
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
