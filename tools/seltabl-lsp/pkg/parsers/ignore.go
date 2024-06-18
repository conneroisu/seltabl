// Package parsers provides a set of functions for parsing html documents.
package parsers

import (
	"bufio"
	"context"
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
func ExtractIgnores(_ context.Context, fileContent string) ([]string, error) {
	var result []string
	var err error
	scanner := bufio.NewScanner(strings.NewReader(fileContent))
	for scanner.Scan() {
		var matches []string
		line := scanner.Text()
		line = strings.TrimSpace(line)
		// Find the ignore elements in the line
		matches = ignoreRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			res := matches[1]
			split := strings.Split(res, ", ")
			if !contains(split, "") && !contains(split, res) {
				result = append(result, split...)
			}
		}
		matches = ignoreRegex2.FindStringSubmatch(line)
		if len(matches) > 0 {
			res := matches[1]
			split := strings.Split(res, ", ")
			if !contains(split, "") && !contains(split, res) {
				result = append(result, split...)
			}
		}
	}
	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}
	if len(result) == 0 {
		return nil, ErrNoIgnoresFound{File: fileContent}
	}
	return result, nil
}
