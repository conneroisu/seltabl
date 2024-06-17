package data

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

// SetupMasterDatabase sets up the database schema
func SetupMasterDatabase(ctx context.Context, getenv func(string) string, db *bun.DB) error {
	var err error
	_, err = db.NewCreateTable().Model((*HTML)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to html table: %w", err)
	}
	_, err = db.NewCreateTable().Model((*URL)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to urls table: %w", err)
	}
	_, err = db.NewCreateTable().Model((*Selector)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to selectors table: %w", err)
	}
	return nil
}
