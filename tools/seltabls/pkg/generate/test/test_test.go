package test

import (
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

// TestNewTestTemplate tests the NewTestTemplate function
func TestNewTestTemplate(t *testing.T) {
	a := assert.New(t)
	ctn, err := NewTestFileContent("test.go", "TestStruct", "0.0.0.0", "main")
	a.NoError(err)
	a.NotEmpty(ctn)
	t.Logf("test: %s", ctn)
}
