// Package master contains the sqlite database schema
package master

import (
	"database/sql"

	_ "embed"

	"github.com/conneroisu/seltabl/tools/seltabls/data/generic"
)

// New creates a new queries type
func New(db generic.DBTX) *Queries {
	return &Queries{db: db}
}

// Queries is the object to use to interact with the database.
type Queries struct {
	db generic.DBTX
}

// WithTx creates a new queries type with a transaction.
func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

// MasterSchema is the schema for the main database
//
//go:embed combined/schema.sql
var MasterSchema string

// Selectors is a slice of selectors
type Selectors []Selector

// Len returns the length of the slice.
func (s Selectors) Len() int {
	return len(s)
}

// Less reports whether the element with
func (s Selectors) Less(i, j int) bool {
	return s[i].Occurances < s[j].Occurances
}

// Swap swaps the elements with indexes i and j.
func (s Selectors) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Strings returns a slice of strings from the Selectors
func (s Selectors) Strings() []string {
	strs := make([]string, len(s))
	for i, v := range s {
		strs[i] = v.Value
	}
	return strs
}
