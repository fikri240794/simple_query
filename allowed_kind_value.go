package simple_query

import "reflect"

var allowedKindValue map[reflect.Kind]bool = map[reflect.Kind]bool{
	reflect.Array:   true,
	reflect.Bool:    true,
	reflect.Float32: true,
	reflect.Float64: true,
	reflect.Int:     true,
	reflect.Int8:    true,
	reflect.Int16:   true,
	reflect.Int32:   true,
	reflect.Int64:   true,
	reflect.Slice:   true,
	reflect.String:  true,
	reflect.Uint:    true,
	reflect.Uint8:   true,
	reflect.Uint16:  true,
	reflect.Uint32:  true,
	reflect.Uint64:  true,
}
