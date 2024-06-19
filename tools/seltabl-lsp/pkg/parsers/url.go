package parsers

import (
	"bufio"
	"context"
	"fmt"
	"regexp"
	"strings"
)

var (
	// urlRegex is a regex for finding urls in a comment with a space
	urlRegex = regexp.MustCompile(`// @url:\s*(https?://[^\s]+)`)
	// urlRegex2 is a regex for finding urls in a comment without a space
	urlRegex2 = regexp.MustCompile(`//@url:\s*(https?://[^\s]+)`)
)

// ExtractUrls collects all the urls with comments using the // @url: syntax
// and returns a slice of strings of only the urls.
//
// It also returns an error if there is an error while scanning the file.
func ExtractUrls(_ context.Context, src string) ([]string, error) {
	var result []string
	result = make([]string, 0)
	// Define the regex to match URL comments
	// Scan the file line by line
	scanner := bufio.NewScanner(strings.NewReader(src))
	for scanner.Scan() {
		line := scanner.Text()
		matches := urlRegex.FindStringSubmatch(line)
		if len(matches) > 0 && matches != nil {
			result = append(result, matches[1])
		}
		matches = urlRegex2.FindStringSubmatch(line)
		if len(matches) > 0 && matches != nil {
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
