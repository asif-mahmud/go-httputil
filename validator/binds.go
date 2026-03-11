package validator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

// BindUrlValues binds url.Values into a struct instance.
// Additionally it runs the struct through mold transformer.
func BindUrlValues(ctx context.Context, v url.Values, s any) error {
	if err := formDecoder.Decode(s, v); err != nil {
		return err
	}

	return runMold(ctx, s)
}

// BindJSON binds body into a struct instance.
// Additionally it runs the struct through mold transformer.
func BindJSON(ctx context.Context, body io.ReadCloser, s any) error {
	defer body.Close()
	if err := json.NewDecoder(body).Decode(s); err != nil {
		return err
	}

	return runMold(ctx, s)
}

// BindPathValues binds values from the URL path to the struct fields
// based on the "path" tag. It uses reflection to dynamically set the field values.
func BindPathValues(ctx context.Context, r *http.Request, s any) error {
	val := reflect.ValueOf(s)
	// If s is a pointer, dereference it to access the actual struct.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Skip fields that cannot be set (e.g., unexported fields).
		if !field.CanSet() {
			continue
		}

		structField := typ.Field(i)

		pathTag := structField.Tag.Get("path")

		// Skip fields that don't have a "path" tag.
		if pathTag == "" {
			continue
		}

		pathValue := r.PathValue(pathTag)

		// If the path parameter is empty, skip this field.
		if pathValue == "" {
			continue
		}

		// Set the field value using the path value.
		if err := setFieldValue(field, pathValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", structField.Name, err)
		}
	}

	return runMold(ctx, s)
}

// runMold applies mold transformations to struct values, recursing into slices if needed
func runMold(ctx context.Context, s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		return conform.Struct(ctx, s)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if elem.Kind() != reflect.Pointer {
				if elem.CanAddr() {
					elem = elem.Addr()
				} else {
					// Cannot take address, skip or handle (usually pointers are stored)
					continue
				}
			}

			// Unwrap element to check if it's a struct
			elemVal := elem
			if elemVal.Kind() == reflect.Pointer {
				elemVal = elemVal.Elem()
			}
			if elemVal.Kind() == reflect.Struct {
				if err := conform.Struct(ctx, elem.Interface()); err != nil {
					return err
				}
			}
		}
		return nil
	default:
		return nil
	}
}

// setFieldValue sets the value of a struct field based on the string value.
// It converts the string value to the appropriate type based on the field's type.
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse int: %w", err)
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %w", err)
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, field.Type().Bits())
		if err != nil {
			return fmt.Errorf("failed to parse float: %w", err)
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("failed to parse bool: %w", err)
		}
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}
