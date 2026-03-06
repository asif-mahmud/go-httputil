package validator_test

import (
	"context"
	"testing"

	"github.com/asif-mahmud/go-httputil/validator"
	vd "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type C struct {
	C1 int `validate:"gte=1"`
}

type D struct {
	D1 int `validate:"gt=0"`
}

type A struct {
	B     int      `json:"b"     validate:"gt=0"`
	C     *C       `             validate:"required"`
	Ds    []D      `json:"allDs" validate:"required,dive"`
	MinDs []D      `json:"minDs" validate:"required,min=1,dive"`
	Books []string `json:"books" validate:"required,dive,min=10"`
}

func TestFormatErrors(t *testing.T) {
	a := A{
		0,
		&C{0},
		[]D{
			{0},
			{1},
		},
		nil,
		[]string{"short", "this is long enough", "short too"},
	}

	expectedMsg := map[string]any{
		"b": "b must be greater than 0",
		"C": map[string]any{"C1": "C1 must be 1 or greater"},
		"allDs": []any{
			map[string]any{"D1": "D1 must be greater than 0"},
		},
		"minDs": "minDs is a required field",
		"books": []any{
			"books[0] must be at least 10 characters in length",
			nil,
			"books[2] must be at least 10 characters in length",
		},
	}

	err := validator.ValidateStruct(context.Background(), a)

	assert.NotNil(t, err)
	assert.IsType(t, vd.ValidationErrors{}, err)

	vErr := err.(vd.ValidationErrors)

	errMsg := validator.FormatErrors(vErr)

	assert.NotNil(t, errMsg)
	assert.Equal(t, expectedMsg, errMsg)
}
