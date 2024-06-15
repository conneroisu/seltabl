package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/conneroisu/seltabl/tools/pkg/lsp"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"

	// Import the database dialects. Automatically registers dialects.
	_ "modernc.org/sqlite"
)

// NewDb sets up the database using the URI and optional options.
// Using generics to return the appropriate type of query struct,
// it creates a new database struct with the sql database and the
// queries struct utilizing the URI and optional options provided.
func NewDb(
	ctx context.Context,
	path string,
	srv lsp.Server,
) (*bun.DB, error) {
	db, err := sql.Open(sqliteshim.ShimName, path)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	bu := bun.NewDB(db, sqlitedialect.New())
	_, err = bu.NewCreateTable().Model((*Selector)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	_, err = bu.NewCreateTable().Model((*HTML)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	bu.AddQueryHook(srv)
	return bu, nil
}

// Selector is a struct for a selector
type Selector struct {
	bun.BaseModel `bun:"table:selectors,alias:s"`
	// ID is the id of the selector
	ID int64 `bun:"id,pk,autoincrement"`
	// Selector is the selector for the selector
	Selector string `bun:"selector"`
	// URL is the url for the selector
	URL string `bun:"url"`
	// Context is the context for the selector
	Context string `bun:"context"`
}

type HTML struct {
	bun.BaseModel `bun:"table:html,alias:h"`
	// ID is the id of the html
	ID int64 `bun:"id,pk,autoincrement"`
	// HTML is the html for the html
	HTML string `bun:"html"`
	// URL is the url for the html
	URL string `bun:"url"`
}
