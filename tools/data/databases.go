package data

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/conneroisu/seltabl/tools/data/generic"
	"github.com/conneroisu/seltabl/tools/data/master"
	"github.com/tursodatabase/go-libsql"
	"golang.org/x/sync/errgroup"

	// Import the database dialects. Automatically registers dialects.
	_ "modernc.org/sqlite"
)

// Database is a struct that holds the sql database and the queries.
// It uses generics to hold the appropriate type of query struct.
type Database[
	T master.Queries,
] struct {
	db      *sql.DB
	Queries *T
}

// NewDatabase creates a new database struct with the sql database and the
// queries struct. It uses generics to return the appropriate type of query
func NewDatabase[
	Q master.Queries,
](
	parentCtx context.Context,
	db *sql.DB,
	newFunc func(generic.DBTX) *Q,
) (*Database[Q], error) {
	eg, ctx := errgroup.WithContext(parentCtx)
	q := &Database[Q]{
		db:      db,
		Queries: newFunc(db),
	}
	eg.Go(func() error {
		<-ctx.Done()
		err := q.db.Close()
		if err != nil {
			return fmt.Errorf("failed to close db: %v", err)
		}
		return nil
	})
	return q, nil
}

// Config is a struct that holds a database configuration.
type Config struct {
	Schema string
	URI    string
	Name   string
	Opts   []libsql.Option
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
	path string,
) (*Database[Q], error) {
	u, err := url.Parse(config.URI)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}
	switch u.Scheme {
	case "sqlite":
		db, err := sql.Open("sqlite", path)
		if err != nil {
			return nil, fmt.Errorf("failed to open db: %v", err)
		}
		if _, err := db.Exec(config.Schema); err != nil {
			return nil, fmt.Errorf("error executing schema: %v", err)
		}
		return NewDatabase(ctx, db, newFunc)
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
}
