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

// HandlerFunc is a function that handles a request by returning a function
type HandlerFunc func(ResponseWriter, *Request)

// Handler is a function that handles a request
type Handler interface {
	ServeRPC(ResponseWriter, *Request)
}

// ServeHTTP serves a request using the handlerfunc
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}

// ResponseWriter is a writer for a response
type ResponseWriter interface {
	Write(data []byte) (int, error)
}

// Route is a interface for routing requests and notifications to a handler
type Route interface {
	Handler(method string, request interface{}) error
}

// Router is a struct for routing requests and notifications to handlers
type Router struct {
	handlers map[string]Handler
}
