package data

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

// TestNewDb tests the NewDb function
func TestSetupMasterDbSchema(t *testing.T) {
	args := []struct {
		name     string
		wantErr  bool
		fn       func(ctx context.Context, getenv func(string) string, db *bun.DB) error
		expected []string
	}{
		{
			name:     "Test Master DB Schema",
			fn:       SetupMasterDatabase,
			expected: []string{"urls", "html"},
		},
	}
	for _, tt := range args {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			getenv := func(string) string {
				return "true"
			}
			db, err := NewDb(
				ctx,
				getenv,
				":memory:",
				tt.fn,
				sqlitedialect.New(),
			)
			if err != nil {
				t.Fatalf("failed to create database: %v", err)
			}
			tables := []string{}
			err = db.NewRaw(
				"SELECT name FROM sqlite_master WHERE type='table'",
			).Scan(ctx, &tables)
			if err != nil {
				t.Fatalf("failed to get URLs: %v", err)
			}
			t.Logf("tables: %v", tables)
			for _, v := range tt.expected {
				assert.Contains(t, tables, v, "expected %s not found in %s", v, tables)
			}
			cancel()
			select {
			case <-ctx.Done():
				t.Logf("context closed")
			case <-time.After(time.Second):
				t.Fatalf("context not closed")
			}
		})
	}
}
