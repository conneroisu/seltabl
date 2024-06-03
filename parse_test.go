package seltabl

import (
	"fmt"
	"testing"

	"github.com/conneroisu/seltabl/testdata"
	"github.com/stretchr/testify/assert"
)

type fixtureStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
}

func TestFixtureTables(t *testing.T) {
	p, err := NewFromString[fixtureStruct](testdata.FixtureABNumTable)
	assert.Nil(t, err)
	for _, pp := range p {
		fmt.Printf("pp %+v\n", pp)
	}
	assert.Equal(t, "1", p[0].A)
	assert.Equal(t, "2", p[0].B)
	assert.Equal(t, "3", p[1].A)
	assert.Equal(t, "4", p[1].B)
	assert.Equal(t, "5", p[2].A)
	assert.Equal(t, "6", p[2].B)
	assert.Equal(t, "7", p[3].A)
	assert.Equal(t, "8", p[3].B)
}
