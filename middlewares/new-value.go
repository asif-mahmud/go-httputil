package middlewares

import (
	"errors"
	"reflect"
)

// newValue creates a pointer to a new zero value of the type of v.
func newValue(v any) (reflect.Value, error) {
	vt := reflect.TypeOf(v)

	switch vt.Kind() {
	case reflect.Pointer:
		return reflect.New(vt.Elem()), nil

	case reflect.Struct:
		return reflect.New(vt), nil

	default:
		return reflect.Value{}, errors.New("unsupported type")
	}
}
