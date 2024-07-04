package lsp

import "github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"

// Actor is an interface for an actor
type Actor interface {
	Method() methods.Method
}
