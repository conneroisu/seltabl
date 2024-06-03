package seltabl

import (
	"fmt"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/conneroisu/seltabl/testdata"
	"github.com/stretchr/testify/assert"
)

// TestFixtureTables tests the parsing of a table with headers.
func TestFixtureTables(t *testing.T) {
	p, err := NewFromString[testdata.FixtureStruct](testdata.FixtureABNumTable)
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

// TestNumberedTable tests the parsing of a table with numbered headers.
func TestNumberedTable(t *testing.T) {
	p, err := NewFromString[testdata.NumberedStruct](testdata.NumberedTable)
	assert.Nil(t, err)
	assert.NoError(t, err)

	assert.Equal(t, "Row 1, Cell 1", p[0].Header1)
	assert.Equal(t, "Row 1, Cell 2", p[0].Header2)
	assert.Equal(t, "Row 1, Cell 3", p[0].Header3)
	assert.Equal(t, "Row 2, Cell 1", p[1].Header1)
	assert.Equal(t, "Row 2, Cell 2", p[1].Header2)
	assert.Equal(t, "Row 2, Cell 3", p[1].Header3)
	assert.Equal(t, "Row 3, Cell 1", p[2].Header1)
	assert.Equal(t, "Row 3, Cell 2", p[2].Header2)
	assert.Equal(t, "Row 3, Cell 3", p[2].Header3)
}

// <table>
//         <tr>
//                 <th>Supernova</th>
//                 <th>Year</th>
//                 <th>Type</th>
//                 <th>Distance (light-years)</th>
//                 <th>Notes</th>
//         </tr>
//         <tr>
//                 <td>SN 1006</td>
//                 <td>1006</td>
//                 <td>Type Ia</td>
//                 <td>7,200</td>
//                 <td>Brightest recorded supernova in history</td>
//         </tr>
//         <tr>
//                 <td>SN 1054 (Crab Nebula)</td>
//                 <td>1054</td>
//                 <td>Type II</td>
//                 <td>6,500</td>
//                 <td>Formed the Crab Nebula and pulsar</td>
//         </tr>
//         <tr>
//                 <td>SN 1572 (Tycho's Supernova)</td>
//                 <td>1572</td>
//                 <td>Type Ia</td>
//                 <td>8,000-10,000</td>
//                 <td>Observed by Tycho Brahe</td>
//         </tr>
//         <tr>
//                 <td>SN 1604 (Kepler's Supernova)</td>
//                 <td>1604</td>
//                 <td>Type Ia</td>
//                 <td>20,000</td>
//                 <td>Last observed supernova in the Milky Way</td>
//         </tr>
//         <tr>
//                 <td>SN 1987A</td>
//                 <td>1987</td>
//                 <td>Type II</td>
//                 <td>168,000</td>
//                 <td>Closest observed supernova since 1604</td>
//         </tr>
//         <tr>
//                 <td>SN 1993J</td>
//                 <td>1993</td>
//                 <td>Type IIb</td>
//                 <td>11,000,000</td>
//                 <td>In the galaxy M81</td>
//         </tr>
// </table>

func TestSuperNovaTable(t *testing.T) {
	p, err := NewFromString[testdata.SuperNovaStruct](testdata.SuperNovaTable)
	assert.Nil(t, err)

	var data testdata.SuperNovaStruct
	// Marshal
	output, err := sonic.Marshal(&data)
}
