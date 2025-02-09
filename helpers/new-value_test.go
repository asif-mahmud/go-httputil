package helpers_test

import (
	"testing"

	"github.com/asif-mahmud/go-httputil/helpers"
	"github.com/stretchr/testify/assert"
)

func TestNewValue(t *testing.T) {
	type testStruct struct {
		Id   int
		Name string
	}

	var expectedZero testStruct

	nz, err := helpers.NewValue(expectedZero)

	assert.Nil(t, err)

	nzi := nz.Interface()

	assert.NotNil(t, nzi)

	nzv, ok := nzi.(*testStruct)

	assert.True(t, ok)

	nzv.Id = 1
	nzv.Name = "1"

	assert.NotEqual(t, expectedZero, *nzv)
}

func TestNewPtrValue(t *testing.T) {
	type testStruct struct {
		Id   int
		Name string
	}

	var expectedZero testStruct

	nz, err := helpers.NewValue(&expectedZero)

	assert.Nil(t, err)

	nzi := nz.Interface()

	assert.NotNil(t, nzi)

	nzv, ok := nzi.(*testStruct)

	assert.True(t, ok)

	nzv.Id = 1
	nzv.Name = "1"

	assert.NotEqual(t, expectedZero, *nzv)
}
