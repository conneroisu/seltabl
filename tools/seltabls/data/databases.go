// Package data provides a set of data types for the database.
package data

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/conneroisu/seltabl/tools/seltabls/data/generic"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/schema"
	"golang.org/x/sync/errgroup"

	// Import the database dialects. Automatically registers dialects.
	_ "modernc.org/sqlite"
)

// SetupFunc is a function for setting up the database
type SetupFunc func(ctx context.Context, getenv func(string) string, db *bun.DB) error

// Database is a struct that holds the sql database and the queries.
// It uses generics to hold the appropriate type of query struct.
type Database[
	T master.Queries,
] struct {
	Queries *T
	Bun     *bun.DB
}

// NewSQLDatabase creates a new database struct with the sql database and the
// queries struct. It uses generics to return the appropriate type of query
func NewSQLDatabase[
	Q master.Queries,
](
	parentCtx context.Context,
	dialect schema.Dialect,
	db *sql.DB,
	newFunc func(generic.DBTX) *Q,
) (*Database[Q], error) {
	eg, ctx := errgroup.WithContext(parentCtx)
	q := &Database[Q]{
		Queries: newFunc(db),
		Bun:     bun.NewDB(db, dialect),
	}
	eg.Go(func() error {
		<-ctx.Done()
		err := db.Close()
		if err != nil {
			return fmt.Errorf("failed to close db: %v", err)
		}
		return nil
	})
	return q, nil
}

// Config is a struct that holds the configuration for a database.
type Config struct {
	// Schema is the schema for the libsql database
	Schema string
	// URI is the uri for the libsql database
	URI string
	// Name is the name for the libsql database
	Name string
	// FileName is the file name for the sqlite database
	FileName string
}

// NewDb sets up the database using the URI and optional options.
// Using generics to return the appropriate type of query struct,
// it creates a new database struct with the sql database and the
// queries struct utilizing the URI and optional options provided.
func NewDb[
	Q master.Queries,
](
	ctx context.Context,
	newFunc func(generic.DBTX) *Q,
	config *Config,
) (*Database[Q], error) {
	u, err := url.Parse(config.URI)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}
	switch u.Scheme {
	case "sqlite":
		fileName := config.FileName
		if fileName == "" {
			return nil, fmt.Errorf("file name is required")
		}
		db, err := sql.Open("sqlite", fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to open db: %v", err)
		}
		if _, err := db.Exec(config.Schema); err != nil {
			return nil, fmt.Errorf(
				"error executing schema %s: %v",
				config.Schema,
				err,
			)
		}
		return NewSQLDatabase(ctx, sqlitedialect.New(), db, newFunc)
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
}
