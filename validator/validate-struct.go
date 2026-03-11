package validator

import (
	"context"
	"reflect"
)

// ValidateStruct validates the incoming structure, and dives into slice elements.
func ValidateStruct(ctx context.Context, s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		return validate.StructCtx(ctx, s)
	case reflect.Slice, reflect.Array:
		return validate.VarCtx(ctx, s, "dive")
	default:
		// Not a struct or slice of structs, no inner validation to perform on root types lacking tags
		return nil
	}
}
