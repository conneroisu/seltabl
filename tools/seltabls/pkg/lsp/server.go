package lsp

// LineRange returns a range of a line in a document
//
// line is the line number
//
// start is the start character of the range
//
// end is the end character of the range
func LineRange(line, start, end int) Range {
	return Range{
		Start: Position{
			Line:      line,
			Character: start,
		},
		End: Position{
			Line:      line,
			Character: end,
		},
	}
}
