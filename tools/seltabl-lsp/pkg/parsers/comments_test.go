package parsers

import (
	"reflect"
	"testing"
)

// Test case structure
type testCase struct {
	name      string
	input     string
	expected  StructCommentData
	wantErr   bool
	wantedErr string
}

// TestParseStructComments tests the ParseStructComments function
func TestParseStructComments(t *testing.T) {
	// Define the test cases
	tests := []testCase{
		{
			name: "Single struct with URL and ignore-elements",
			input: `package main

// @url: https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html
// @ignore-elements: div, script
type Structure struct {
	Fields []Field
}`,
			expected: StructCommentData{
				URLs: []string{
					"https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
				},
				IgnoreElements: []string{"div", "script"},
			},
			wantErr: false,
		},
		{
			name: "Multiple structs with different comments",
			input: `package main

// @url: https://example.com/one
// @ignore-elements: p, span
type FirstStruct struct {
	Field1 string
}

// @url: https://example.com/two
// @ignore-elements: div, style
type SecondStruct struct {
	Field2 int
}`,
			expected: StructCommentData{
				URLs: []string{
					"https://example.com/one",
					"https://example.com/two",
				},
				IgnoreElements: []string{"p", "span", "div", "style"},
			},
			wantErr: false,
		},
		{
			name: "Struct without comments",
			input: `package main

type NoCommentStruct struct {
	Field string
}`,
			expected: StructCommentData{},
			wantErr:  true,
		},
		{
			name: "Struct with only URL",
			input: `package main

// @url: https://example.com/onlyurl
type URLOnlyStruct struct {
	Field string
}`,
			expected: StructCommentData{
				URLs: []string{"https://example.com/onlyurl"},
			},
			wantErr: false,
		},
		{
			name: "Struct with only ignore-elements",
			input: `package main

// @ignore-elements: img, link
type IgnoreOnlyStruct struct {
	Field string
}`,
			expected: StructCommentData{},
			wantErr:  true,
		},
		{
			name: "Struct with multiple URLs and ignore-elements",
			input: `package main

// @url: https://example.com/multiple1
// @url: https://example.com/multiple2
// @ignore-elements: header, footer
// @ignore-elements: nav
type MultipleAnnotationsStruct struct {
	Field string
}`,
			expected: StructCommentData{
				URLs: []string{
					"https://example.com/multiple1",
					"https://example.com/multiple2",
				},
				IgnoreElements: []string{"header", "footer", "nav"},
			},
			wantErr: false,
		},
		{
			name: "Mixed comments",
			input: `package main

// Some general comment
// @url: https://example.com/mixed
// Another general comment
// @ignore-elements: aside, section
type MixedCommentStruct struct {
	Field string
}`,
			expected: StructCommentData{
				URLs:           []string{"https://example.com/mixed"},
				IgnoreElements: []string{"aside", "section"},
			},
			wantErr: false,
		},
		{
			name: "Struct with spaced comments",
			input: `package main

// @url:    https://example.com/spaced
// @ignore-elements:   div ,  span 
type SpacedCommentStruct struct {
	Field string
}`,
			expected: StructCommentData{
				URLs:           []string{"https://example.com/spaced"},
				IgnoreElements: []string{"div", "span"},
			},
			wantErr: false,
		},
		{
			name: "Struct with invalid URL and ignore-elements formats",
			input: `package main

// @url: https://example.com/valid
// @ignore-elements: valid, invalid elements here
type InvalidFormatStruct struct {
	Field string
}`,
			expected: StructCommentData{
				URLs:           []string{"https://example.com/valid"},
				IgnoreElements: []string{"valid", "invalid elements here"},
			},
			wantErr: false,
		},
		{
			name: "Struct with additional comment markers",
			input: `package main

// This is a comment with a @url: https://example.com/additional
// and this @ignore-elements: div, span
type AdditionalCommentMarkersStruct struct {
	Field string
}`,
			expected: StructCommentData{
				URLs:           []string{"https://example.com/additional"},
				IgnoreElements: []string{"div", "span"},
			},
			wantErr: false,
		},
		{
			name: "Multiple structs with mixed comments",
			input: `package main

// @url: https://example.com/first
type FirstMixedStruct struct {
	Field string
}

// @ignore-elements: header, footer
type SecondMixedStruct struct {
	Field string
}`,
			expected: StructCommentData{
				URLs:           []string{"https://example.com/first"},
				IgnoreElements: []string{"header", "footer"},
			},
			wantErr: false,
		},
		{
			name: "Struct with no annotations",
			input: `package main

// This is a regular comment
// It should not be picked up
type NoAnnotationsStruct struct {
	Field string
}`,
			expected: StructCommentData{},
			wantErr:  true,
		},
		{
			name: "Struct with special characters in comments",
			input: `package main

// @url: https://example.com/special?chars&test=true
// @ignore-elements: #id, .class
type SpecialCharsStruct struct {
	Field string
}`,
			expected: StructCommentData{
				URLs: []string{
					"https://example.com/special?chars&test=true",
				},
				IgnoreElements: []string{"#id", ".class"},
			},
			wantErr: false,
		},
		{
			name: "Struct with comments in mixed order",
			input: `package main

// @ignore-elements: aside, section
// Some general comment
// @url: https://example.com/mixedorder
type MixedOrderStruct struct {
	Field string
}`,
			expected: StructCommentData{
				URLs:           []string{"https://example.com/mixedorder"},
				IgnoreElements: []string{"aside", "section"},
			},
			wantErr: false,
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseStructComments(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"ParseStructComments() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
			if !reflect.DeepEqual(actual, tt.expected) && tt.wantedErr == "" {
				t.Errorf(
					"ParseStructComments() = %v, want %v",
					actual,
					tt.expected,
				)
			} else if tt.wantedErr != "" {
				if err.Error() != tt.wantedErr {
					t.Errorf("ParseStructComments() error = %v, wantErr %v", err, tt.wantedErr)
				}
			}
		})
	}
}
