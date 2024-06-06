package editor

import "fmt"

// Option defines an editor option.
//
// An Option may act differently in some editors, or not be supported in
// some of them.
type Option func(editor, filename string) (args []string, pathInArgs bool)

// OpenAtLine opens the file at the given line number in supported editors.
func OpenAtLine(number uint) Option {
	plusLineEditors := []string{"vi", "vim", "nvim", "nano", "emacs", "kak", "gedit"}
	return func(editor, filename string) ([]string, bool) {
		for _, e := range plusLineEditors {
			if editor == e {
				return []string{fmt.Sprintf("+%d", number)}, false
			}
		}
		if editor == "code" {
			return []string{
				"--goto",
				fmt.Sprintf("%s:%d", filename, number),
			}, true
		}
		return nil, false
	}
}
