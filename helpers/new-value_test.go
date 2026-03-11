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

func TestNewValue_WithInt(t *testing.T) {
	nzInt, err := helpers.NewValue(10)
	assert.Nil(t, err)
	assert.NotNil(t, nzInt.Interface())
	_, ok := nzInt.Interface().(*int)
	assert.True(t, ok)
}

func TestNewValue_WithSlice(t *testing.T) {
	type Data struct{ Value string }
	var expectedSlice []*Data
	nzSlice, err := helpers.NewValue(expectedSlice)
	assert.Nil(t, err)
	assert.NotNil(t, nzSlice.Interface())
	_, ok := nzSlice.Interface().(*[]*Data)
	assert.True(t, ok)
}

func TestNewValue_WithNil(t *testing.T) {
	_, err := helpers.NewValue(nil)
	assert.NotNil(t, err)
	assert.Equal(t, "unsupported type: nil interface", err.Error())
}
