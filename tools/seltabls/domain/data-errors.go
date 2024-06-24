package domain

import "fmt"

// ErrDataNilFileName is an error for when a file name is nil
type ErrDataNilFileName struct {
	Message string
}

// Error implements the error interface
func (e ErrDataNilFileName) Error() string {
	return e.Message
}

// ErrDataSchemaFailed is an error for when a database schema fails to execute
type ErrDataSchemaFailed struct {
	Message string
	Schema  string
}

// Error implements the error interface
func (e ErrDataSchemaFailed) Error() string {
	return fmt.Sprintf(
		"failed to execute database schema: %s\n%s",
		e.Message,
		e.Schema,
	)
}

// ErrDataOpenFailed is an error for when a database file fails to open
type ErrDataOpenFailed struct {
	Message  string
	FileName string
}

// Error implements the error interface
func (e ErrDataOpenFailed) Error() string {
	return fmt.Sprintf(
		"failed to open database file, '%s', with error: %s",
		e.FileName,
		e.Message,
	)
}

// ErrDataFailedLocateConfig is an error for when a database config file fails to locate
type ErrDataFailedLocateConfig struct {
	Message string
}

// Error implements the error interface
func (e ErrDataFailedLocateConfig) Error() string {
	return e.Message
}
