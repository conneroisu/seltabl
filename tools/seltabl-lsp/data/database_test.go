package data

import (
	"context"
	"testing"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/data/master"
	"github.com/stretchr/testify/assert"
)

// TestSetupMasterDbSchema tests the SetupMasterDatabase function
func TestSetupMasterDbSchema(t *testing.T) {
	db, err := NewDb(
		context.Background(),
		master.New,
		&Config{
			Schema:   master.MasterSchema,
			URI:      "sqlite://urls.sqlite",
			FileName: "urls.sqlite",
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}
