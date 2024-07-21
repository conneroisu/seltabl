package lsp

// Error is a struct for the error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorCode represents the error codes defined by JSON-RPC and LSP.
type ErrorCode int

const (
	// JSON-RPC Error Codes

	// CodeParseError is the error code for when a parse error occurs
	CodeParseError ErrorCode = -32700

	// CodeInvalidRequest is the error code for when an invalid request is sent
	CodeInvalidRequest ErrorCode = -32600
	// CodeMethodNotFound is the error code for when a method is not found
	CodeMethodNotFound ErrorCode = -32601
	// CodeInvalidParams is the error code for when invalid parameters are passed
	CodeInvalidParams ErrorCode = -32602
	// CodeInternalError is the error code for when an internal error occurs
	CodeInternalError ErrorCode = -32603
	// CodeJsonrpcReservedErrorRangeStart is the start of the JSON-RPC error codes
	CodeJsonrpcReservedErrorRangeStart ErrorCode = -32099
	// CodeServerErrorStart is the start
	CodeServerErrorStart ErrorCode = CodeJsonrpcReservedErrorRangeStart
	// CodeServerNotInitialized is the error code for when a server is not initialized
	CodeServerNotInitialized ErrorCode = -32002
	// CodeUnknownErrorCode is the error code for when an unknown error occurs
	CodeUnknownErrorCode ErrorCode = -32001
	// CodeJsonrpcReservedErrorRangeEnd is the end of the JSON-RPC error codes
	CodeJsonrpcReservedErrorRangeEnd ErrorCode = -32000
	// CodeServerErrorEnd is the end of the JSON-RPC error codes
	CodeServerErrorEnd ErrorCode = CodeJsonrpcReservedErrorRangeEnd
	// CodeLspReservedErrorRangeStart is the start of the LSP error codes
	CodeLspReservedErrorRangeStart ErrorCode = -32899
	// CodeRequestFailed is the error code for when a request fails
	CodeRequestFailed ErrorCode = -32803
	// CodeServerCancelled is the error code for when a server is cancelled
	CodeServerCancelled ErrorCode = -32802
	// CodeContentModified is the error code for when the content of a document is modified
	CodeContentModified ErrorCode = -32801
	// CodeRequestCancelled is the error code for when a request is cancelled
	CodeRequestCancelled ErrorCode = -32800
	// CodeLspReservedErrorRangeEnd is the end of the LSP error codes
	CodeLspReservedErrorRangeEnd ErrorCode = -32800
)
