package generic

import (
	"context"
	"database/sql"
)

// DBTX is the interface for the database/sql.Tx type
// and is used to simplify the queries interface by
// allowing the queries to be run within a transaction.
//
// Example:
//
//	tx, err := db.Begin()
//	if err != nil {
//	    return err
//	}
//	q := data.New(tx)
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
