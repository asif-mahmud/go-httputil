package helpers_test

import (
	"encoding/json"
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

func TestNewSliceValue(t *testing.T) {
	type testStruct struct {
		Id   int
		Name string
	}

	var expectedZero []testStruct

	nz, err := helpers.NewValue(expectedZero)

	assert.Nil(t, err)

	nzi := nz.Interface()

	assert.NotNil(t, nzi)

	nzvp, ok := nzi.(*[]testStruct)

	assert.True(t, ok)

	data := `[{"Id":1,"Name":"1"}]`
	expected := testStruct{Id: 1, Name: "1"}

	e := json.Unmarshal([]byte(data), nzvp)

	assert.Nil(t, e)
	assert.Len(t, *nzvp, 1)
	assert.Equal(t, expected, (*nzvp)[0])
}
