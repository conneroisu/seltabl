package seltabl

import (
	"reflect"
)

var (
	// cSels is a list of supported control selectors
	cSels = []string{ctlInnerTextSelector, ctlAttrSelector}
)

const (
	// selectorDataTag is the tag used to mark a data cell.
	selectorDataTag = "dSel"
	// selectorHeaderTag is the tag used to mark a header selector.
	selectorHeaderTag = "hSel"
	// selectorTag is the tag used to mark a selector.
	selectorQueryTag = "qSel"
	// selectorMustBePresentTag is the tag used to mark text that must be present in a given content.
	selectorMustBePresentTag = "must"

	// selectorControlTag is the tag used to signify selecting aspects of a cell
	selectorControlTag = "ctl"

	// cSelInnerTextSelector is the selector used to extract text from a cell.
	ctlInnerTextSelector = "text"
	// cSelAttrSelector is the selector used to extract attributes from a cell.
	ctlAttrSelector = "query"
)

// SelectorConfig is a struct for configuring a selector
type SelectorConfig struct {
	DataSelector  string // selector for the data cell
	HeadSelector  string // selector for the header cell
	QuerySelector string // selector for the data cell
	ControlTag    string // tag used to signify selecting aspects of a cell
	MustBePresent string // text that must be present in a given content
}

// NewSelectorConfig parses a struct tag and returns a SelectorConfig
func NewSelectorConfig(tag reflect.StructTag) *SelectorConfig {
	cfg := &SelectorConfig{
		HeadSelector:  tag.Get(selectorHeaderTag),
		DataSelector:  tag.Get(selectorDataTag),
		QuerySelector: tag.Get(selectorQueryTag),
		ControlTag:    tag.Get(selectorControlTag),
		MustBePresent: tag.Get(selectorMustBePresentTag),
	}
	if cfg.QuerySelector == "" || cfg.DataSelector == ctlAttrSelector {
		cfg.QuerySelector, cfg.ControlTag =
			ctlInnerTextSelector,
			ctlInnerTextSelector
	}
	return cfg
}
