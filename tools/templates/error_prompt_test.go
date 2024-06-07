package templates

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFillErrorPrompt tests the FillErrorPrompt function
func TestFillErrorPrompt(t *testing.T) {
	err := fmt.Errorf("failed to parse html: %w", nil)
	html := "<table><tr><td>a</td></tr></table>"
	input := "input"
	fields := []string{"a"}
	output, err := FillErrorPrompt(err, html, input)
	if err != nil {
		t.Fatalf("failed to fill error prompt: %v", err)
	}
	assert.NoError(t, err)
}
