package seltabl

import (
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// TestStruct is a test struct
type TestStruct struct {
	Name  string
	Age   int
	Score float64
}

// createSelectionFromHTML creates a goquery selection from a html string
func createSelectionFromHTML(html string) *goquery.Selection {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	return doc.Selection
}

// TestSetStructField tests the SetStructField function
// ensuring that it sets the correct value for the given fieldk
func TestSetStructField(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		html      string
		selector  string
		expected  any
		wantErr   bool
	}{
		{
			name:      "Set string field",
			fieldName: "Name",
			html:      `<div id="name">John Doe</div>`,
			selector:  innerTextSelector,
			expected:  "John Doe",
			wantErr:   false,
		},
		{
			name:      "Set int field",
			fieldName: "Age",
			html:      `<div id="age">30</div>`,
			selector:  innerTextSelector,
			expected:  30,
			wantErr:   false,
		},
		{
			name:      "Set float field",
			fieldName: "Score",
			html:      `<div id="score">99.5</div>`,
			selector:  innerTextSelector,
			expected:  99.5,
			wantErr:   false,
		},
		{
			name:      "Invalid field name",
			fieldName: "InvalidField",
			html:      `<div id="invalid">Invalid</div>`,
			selector:  innerTextSelector,
			expected:  nil,
			wantErr:   true,
		},
		{
			name:      "Invalid int value",
			fieldName: "Age",
			html:      `<div id="age">invalid</div>`,
			selector:  innerTextSelector,
			expected:  nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TestStruct{}
			cellValue := createSelectionFromHTML(tt.html)
			err := SetStructField(ts, tt.fieldName, cellValue, tt.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetStructField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				v := reflect.ValueOf(ts).Elem().FieldByName(tt.fieldName).Interface()
				if !reflect.DeepEqual(v, tt.expected) {
					t.Errorf("SetStructField() got = %v, want %v", v, tt.expected)
				}
			}
		})
	}
}
