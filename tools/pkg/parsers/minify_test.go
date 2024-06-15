package parsers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
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
			t,
		)
		assert.NoError(t, err)
	})
}

// run takes a URL as input and outputs a cleaned up HTML file.
func run(ctx context.Context, args []string, t *testing.T) error {
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
	content, err := cleanHTML(
		body,
		[]string{"script", "style", "link", "img", "footer", "header"},
	)
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(
			string(content),
		),
	)
	// the length of script, style, link, img, footer, header should all be 0
	scripts := doc.Find("script")
	assert.Equal(t, 0, scripts.Length())
	styles := doc.Find("style")
	assert.Equal(t, 0, styles.Length())
	links := doc.Find("link")
	assert.Equal(t, 0, links.Length())
	images := doc.Find("img")
	assert.Equal(t, 0, images.Length())
	footers := doc.Find("footer")
	assert.Equal(t, 0, footers.Length())
	headers := doc.Find("header")
	assert.Equal(t, 0, headers.Length())
	return nil
}
