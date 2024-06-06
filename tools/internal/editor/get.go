package editor

import (
	"os"
	"strings"
)

// defaultEditor is the default editor to use
const defaultEditor = "nvim"

// GetEditor returns the editor and the arguments to pass to it
func GetEditor() (string, []string) {
	editor := strings.Fields(os.Getenv("EDITOR"))
	if len(editor) > 1 {
		return editor[0], editor[1:]
	}
	if len(editor) == 1 {
		return editor[0], []string{}
	}
	return defaultEditor, []string{}
}
