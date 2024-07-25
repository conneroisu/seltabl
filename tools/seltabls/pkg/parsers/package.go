package parsers

import (
	"fmt"
	"net/url"
	"unicode"
	"unicode/utf8"
)

// ValidatePackageName checks if the given string is a valid Go package name.
func ValidatePackageName(name string) error {
	if name == "" {
		return fmt.Errorf("package name cannot be empty")
	}

	// Check if the first character is a letter
	firstChar, _ := utf8.DecodeRuneInString(name)
	if !unicode.IsLetter(firstChar) {
		return fmt.Errorf("package name must start with a letter")
	}

	// Check the rest of the characters
	for _, ch := range name[1:] {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return fmt.Errorf("package name must only contain letters, digits, and underscores")
		}
	}

	// Check if it's a Go keyword
	keywords := map[string]bool{
		"break": true, "default": true, "func": true, "interface": true, "select": true,
		"case": true, "defer": true, "go": true, "map": true, "struct": true,
		"chan": true, "else": true, "goto": true, "package": true, "switch": true,
		"const": true, "fallthrough": true, "if": true, "range": true, "type": true,
		"continue": true, "for": true, "import": true, "return": true, "var": true,
	}

	if keywords[name] {
		return fmt.Errorf("package name cannot be a Go keyword")
	}

	return nil
}

// ValidateFileName checks if the given string is a valid file name.
func ValidateFileName(name string) error {
	if name == "" {
		return fmt.Errorf("file name cannot be empty")
	}

	// Check if the first character is a letter
	firstChar, _ := utf8.DecodeRuneInString(name)
	if !unicode.IsLetter(firstChar) {
		return fmt.Errorf("file name must start with a letter")
	}

	// Check the rest of the characters
	for _, ch := range name[1:] {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return fmt.Errorf("file name must only contain letters, digits, and underscores")
		}
	}
	return nil
}

// ValidateURL checks if the given string is a valid URL.
func ValidateURL(uurl string) error {
	_, err := url.Parse(uurl)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}
	return nil
}
