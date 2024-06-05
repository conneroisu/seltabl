package seltabl

import (
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// TestStruct is a test struct
type TestStruct struct {
	Name        string
	Age         int
	Score       float64
	IntField    int
	FloatField  float64
	FloatField2 float32
	StringField string
	UIntField   uint
	Uint8Field  uint8
	Uint16Field uint16
	Uint32Field uint32
	Uint64Field uint64
}

// TestSetStructField tests the SetStructField function
func TestSetStructField(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		structPtr *TestStruct
		fieldName string
		cellHTML  string
		selector  string
		wantErr   bool
		expected  interface{}
	}{
		{
			name:      "Set string field",
			structPtr: &TestStruct{},
			fieldName: "Name",
			cellHTML:  `<div id="name">John Doe</div>`,
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  "John Doe",
		},
		{
			name:      "Set uint field",
			structPtr: &TestStruct{},
			fieldName: "UIntField",
			cellHTML:  `<div id="uint">123</div>`,
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  uint(123),
		},
		{
			name:      "Set uint8 field",
			structPtr: &TestStruct{},
			fieldName: "Uint8Field",
			cellHTML:  `<div id="uint8">123</div>`,
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  uint8(123),
		},
		{
			name:      "Set uint16 field",
			structPtr: &TestStruct{},
			fieldName: "Uint16Field",
			cellHTML:  `<div id="uint16">123</div>`,
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  uint16(123),
		},
		{
			name:      "Set uint32 field",
			structPtr: &TestStruct{},
			fieldName: "Uint32Field",
			cellHTML:  `<div id="uint32">123</div>`,
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  uint32(123),
		},
		{
			name:      "Set uint64 field",
			structPtr: &TestStruct{},
			fieldName: "Uint64Field",
			cellHTML:  `<div id="uint64">123</div>`,
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  uint64(123),
		},
		{
			name:      "Set int field",
			structPtr: &TestStruct{},
			fieldName: "Age",
			cellHTML:  `<div id="age">30</div>`,
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  30,
		},
		{
			name:      "Set float field (float64)",
			structPtr: &TestStruct{},
			fieldName: "Score",
			cellHTML:  `<div id="score">99.5</div>`,
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  99.5,
		},
		{
			name:      "Set float field (float32)",
			structPtr: &TestStruct{},
			fieldName: "Score",
			cellHTML:  `<div id="score">99.5</div>`,
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  99.5,
		},
		{
			name:      "Invalid field name",
			structPtr: &TestStruct{},
			fieldName: "InvalidField",
			cellHTML:  `<div id="invalid">Invalid</div>`,
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid int value",
			structPtr: &TestStruct{},
			fieldName: "Age",
			cellHTML:  `<div id="age">invalid</div>`,
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid uint value",
			structPtr: &TestStruct{},
			fieldName: "UIntField",
			cellHTML:  `<div id="uint">invalid</div>`,
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid float value",
			structPtr: &TestStruct{},
			fieldName: "Score",
			cellHTML:  `<div id="score">invalid</div>`,
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint8",
			structPtr: &TestStruct{},
			fieldName: "Uint8Field",
			cellHTML:  `<div id="uint8">invalid</div>`,
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Float32",
			structPtr: &TestStruct{},
			fieldName: "FloatField2",
			cellHTML:  `<div id="float">invalid</div>`,
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint16",
			structPtr: &TestStruct{},
			fieldName: "Uint16Field",
			cellHTML:  `<div id="uint16">invalid</div>`,
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint32",
			structPtr: &TestStruct{},
			fieldName: "Uint32Field",
			cellHTML:  `<div id="uint32">invalid</div>`,
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint64",
			structPtr: &TestStruct{},
			fieldName: "Uint64Field",
			cellHTML:  `<div id="uint64">invalid</div>`,
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Set string field",
			structPtr: &TestStruct{},
			fieldName: "StringField",
			cellHTML:  "<div>Test String</div>",
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  "Test String",
		},
		{
			name:      "Set int field",
			structPtr: &TestStruct{},
			fieldName: "IntField",
			cellHTML:  "<div>123</div>",
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  123,
		},
		{
			name:      "Set float field",
			structPtr: &TestStruct{},
			fieldName: "FloatField",
			cellHTML:  "<div>123.45</div>",
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  123.45,
		},
		{
			name:      "Field does not exist",
			structPtr: &TestStruct{},
			fieldName: "NonExistentField",
			cellHTML:  "<div>Test</div>",
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Field cannot be set",
			structPtr: &TestStruct{},
			fieldName: "UnexportedField",
			cellHTML:  "<div>456</div>",
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Selector not found (innerText)",
			structPtr: &TestStruct{},
			fieldName: "aa",
			cellHTML:  "<div></div>",
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Selector not found (attr)",
			structPtr: &TestStruct{},
			fieldName: "StringField",
			cellHTML:  "<div></div>",
			selector:  attrSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid int value",
			structPtr: &TestStruct{},
			fieldName: "IntField",
			cellHTML:  "<div>invalid</div>",
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid float value",
			structPtr: &TestStruct{},
			fieldName: "FloatField",
			cellHTML:  "<div>invalid</div>",
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "valid float32",
			structPtr: &TestStruct{},
			fieldName: "FloatField2",
			cellHTML:  "<div>1.45</div>",
			selector:  innerTextSelector,
			wantErr:   false,
			expected:  float32(1.45),
		},

		{
			name:      "invalid float32",
			structPtr: &TestStruct{},
			fieldName: "FloatField2",
			cellHTML:  "<div>1.23.45</div>",
			selector:  innerTextSelector,
			wantErr:   true,
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.cellHTML))
			if err != nil {
				t.Fatalf("failed to create document: %v", err)
			}
			cellValue := doc.Find("div")

			err = SetStructField(tt.structPtr, tt.fieldName, cellValue, tt.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetStructField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				v := reflect.ValueOf(tt.structPtr).Elem().FieldByName(tt.fieldName).Interface()
				if !reflect.DeepEqual(v, tt.expected) {
					t.Errorf("SetStructField() = %v, expected %v", v, tt.expected)
				}
			}
		})
	}
}
