package parsers

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var (
	// ignoreRegex is a regex for finding ignore elements in a comment with a space
	ignoreRegex = regexp.MustCompile(`// @ignore-elements: (.*)`)
	// ignoreRegex2 is a regex for finding ignore elements in a comment without a space
	ignoreRegex2 = regexp.MustCompile(`//@ignore-elements: (.*)`)
)

// ErrNoIgnoresFound is an error for when no ignores are found
type ErrNoIgnoresFound struct {
	File string
}

// Error implements the error interface
func (e ErrNoIgnoresFound) Error() string {
	return fmt.Sprintf("no ignores found in file: %s", e.File)
}

// ExtractIgnores extracts the ignores from the given text
func ExtractIgnores(fileContent string) ([]string, error) {
	var result []string
	scanner := bufio.NewScanner(strings.NewReader(fileContent))
	for scanner.Scan() {
		line := scanner.Text()
		// Find the ignore elements in the line
		matches := ignoreRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			res := matches[1]
			split := strings.Split(res, ", ")
			if !contains(split, "") && !contains(split, res) {
				result = append(result, res)
			}
		}
		matches = ignoreRegex2.FindStringSubmatch(line)
		if len(matches) > 0 {
			result = append(result, matches[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}
	if len(result) == 0 {
		return nil, ErrNoIgnoresFound{File: fileContent}
	}
	return result, nil
}
