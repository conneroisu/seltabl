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
	baseModelType := reflect.TypeOf(TestStruct{})
	tests := []struct {
		name        string
		structPtr   *TestStruct
		cellHTML    string
		structField reflect.StructField
		selector    SelectorI
		wantErr     bool
		expected    interface{}
	}{
		{
			name:        "Set string field",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(0),
			cellHTML:    `<div id="name">John Doe</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    "John Doe",
		},
		{
			name:        "Set uint field",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(7),
			cellHTML:    `<div id="uint">123</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    uint(123),
		},
		{
			name:        "Set uint8 field",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(8),
			cellHTML:    `<div id="uint8">123</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    uint8(123),
		},
		{
			name:        "Set uint16 field",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(9),
			cellHTML:    `<div id="uint16">123</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    uint16(123),
		},
		{
			name:        "Set uint32 field",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(10),
			cellHTML:    `<div id="uint32">123</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    uint32(123),
		},
		{
			name:        "Set uint64 field",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(11),
			cellHTML:    `<div id="uint64">123</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    uint64(123),
		},
		{
			name:        "Set int field",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(1),
			cellHTML:    `<div id="age">30</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    30,
		},
		{
			name:        "Set float field (float64)",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(2),
			cellHTML:    `<div id="score">99.5</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    99.5,
		},
		{
			name:        "Set float field (float32)",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(5),
			cellHTML:    `<div id="score">99.5</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    float32(99.5),
		},
		{
			name:        "Invalid int value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(1),
			cellHTML:    `<div id="age">invalid</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Invalid uint value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(7),
			cellHTML:    `<div id="uint">invalid</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Invalid float value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(2),
			cellHTML:    `<div id="score">invalid</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Invalid Uint8",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(8),
			cellHTML:    `<div id="uint8">invalid</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Invalid Float32",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(5),
			cellHTML:    `<div id="float">invalid</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Invalid Uint16",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(9),
			cellHTML:    `<div id="uint16">invalid</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Invalid Uint32",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(10),
			cellHTML:    `<div id="uint32">invalid</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Invalid Uint64",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(11),
			cellHTML:    `<div id="uint64">invalid</div>`,
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Set string field",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(0),
			cellHTML:    "<div>Test String</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    "Test String",
		},
		{
			name:        "Set int field",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(1),
			cellHTML:    "<div>123</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    123,
		},
		{
			name:        "Set float field",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(4),
			cellHTML:    "<div>123.45</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    123.45,
		},
		{
			name:        "set float32 field",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(6),
			cellHTML:    "<div>1.45</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    float32(1.45),
		},
		{
			name:        "Selector not found (attr)",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(1),
			cellHTML:    "<div></div>",
			selector:    selector{control: ctlAttrSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Invalid int value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(1),
			cellHTML:    "<div>invalid</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Invalid float value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(2),
			cellHTML:    "<div>invalid</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:      "Invalid int8 value",
			structPtr: &TestStruct{},
			cellHTML:  "<div>invalid</div>",
			selector:  selector{control: ctlInnerTextSelector},
			wantErr:   true,
			expected:  nil,
		},
		{
			name:        "Invalid int16 value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(13),
			cellHTML:    "<div>invalid</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Invalid int32 value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(14),
			cellHTML:    "<div>invalid</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Invalid int64 value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(15),
			cellHTML:    "<div>invalid</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     true,
			expected:    nil,
		},
		{
			name:        "Valid int8 value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(12),
			cellHTML:    "<div>123</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    int8(123),
		},
		{
			name:        "Valid int16 value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(13),
			cellHTML:    "<div>123</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    int16(123),
		},
		{
			name:        "Valid int32 value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(14),
			cellHTML:    "<div>456</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    int32(456),
		},
		{
			name:        "Valid int64 value",
			structPtr:   &TestStruct{},
			structField: baseModelType.Field(15),
			cellHTML:    "<div>25565</div>",
			selector:    selector{control: ctlInnerTextSelector},
			wantErr:     false,
			expected:    int64(25565),
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
				tt.structField,
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
		})
	}
}
