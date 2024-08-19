package parsers

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	urlPattern         = `@url:\s*(https?://[^\s]+)`
	ignorePattern      = `@ignore:\s*(.*)`
	occurrencesPattern = `@occurrences:\s*(\d+)`
)

var (
	regUrlPattern         = regexp.MustCompile(urlPattern)
	regIgnorePattern      = regexp.MustCompile(ignorePattern)
	regOccurrencesPattern = regexp.MustCompile(occurrencesPattern)
)

// ParseStructComments parses the comments from struct type declarations in the
// provided Go source code and extracts @url and @ignore-elements into separate
// arrays.
func parseStructComments(commentSrc string) (
	url string,
	ignoreElements []string,
	occurrences int,
	err error,
) {
	// Apply the regular expressions to the comment source
	url = regUrlPattern.FindString(commentSrc)

	ignoreElements = regIgnorePattern.FindAllString(commentSrc, -1)

	occurrencesStr := regOccurrencesPattern.FindString(commentSrc)
	if occurrencesStr == "" {
		occurrencesStr = "1"
	}
	occurrences, err = strconv.Atoi(occurrencesStr)
	if err != nil {
		return "", nil, 0, fmt.Errorf("failed to parse occurrences: %w", err)
	}

	return url, ignoreElements, occurrences, nil
}
