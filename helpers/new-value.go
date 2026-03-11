package helpers

import (
	"errors"
	"reflect"
)

func NewValue(v any) (reflect.Value, error) {
	if v == nil {
		return reflect.Value{}, errors.New("unsupported type: nil interface")
	}

	vt := reflect.TypeOf(v)

	if vt.Kind() == reflect.Pointer {
		return reflect.New(vt.Elem()), nil
	}

	return reflect.New(vt), nil
}
