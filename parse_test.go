package seltabl

import (
	"testing"

	"github.com/conneroisu/seltabl/testdata"
	"github.com/stretchr/testify/assert"
)

// TestFixtureTables tests the parsing of a table with headers.
func TestFixtureTables(t *testing.T) {
	p, err := NewFromString[testdata.FixtureStruct](testdata.FixtureABNumTable)
	assert.Nil(t, err)
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

// TestSuperNovaTable tests the parsing of a table with supernova data.
func TestSuperNovaTable(t *testing.T) {
	p, err := NewFromString[testdata.SuperNovaStruct](testdata.SuperNovaTable)
	assert.Nil(t, err)
	assert.Equal(t, "SN 1006", p[0].Supernova)
	assert.Equal(t, "1006", p[0].Year)
	assert.Equal(t, "Type Ia", p[0].Type)
	assert.Equal(t, "7,200", p[0].Distance)
	assert.Equal(t, "Brightest recorded supernova in history", p[0].Notes)

	assert.Equal(t, "SN 1054 (Crab Nebula)", p[1].Supernova)
	assert.Equal(t, "1054", p[1].Year)
	assert.Equal(t, "Type II", p[1].Type)
	assert.Equal(t, "6,500", p[1].Distance)
	assert.Equal(t, "Formed the Crab Nebula and pulsar", p[1].Notes)

	assert.Equal(t, "SN 1572 (Tycho's Supernova)", p[2].Supernova)
	assert.Equal(t, "1572", p[2].Year)
	assert.Equal(t, "Type Ia", p[2].Type)
	assert.Equal(t, "8,000-10,000", p[2].Distance)
	assert.Equal(t, "Observed by Tycho Brahe", p[2].Notes)

	assert.Equal(t, "SN 1604 (Kepler's Supernova)", p[3].Supernova)
	assert.Equal(t, "1604", p[3].Year)
	assert.Equal(t, "Type Ia", p[3].Type)
	assert.Equal(t, "20,000", p[3].Distance)
	assert.Equal(t, "Last observed supernova in the Milky Way", p[3].Notes)

}
