package generate

import (
	"fmt"
	"net/url"

	"github.com/conneroisu/seltabl/tools/seltabls/domain"
)

// isvalidgotype checks if the given type is a valid go type
func isvalidgotype(t string) bool {
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

// isValidTreeWidth checks if the given tree width is valid
func isValidTreeWidth(treeWidth int) bool {
	if treeWidth%2 != 0 || treeWidth < 1 {
		return false
	}
	return true
}

// isValidTreeDepth checks if the given tree depth is valid
func isValidTreeDepth(treeDepth int) bool {
	if treeDepth < 1 {
		return false
	}
	return true
}

// validateConfig validates the given config file
func validateConfig(cfg *domain.ConfigFile) error {
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
	if cfg.Recycle {
		return fmt.Errorf("recycle is required")
	}
	return nil
}
