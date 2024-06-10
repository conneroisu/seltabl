package seltabl

import (
	"reflect"
)

var (
	cSels = []string{cSelInnerTextSelector, cSelAttrSelector}
)

const (
	cSelInnerTextSelector    = "$text"   // cSelInnerTextSelector is the selector used to extract text from a cell.
	cSelAttrSelector         = "$query"  // cSelAttrSelector is the selector used to extract attributes from a cell.
	headerTag                = "seltabl" // headerTag is the tag used to mark a header cell.
	selectorDataTag          = "dSel"    // selectorDataTag is the tag used to mark a data cell.
	selectorHeaderTag        = "hSel"    // selectorHeaderTag is the tag used to mark a header selector.
	selectorControlTag       = "cSel"    // selectorControlTag is the tag used to mark a data selector.
	selectorQueryTag         = "qSel"    // selectorTag is the tag used to mark a selector.
	selectorMustBePresentTag = "must"    // selectorMustBePresentTag is the tag used to mark a selector.
)

// SelectorConfig is a struct for configuring a selector
type SelectorConfig struct {
	HeadName      string
	DataSelector  string
	HeadSelector  string
	QuerySelector string
	ControlTag    string
	MustBePresent string
}

// NewSelectorConfig parses a struct tag and returns a SelectorConfig
func NewSelectorConfig(tag reflect.StructTag) *SelectorConfig {
	cfg := &SelectorConfig{
		HeadName:      tag.Get(headerTag),
		HeadSelector:  tag.Get(selectorHeaderTag),
		DataSelector:  tag.Get(selectorDataTag),
		QuerySelector: tag.Get(selectorQueryTag),
		ControlTag:    tag.Get(selectorControlTag),
		MustBePresent: tag.Get(selectorMustBePresentTag),
	}
	if cfg.QuerySelector == "" || cfg.DataSelector == cSelAttrSelector {
		cfg.QuerySelector, cfg.ControlTag = cSelInnerTextSelector, cSelInnerTextSelector
	}
	return cfg
}
