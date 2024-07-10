package domain

import "fmt"

// ErrHTTPParse is an error for when a http request fails
type ErrHTTPParse struct {
	URL        string
	StatusCode int
}

// Error implements the error interface
func (e ErrHTTPParse) Error() string {
	return fmt.Sprintf(
		"failed to parse response frmo url: %s with status code: %d",
		e.URL,
		e.StatusCode,
	)
}

// ErrDocumentFromReader is an error for when a document cannot be created from
// a reader with the given url and content.
type ErrDocumentFromReader struct {
	URL     string
	Content string
}

// Error implements the error interface
func (e ErrDocumentFromReader) Error() string {
	return fmt.Sprintf(
		"failed to create document from reader with url: %s and content: \n```html\n%s\n```",
		e.URL,
		e.Content,
	)
}
