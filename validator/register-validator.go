package validator

import (
	vd "github.com/go-playground/validator/v10"
)

// RegisterValidator register custiom validation function.
func RegisterValidator(tag string, fn vd.Func, callValidationEvenIfNull ...bool) {
	validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}
