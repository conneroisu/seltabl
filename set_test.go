package seltabl

import (
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// TestStruct is a test struct
// fields in this struct are used in the tests
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
	Int8Field   int8
	Int16Field  int16
	Int32Field  int32
	Int64Field  int64
}

// TestSetStructField tests the SetStructField function
func TestSetStructField(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		structPtr *TestStruct
		fieldName string
		cellHTML  string
		selector  SelectorI
		wantErr   bool
		expected  interface{}
	}{
		{
			name:      "Set string field",
			structPtr: &TestStruct{},
			fieldName: "Name",
			cellHTML:  `<div id="name">John Doe</div>`,
			selector:  selector{control: ctlInnerTextSelector},
			wantErr:   false,
			expected:  "John Doe",
		},
		{
			name:      "Set uint field",
			structPtr: &TestStruct{},
			fieldName: "UIntField",
			cellHTML:  `<div id="uint">123</div>`,
			selector:  selector{control: ctlInnerTextSelector},
			wantErr:   false,
			expected:  uint(123),
		},
		{
			name:      "Set uint8 field",
			structPtr: &TestStruct{},
			fieldName: "Uint8Field",
			cellHTML:  `<div id="uint8">123</div>`,
			selector:  selector{control: ctlInnerTextSelector},
			wantErr:   false,
			expected:  uint8(123),
		},
		{
			name:      "Set uint16 field",
			structPtr: &TestStruct{},
			fieldName: "Uint16Field",
			cellHTML:  `<div id="uint16">123</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  uint16(123),
		},
		{
			name:      "Set uint32 field",
			structPtr: &TestStruct{},
			fieldName: "Uint32Field",
			cellHTML:  `<div id="uint32">123</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  uint32(123),
		},
		{
			name:      "Set uint64 field",
			structPtr: &TestStruct{},
			fieldName: "Uint64Field",
			cellHTML:  `<div id="uint64">123</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  uint64(123),
		},
		{
			name:      "Set int field",
			structPtr: &TestStruct{},
			fieldName: "Age",
			cellHTML:  `<div id="age">30</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  30,
		},
		{
			name:      "Set float field (float64)",
			structPtr: &TestStruct{},
			fieldName: "Score",
			cellHTML:  `<div id="score">99.5</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  99.5,
		},
		{
			name:      "Set float field (float32)",
			structPtr: &TestStruct{},
			fieldName: "Score",
			cellHTML:  `<div id="score">99.5</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  99.5,
		},
		{
			name:      "Invalid field name",
			structPtr: &TestStruct{},
			fieldName: "InvalidField",
			cellHTML:  `<div id="invalid">Invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid int value",
			structPtr: &TestStruct{},
			fieldName: "Age",
			cellHTML:  `<div id="age">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid uint value",
			structPtr: &TestStruct{},
			fieldName: "UIntField",
			cellHTML:  `<div id="uint">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid float value",
			structPtr: &TestStruct{},
			fieldName: "Score",
			cellHTML:  `<div id="score">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint8",
			structPtr: &TestStruct{},
			fieldName: "Uint8Field",
			cellHTML:  `<div id="uint8">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Float32",
			structPtr: &TestStruct{},
			fieldName: "FloatField2",
			cellHTML:  `<div id="float">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint16",
			structPtr: &TestStruct{},
			fieldName: "Uint16Field",
			cellHTML:  `<div id="uint16">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint32",
			structPtr: &TestStruct{},
			fieldName: "Uint32Field",
			cellHTML:  `<div id="uint32">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint64",
			structPtr: &TestStruct{},
			fieldName: "Uint64Field",
			cellHTML:  `<div id="uint64">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Set string field",
			structPtr: &TestStruct{},
			fieldName: "StringField",
			cellHTML:  "<div>Test String</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  "Test String",
		},
		{
			name:      "Set int field",
			structPtr: &TestStruct{},
			fieldName: "IntField",
			cellHTML:  "<div>123</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  123,
		},
		{
			name:      "Set float field",
			structPtr: &TestStruct{},
			fieldName: "FloatField",
			cellHTML:  "<div>123.45</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  123.45,
		},
		{
			name:      "set float32 field",
			structPtr: &TestStruct{},
			fieldName: "FloatField2",
			cellHTML:  "<div>1.45</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  float32(1.45),
		},
		{
			name:      "Field does not exist",
			structPtr: &TestStruct{},
			fieldName: "NonExistentField",
			cellHTML:  "<div>Test</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Field cannot be set",
			structPtr: &TestStruct{},
			fieldName: "UnexportedField",
			cellHTML:  "<div>456</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Selector not found (innerText)",
			structPtr: &TestStruct{},
			fieldName: "aa",
			cellHTML:  "<div></div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Selector not found (attr)",
			structPtr: &TestStruct{},
			fieldName: "StringField",
			cellHTML:  "<div></div>",
			selector:  selector{ctlAttrSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid int value",
			structPtr: &TestStruct{},
			fieldName: "IntField",
			cellHTML:  "<div>invalid</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid float value",
			structPtr: &TestStruct{},
			fieldName: "FloatField",
			cellHTML:  "<div>invalid</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "invalid float32",
			structPtr: &TestStruct{},
			fieldName: "FloatField2",
			cellHTML:  "<div>1.23.45</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid int8 value",
			structPtr: &TestStruct{},
			fieldName: "Int8Field",
			cellHTML:  "<div>invalid</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid int16 value",
			structPtr: &TestStruct{},
			fieldName: "Int16Field",
			cellHTML:  "<div>invalid</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid int32 value",
			structPtr: &TestStruct{},
			fieldName: "Int32Field",
			cellHTML:  "<div>invalid</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid int64 value",
			structPtr: &TestStruct{},
			fieldName: "Int64Field",
			cellHTML:  "<div>invalid</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Valid int8 value",
			structPtr: &TestStruct{},
			fieldName: "Int8Field",
			cellHTML:  "<div>123</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  int8(123),
		},
		{
			name:      "Valid int16 value",
			structPtr: &TestStruct{},
			fieldName: "Int16Field",
			cellHTML:  "<div>123</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  int16(123),
		},
		{
			name:      "Valid int32 value",
			structPtr: &TestStruct{},
			fieldName: "Int32Field",
			cellHTML:  "<div>456</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  int32(456),
		},
		{
			name:      "Valid int64 value",
			structPtr: &TestStruct{},
			fieldName: "Int64Field",
			cellHTML:  "<div>25565</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  int64(25565),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			doc, err := goquery.NewDocumentFromReader(
				strings.NewReader(tt.cellHTML),
			)
			if err != nil {
				t.Fatalf("failed to create document: %v", err)
			}
			cellValue := doc.Find("div")
			err = SetStructField(
				tt.structPtr,
				tt.fieldName,
				cellValue,
				tt.selector,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"SetStructField() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if err == nil {
				v := reflect.ValueOf(tt.structPtr).
					Elem().
					FieldByName(tt.fieldName).
					Interface()
				if !reflect.DeepEqual(v, tt.expected) {
					t.Errorf(
						"SetStructField() = %v, expected %v",
						v,
						tt.expected,
					)
				}
			}
		})
	}
}

// BenchStruct is a test struct
// fields in this struct are used in the tests
type BenchStruct struct {
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
	StructField *BenchStruct
}

// BenchSetStructField tests the SetStructField function
func BenchSetStructField(t *testing.B) {
	Benchs := []struct {
		name      string
		structPtr *BenchStruct
		fieldName string
		cellHTML  string
		selector  SelectorI
		wantErr   bool
		expected  interface{}
	}{
		{
			name:      "Set string field",
			structPtr: &BenchStruct{},
			fieldName: "Name",
			cellHTML:  `<div id="name">John Doe</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  "John Doe",
		},
		{
			name:      "Set uint field",
			structPtr: &BenchStruct{},
			fieldName: "UIntField",
			cellHTML:  `<div id="uint">123</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  uint(123),
		},
		{
			name:      "Set uint8 field",
			structPtr: &BenchStruct{},
			fieldName: "Uint8Field",
			cellHTML:  `<div id="uint8">123</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  uint8(123),
		},
		{
			name:      "Set uint16 field",
			structPtr: &BenchStruct{},
			fieldName: "Uint16Field",
			cellHTML:  `<div id="uint16">123</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  uint16(123),
		},
		{
			name:      "Set uint32 field",
			structPtr: &BenchStruct{},
			fieldName: "Uint32Field",
			cellHTML:  `<div id="uint32">123</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  uint32(123),
		},
		{
			name:      "Set uint64 field",
			structPtr: &BenchStruct{},
			fieldName: "Uint64Field",
			cellHTML:  `<div id="uint64">123</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  uint64(123),
		},
		{
			name:      "Set int field",
			structPtr: &BenchStruct{},
			fieldName: "Age",
			cellHTML:  `<div id="age">30</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  30,
		},
		{
			name:      "Set float field (float64)",
			structPtr: &BenchStruct{},
			fieldName: "Score",
			cellHTML:  `<div id="score">99.5</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  99.5,
		},
		{
			name:      "Set float field (float32)",
			structPtr: &BenchStruct{},
			fieldName: "Score",
			cellHTML:  `<div id="score">99.5</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  99.5,
		},
		{
			name:      "Invalid field name",
			structPtr: &BenchStruct{},
			fieldName: "InvalidField",
			cellHTML:  `<div id="invalid">Invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid int value",
			structPtr: &BenchStruct{},
			fieldName: "Age",
			cellHTML:  `<div id="age">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid uint value",
			structPtr: &BenchStruct{},
			fieldName: "UIntField",
			cellHTML:  `<div id="uint">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid float value",
			structPtr: &BenchStruct{},
			fieldName: "Score",
			cellHTML:  `<div id="score">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint8",
			structPtr: &BenchStruct{},
			fieldName: "Uint8Field",
			cellHTML:  `<div id="uint8">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Float32",
			structPtr: &BenchStruct{},
			fieldName: "FloatField2",
			cellHTML:  `<div id="float">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint16",
			structPtr: &BenchStruct{},
			fieldName: "Uint16Field",
			cellHTML:  `<div id="uint16">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint32",
			structPtr: &BenchStruct{},
			fieldName: "Uint32Field",
			cellHTML:  `<div id="uint32">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid Uint64",
			structPtr: &BenchStruct{},
			fieldName: "Uint64Field",
			cellHTML:  `<div id="uint64">invalid</div>`,
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Set string field",
			structPtr: &BenchStruct{},
			fieldName: "StringField",
			cellHTML:  "<div>Bench String</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  "Bench String",
		},
		{
			name:      "Set int field",
			structPtr: &BenchStruct{},
			fieldName: "IntField",
			cellHTML:  "<div>123</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  123,
		},
		{
			name:      "Set float field",
			structPtr: &BenchStruct{},
			fieldName: "FloatField",
			cellHTML:  "<div>123.45</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  123.45,
		},
		{
			name:      "set float32 field",
			structPtr: &BenchStruct{},
			fieldName: "FloatField2",
			cellHTML:  "<div>1.45</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   false,
			expected:  float32(1.45),
		},
		{
			name:      "Field does not exist",
			structPtr: &BenchStruct{},
			fieldName: "NonExistentField",
			cellHTML:  "<div>Bench</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Field cannot be set",
			structPtr: &BenchStruct{},
			fieldName: "UnexportedField",
			cellHTML:  "<div>456</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Selector not found (innerText)",
			structPtr: &BenchStruct{},
			fieldName: "aa",
			cellHTML:  "<div></div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Selector not found (attr)",
			structPtr: &BenchStruct{},
			fieldName: "StringField",
			cellHTML:  "<div></div>",
			selector:  selector{ctlAttrSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid int value",
			structPtr: &BenchStruct{},
			fieldName: "IntField",
			cellHTML:  "<div>invalid</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "Invalid float value",
			structPtr: &BenchStruct{},
			fieldName: "FloatField",
			cellHTML:  "<div>invalid</div>",
			selector:  selector{ctlInnerTextSelector, ""},
			wantErr:   true,
			expected:  nil,
		},
	}
	for _, tt := range Benchs {
		t.Run(tt.name, func(t *testing.B) {
			tt := tt
			doc, err := goquery.NewDocumentFromReader(
				strings.NewReader(tt.cellHTML),
			)
			if err != nil {
				t.Fatalf("failed to create document: %v", err)
			}
			cellValue := doc.Find("div")
			err = SetStructField(
				tt.structPtr,
				tt.fieldName,
				cellValue,
				tt.selector,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"SetStructField() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if err == nil {
				v := reflect.ValueOf(tt.structPtr).
					Elem().
					FieldByName(tt.fieldName).
					Interface()
				if !reflect.DeepEqual(v, tt.expected) {
					t.Errorf(
						"SetStructField() = %v, expected %v",
						v,
						tt.expected,
					)
				}
			}
		})
	}

}
