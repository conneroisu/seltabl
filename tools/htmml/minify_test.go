package htmml

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
	"github.com/yosssi/gohtml"
)

// TestClean tests the clean function
func TestClean(t *testing.T) {
	t.Run("Test Clean", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		err := run(
			ctx,
			[]string{"./clean", "https://stats.ncaa.org/teams/572260"},
		)
		assert.NoError(t, err)
	})
}

// run takes a URL as input and outputs a cleaned up HTML file.
func run(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("missing url")
	}
	// take the url from the command line
	url := args[1]
	fmt.Println("url:", url)
	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	sanitized := colly.SanitizeFileName(url)
	// create a new file
	f, err := os.Create(sanitized)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()
	content, err := cleanHTML(body)
	bytes := []byte(gohtml.Format(content))
	// write the body to the file
	_, err = f.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	fmt.Println("file written n-bytes: ", len(bytes), " bytes")
	return nil
}
