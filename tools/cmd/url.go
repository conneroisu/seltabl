package cmd

import (
	"bufio"
	"regexp"
	"strings"
	"testing"
)

const (
	commentPrefix = "//"
	commentSuffix = "\n"

	urlComment = "@url:"
)

// ExtractUrls splits a string into a slice of strings of only the comments in a given file's content's
// url comments.
func ExtractUrls(fileContent string, t *testing.T) ([]string, error) {
	var result []string
	result = make([]string, 0)
	t.Logf("Splitting file content: %s", fileContent)
	// Define the regex to match URL comments
	urlRegex := regexp.MustCompile(`// @url:\s*(https?://[^\s]+)`)
	urlRegex2 := regexp.MustCompile(`//@url:\s*(https?://[^\s]+)`)
	// Scan the file line by line
	scanner := bufio.NewScanner(strings.NewReader(fileContent))
	for scanner.Scan() {
		line := scanner.Text()
		// Find the URL in the line
		matches := urlRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			t.Logf("Found URL: %s", matches[1])
			result = append(result, matches[1])
		}
		matches = urlRegex2.FindStringSubmatch(line)
		if len(matches) > 0 {
			t.Logf("Found URL: %s", matches[1])
			result = append(result, matches[1])
		}
	}
	return result, nil
}
