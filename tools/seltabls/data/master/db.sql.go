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
