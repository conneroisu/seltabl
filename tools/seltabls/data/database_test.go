package data

import (
	"context"
	"testing"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
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
			FileName: ":memory:",
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

// TestSetupMasterDbSchemaInvalid tests the SetupMasterDatabase function with an invalid schema
func TestSetupMasterDbSchemaInvalid(t *testing.T) {
	db, err := NewDb(
		context.Background(),
		master.New,
		&Config{
			Schema:   "invalid",
			URI:      "sqlite://urls.sqlite",
			FileName: ":memory:",
		},
	)
	assert.Error(t, err)
	assert.Nil(t, db)
}
