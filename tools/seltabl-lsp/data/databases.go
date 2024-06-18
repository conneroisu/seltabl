// Package data provides a set of data types for the database.
package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/schema"

	// Import the database dialects. Automatically registers dialects.
	_ "modernc.org/sqlite"
)

// SetupFunc is a function for setting up the database
type SetupFunc func(ctx context.Context, getenv func(string) string, db *bun.DB) error

// NewDb sets up the database using the URI and optional options.
// Using generics to return the appropriate type of query struct,
// it creates a new database struct with the sql database and the
// queries struct utilizing the URI and optional options provided.
func NewDb(
	ctx context.Context,
	getenv func(string) string,
	path string,
	fn SetupFunc,
	dialect schema.Dialect,
) (*bun.DB, error) {
	db, err := sql.Open(sqliteshim.ShimName, path)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	bu := bun.NewDB(db, dialect)
	if err = fn(ctx, getenv, bu); err != nil {
		return nil, fmt.Errorf("failed to apply schema: %w", err)
	}
	return bu, nil
}

// Selector is a struct for a selector
type Selector struct {
	bun.BaseModel `bun:"table:selectors,alias:s"`
	// ID is the id of the selector
	ID int64 `bun:"id,pk,autoincrement"`
	// Selector is the selector for the selector
	Selector string `bun:"selector"`
	// Context is the html context for the selector
	Context string `bun:"context"`
	// URL is the url for the selector
	URLID int64 `bun:"url_id"`
	// URL is the url for the selector
	URL URL `bun:"rel:belongs-to,join:url_id=id"`
}

// URL is a struct for a url
type URL struct {
	bun.BaseModel `bun:"table:urls"`
	// ID is the id of the url
	ID int64 `bun:"id,pk,autoincrement"`
	// URL is the url for the url
	URL string `bun:"url"`
	// HTMLID is the id of the html
	HTMLID int64 `bun:"html_id"`
	// Content is the content of the url response
	HTML HTML `bun:"rel:belongs-to,join:html_id=id"`
}

// HTML is a struct for a html
type HTML struct {
	bun.BaseModel `bun:"table:html,alias:h"`
	// ID is the id of the html
	ID int64 `bun:"id,pk,autoincrement"`
	// HTML is the html for the html
	HTML string `bun:"html"`
	// URL is the url for the html
	URL string `bun:"url"`
}
