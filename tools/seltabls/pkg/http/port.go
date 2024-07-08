package http

import (
	"os/user"
)

// HomeDir returns the home directory of the user that owns the current process.
//
// Example:
//
//	fmt.Println(HomeDir())
//	// Output:
//	/home/user
func HomeDir() string {
	if homeDirSet {
		return homeDir
	}
	if user, err := user.Current(); err == nil {
		homeDir = user.HomeDir
	}
	homeDirSet = true
	return homeDir
}
