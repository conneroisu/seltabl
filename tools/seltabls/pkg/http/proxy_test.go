package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStartContainer tests the StartContainer function
func TestStartContainer(t *testing.T) {
	ctx := context.Background()
	container, err := StartContainer(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, container)
	assert.NotEmpty(t, container.URI)
	client := &http.Client{}
	resp, err := client.Get(container.URI)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, 200)
}
