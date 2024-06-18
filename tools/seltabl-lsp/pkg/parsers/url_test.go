package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSplit tests the ExtractUrls function
func TestSplit(t *testing.T) {
	t.Run("ExtractUrls with no separator", func(t *testing.T) {
		t.Parallel()
		input := `
		// @url: https://asdf.com
		`
		expected := []string{"https://asdf.com"}
		result, err := ExtractUrls(input)
		assert.Nil(t, err)
		for _, line := range expected {
			assert.Contains(t, result, line)
		}
	})

	t.Run("ExtractUrls with a separator", func(t *testing.T) {
		t.Parallel()
		input := `
		// @url: https://elon.com
		asdf
		// @url: https://musk.com
		https://mars.com
		`
		expected := []string{"https://elon.com", "https://musk.com"}
		unexpectd := []string{"https://mars.com"}
		result, err := ExtractUrls(input)
		assert.Nil(t, err)
		assert.Equal(t, expected, result)
		for _, line := range unexpectd {
			assert.NotContains(t, result, line)
		}
	})

	t.Run("ExtractUrls with a separator and no url", func(t *testing.T) {
		t.Parallel()
		args := []struct {
			name      string
			input     string
			wantErr   bool
			expected  []string
			unexpectd []string
		}{
			{
				name: "Split with no separator",
				input: `
				//@url: https://elon.com
				asdf
				// @url: https://musk.com
				https://mars.com
				`,
				expected:  []string{"https://elon.com", "https://musk.com"},
				unexpectd: []string{"https://mars.com"},
			},
			{
				name: "markdown split",
				input: `
				# Title
				## Subtitle
				### Subsubtitle
				#### Subsubsubtitle
				##### Subsubsubsubtitle
				###### Subsubsubsubsubtitle 
				
				# @url: https://rocks.com
				
				// @url: https://elon.com
				asdf
				//@url: https://musk.com
				https://mars.com
				`,
				expected:  []string{"https://elon.com", "https://musk.com"},
				unexpectd: []string{"https://mars.com", "https://rocks.com"},
			},
			{
				name: "gibberish split",
				input: `
				// @url: https://elon.com
				lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in https://mars.com
				// @url: https://musk.com
				`,
				expected:  []string{"https://elon.com", "https://musk.com"},
				unexpectd: []string{"https://mars.com"},
			},
		}
		for _, tt := range args {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result, err := ExtractUrls(tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"ExtractUrls() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				assert.Equal(t, tt.expected, result)
				for _, line := range tt.unexpectd {
					assert.NotContains(t, result, line)
				}
			})
		}
	})

}
