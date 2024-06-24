// Package logs contains the sqlite database schema
package logs

import (
	"database/sql"
	"encoding/json"

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

// // LevelWriter defines as interface a writer may implement in order
// // to receive level information with payload.
//
//	type LevelWriter interface {
//	        io.Writer
//	        WriteLevel(level Level, p []byte) (n int, err error)
//	}

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

// Prime is the prime for the logs table
type Prime struct {
	Log          Log          `json:"log"`
	LogLevel     LogLevel     `json:"log_level"`
	Notification Notification `json:"notification,omitempty"`
	Request      Request      `json:"request,omitempty"`
	Response     Response     `json:"response,omitempty"`
}

// Write writes the given bytes to the database
func (q *Queries) Write(p []byte) (n int, err error) {
	// decode the bytes into a string
	var log Prime
	err = json.Unmarshal(p, log)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
