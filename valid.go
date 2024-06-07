package seltabl

import "reflect"

// isTypeSupported checks if the type is supported
func isTypeSupported(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.String:
		return true
	case reflect.Int:
		return true
	case reflect.Int8:
		return true
	case reflect.Int16:
		return true
	case reflect.Int32:
		return true
	case reflect.Int64:
		return true
	case reflect.Uint:
		return true
	case reflect.Uint8:
		return true
	case reflect.Uint16:
		return true
	case reflect.Uint32:
		return true
	case reflect.Uint64:
		return true
	case reflect.Float32:
		return true
	case reflect.Float64:
		return true
	case reflect.Struct:
		return false
	default:
		return false
	}
}
