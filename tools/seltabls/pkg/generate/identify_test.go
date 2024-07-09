package generate

import (
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

// TestNewIdentifyPrompt tests the NewIdentifyPrompt function
func TestNewIdentifyPrompt(t *testing.T) {
	a := assert.New(t)
	out, err := NewIdentifyPrompt(
		"https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
		"<html></html>",
	)
	a.NoError(err)
	a.NotEmpty(out)
	t.Logf("identify: %s", out)
}

func TestNewIdentifyAggregatePrompt(t *testing.T) {
	a := assert.New(t)
	out, err := NewIdentifyAggregatePrompt(
		"<html></html>",
		[]string{
			`{
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
			}`,
		},
	)
	a.NoError(err)
	a.NotEmpty(out)
	t.Logf("identify: %s", out)
}
