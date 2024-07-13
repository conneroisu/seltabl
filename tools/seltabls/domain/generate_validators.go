package domain

import (
	"fmt"
	"reflect"
)

// ValidateConfig validates the given config file.
func ValidateConfig(cfg *PromptFile) error {
	if cfg.URL == "" {
		return fmt.Errorf("url is required")
	}
	if cfg.IgnoreElements == nil {
		return fmt.Errorf("ignore-elements is required")
	}
	if cfg.Selectors == nil {
		return fmt.Errorf("selectors is required")
	}
	if cfg.RuledHTMLBody == "" {
		return fmt.Errorf("numbered-html-body is required")
	}
	if cfg.SmartModel == "" {
		return fmt.Errorf("smart-model is required")
	}
	if cfg.FastModel == "" {
		return fmt.Errorf("fast-model is required")
	}
	return nil
}

// Validate checks if the sections are valid.
func (s *SectionsResponse) Validate() error {
	for i, section := range s.Sections {
		v := reflect.ValueOf(section)
		t := reflect.TypeOf(section)
		for j := 0; j < v.NumField(); j++ {
			field := v.Field(j)
			fieldName := t.Field(j).Name
			fieldType := field.Type().Kind()
			switch fieldType {
			case reflect.String:
				if field.String() == "" {
					return fmt.Errorf("section %d has no %s", i, fieldName)
				}
			case reflect.Int:
				if fieldName == "Start" && field.Int() < 0 {
					return fmt.Errorf("section %d has a negative %s", i, fieldName)
				}
				if fieldName == "End" && field.Int() < 0 {
					return fmt.Errorf("section %d has a negative %s", i, fieldName)
				}
				if section.Start >= section.End {
					return fmt.Errorf("section %d has a start greater than or equal to the end", i)
				}
			case reflect.Slice:
				if field.IsNil() || field.Len() == 0 {
					return fmt.Errorf("section %d has no %s", i, fieldName)
				}
			}
		}
	}
	return nil
}
