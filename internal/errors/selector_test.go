package errors

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
	output, err := genStructReprAndHighlight(&stc, "tr:nth-child(1) th:nth-child(1)")
	assert.Nil(t, err)
	file, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	file.WriteString(output)
	defer file.Close()
	t.Fail()
}
