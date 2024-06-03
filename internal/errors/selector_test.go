package errors

import (
	"os"
	"reflect"
	"testing"
)

type SuperNovaStruct struct {
	Supernova string `json:"Supernova" seltabl:"Supernova" hSel:"tr:nth-child(1) th:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	Year      string `json:"Year" seltabl:"Year" hSel:"tr:nth-child(1) th:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
	Type      string `json:"Type" seltabl:"Type" hSel:"tr:nth-child(1) th:nth-child(3)" dSel:"tr td:nth-child(3)" cSel:"$text"`
	Distance  string `seltabl:"Distance" hSel:"tr:nth-child(1) th:nth-child(4)" dSel:"tr td:nth-child(4)" json:"Distance" `
	Notes     string `json:"Notes" seltabl:"Notes" hSel:"tr:nth-child(1) th:nth-child(5)" dSel:"tr td:nth-child(5)"`
}

func TestStructErrors(t *testing.T) {
	stc := SuperNovaStruct{}
	output := genStructReprAndHighlight(&stc, "tr:nth-child(1) th:nth-child(1)")
	file, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	file.WriteString(output)
	defer file.Close()
	t.Fail()
}

func TestGenStructTagString(t *testing.T) {

	type TestCase struct {
		Field             reflect.StructField
		HighlightSelector string
		Expected          string
	}

	tests := []TestCase{
		{
			Field: reflect.StructField{
				Name: "Field1",
				Type: reflect.TypeOf(""),
				Tag:  `json:"field1" xml:"field1"`,
			},
			HighlightSelector: "json",
			Expected:          " Field1: string `hSel:json:\"field1\" xml:\"field1\" `",
		},
		{
			Field: reflect.StructField{
				Name: "Field2",
				Type: reflect.TypeOf(0),
				Tag:  `db:"field2" xml:"field2"`,
			},
			HighlightSelector: "db",
			Expected:          " Field2: int `hSel:db:\"field2\" xml:\"field2\" `",
		},
		{
			Field: reflect.StructField{
				Name: "Field3",
				Type: reflect.TypeOf(true),
				Tag:  `json:"field3" xml:"field3"`,
			},
			HighlightSelector: "xml",
			Expected:          " Field3: bool `json:\"field3\" hSel:xml:\"field3\" `",
		},
	}

	for _, test := range tests {
		result := genStructTagString(test.Field, test.HighlightSelector)
		if result != test.Expected {
			t.Errorf("genStructTagString(%v, %v) = %v; expected %v",
				test.Field, test.HighlightSelector, result, test.Expected)
		}
	}
}
