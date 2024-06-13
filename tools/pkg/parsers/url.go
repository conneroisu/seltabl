package parsers

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var (
	urlRegex  = regexp.MustCompile(`// @url:\s*(https?://[^\s]+)`)
	urlRegex2 = regexp.MustCompile(`//@url:\s*(https?://[^\s]+)`)
)

// ExtractUrls splits a string into a slice of strings of only the comments in a given file's content's
// url comments.
func ExtractUrls(fileContent string) ([]string, error) {
	var result []string
	result = make([]string, 0)
	// Define the regex to match URL comments
	// Scan the file line by line
	scanner := bufio.NewScanner(strings.NewReader(fileContent))
	for scanner.Scan() {
		line := scanner.Text()
		// Find the URL in the line
		matches := urlRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			result = append(result, matches[1])
		}
		matches = urlRegex2.FindStringSubmatch(line)
		if len(matches) > 0 {
			result = append(result, matches[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no urls found in file")
	}
	return result, nil
}
