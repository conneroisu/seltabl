package seltabl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fixture = `
	<table>
	<tr><td>W a</td><td>b</td></tr>
	<tr><td> 1 </td><td>2</td></tr>
	<tr><td>3  </td><td>4   </td></tr>
	<tr><td> 5 </td><td>   6</td></tr>
	<tr><td>7  </td><td>   8</td></tr>
	</table>
`

type fixtureStruct struct {
	A string `seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)"`
	B string `seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)"`
}

func TestFindsAllTables(t *testing.T) {
	p, err := NewFromString[fixtureStruct](fixture)
	assert.Nil(t, err)

	assert.NotEmpty(t, p)
	fmt.Printf("Fixure A: %s\n", p.A)
	fmt.Printf("Fixure B: %s\n", p.B)
	t.Fail()
}
