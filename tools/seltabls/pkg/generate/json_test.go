package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewJSONPrompt tests the NewJSONPrompt function
func TestNewJSONPrompt(t *testing.T) {
	a := assert.New(t)
	content := `{
	"fields": [
		{
			"name": "A",
			"type": "string",
			"description": "A description of the field",
			"header-selector": "tr:nth-child(1) td:nth-child(1)",
			"data-selector": "tr td:nth-child(1)",
			"control-selector": "$text",
			"must-be-present": "NCAA Codes"
		},
		{
			"name": "B",
			"type": "int",
			"description": "A description of the field",
			"header-selector": "tr:nth-child(1) td:nth-child(2)",
			"data-selector": "tr td:nth-child(2)",
			"control-selector": "$text",
			"must-be-present": "NCAA Codes"
		}
	]
}`
	got, err := NewJSONPrompt(content)
	a.NoError(err)
	a.NotEqual(got, "")

	t.Logf("json: %s", got)
}
