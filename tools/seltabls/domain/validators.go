package domain

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// IsValidGoType checks if the given type is a valid go type
func IsValidGoType(t string) bool {
	switch t {
	case "string":
		return true
	case "int":
		return true
	case "int8":
		return true
	case "int16":
		return true
	case "int32":
		return true
	case "int64":
		return true
	case "uint":
		return true
	case "uint8":
		return true
	case "uint16":
		return true
	case "uint32":
		return true
	case "uint64":
		return true
	case "float32":
		return true
	case "float64":
		return true
	case "bool":
		return true
	default:
		return false
	}
}

// IsURL checks if the given string is a valid url
func IsURL(toValidateURL string) (err error) {
	_, err = url.ParseRequestURI(toValidateURL)
	return err
}

// IsValidTreeWidth checks if the given tree width is valid
func IsValidTreeWidth(treeWidth int) bool {
	if treeWidth%2 != 0 || treeWidth < 1 {
		return false
	}
	return true
}

// ValidateConfig validates the given config file
func ValidateConfig(cfg *ConfigFile) error {
	if cfg.Name == "" {
		return fmt.Errorf("name is required")
	}
	if cfg.URL == "" {
		return fmt.Errorf("url is required")
	}
	if cfg.IgnoreElements == nil {
		return fmt.Errorf("ignore-elements is required")
	}
	if cfg.Selectors == nil {
		return fmt.Errorf("selectors is required")
	}
	if cfg.HTMLBody == "" {
		return fmt.Errorf("html-body is required")
	}
	if cfg.NumberedHTMLBody == "" {
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

// Verify checks if the selectors are in the html
func (f *Field) Verify(htmlBody string) error {
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(htmlBody),
	)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}
	if f.DataSelector != "" {
		sel := doc.Find(f.DataSelector)
		if sel.Length() == 0 {
			return fmt.Errorf("failed to find selector: %s", f.DataSelector)
		}
	} else {
		return fmt.Errorf("no data found for selector %s with type %s in field %s with type %s", f.DataSelector, f.Type, f.Name, f.Type)
	}
	if f.ControlSelector != "" {
		sel := doc.Find(f.ControlSelector)
		if sel.Length() == 0 {
			return fmt.Errorf("failed to find selector: %s", f.ControlSelector)
		}
	} else {
		return fmt.Errorf("no control found for selector %s with type %s in field %s with type %s", f.ControlSelector, f.Type, f.Name, f.Type)
	}
	if f.QuerySelector != "" {
		sel := doc.Find(f.QuerySelector)
		if sel.Length() == 0 {
			return fmt.Errorf("failed to find selector: %s", f.QuerySelector)
		}
	} else {
		return fmt.Errorf("no query found for selector %s with type %s in field %s with type %s", f.QuerySelector, f.Type, f.Name, f.Type)
	}
	if f.HeaderSelector != "" {
		sel := doc.Find(f.HeaderSelector)
		if sel.Length() == 0 {
			return fmt.Errorf("failed to find selector: %s", f.HeaderSelector)
		}
	}
	mbp := f.MustBePresent
	docTxt := doc.Text()
	if !strings.Contains(docTxt, mbp) {
		return fmt.Errorf(
			"must be present (%s) not found for field %s with type %s",
			mbp,
			f.Name,
			f.Type,
		)
	}
	return nil
}

// force type cast for Responder
var _ responder = (*IdentifyResponse)(nil)
var _ responder = (*FieldsResponse)(nil)
