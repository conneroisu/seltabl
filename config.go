package seltabl

import (
	"reflect"
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
	return cfg
}
