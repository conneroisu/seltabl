package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnore(t *testing.T) {
	t.Run("Test Ignore Valid Case", func(t *testing.T) {
		t.Parallel()
		input := `
		// @ignore-elements: div, script
		`
		expected := []string{"div", "script"}
		result, err := ExtractIgnores(input)
		assert.Nil(t, err)
		for _, v := range expected {
			assert.Contains(t, result, v)
		}
	})

	t.Run("Test Ignore Invalid Case", func(t *testing.T) {
		t.Parallel()
		input := `
		// @ignore-elements: 
		`
		var expected []string
		result, err := ExtractIgnores(input)
		assert.NotNil(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Test Ignore Partially Invalid Case 2", func(t *testing.T) {
		t.Parallel()
		input := `
		// @ignore-elements: div, script, 
		`
		expected := []string{"div", "script"}
		result, err := ExtractIgnores(input)
		assert.Nil(t, err)
		assert.Equal(t, expected, result)
		assert.NotContains(t, result, "")
		assert.NotContains(t, result, " ")
	})
}
