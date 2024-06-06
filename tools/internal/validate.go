package internal

import (
	"fmt"
	"net/url"
	"unicode"
)

// isValidPackageName validates the package name
func isValidPackageName(val string) error {
	if val == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	// Check if the first character is a letter
	if !unicode.IsLetter(rune(val[0])) {
		return fmt.Errorf("package name must start with a letter")
	}
	// Check the remaining characters
	for _, r := range val[1:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return fmt.Errorf("package name can only contain letters and digits")
		}
	}
	return nil
}

// isValidURL validates the url
func isValidURL(val string) error {
	if val == "" {
		return fmt.Errorf("url cannot be empty")
	}
	_, err := url.ParseRequestURI(val)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	return nil
}
