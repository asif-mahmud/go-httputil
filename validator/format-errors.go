package validator

import (
	"strconv"
	"strings"

	vd "github.com/go-playground/validator/v10"
)

func setValueInMap(m map[string]interface{}, key string, value string) {
	keys := strings.Split(key, ".")
	currMap := m

	for i := 1; i < len(keys); i++ { // Start from 1 to skip "Data"
		// key := toCamelCase(keys[i])
		key := keys[i]
		last := i == len(keys)-1

		if strings.Contains(key, "[") {
			currMap = handleArrayKey(currMap, key, last, value)
		} else {
			currMap = handleMapKey(currMap, key, last, value)
		}
	}
}

func handleArrayKey(
	currMap map[string]interface{},
	key string,
	last bool,
	value string,
) map[string]interface{} {
	parts := strings.Split(key, "[")
	key = parts[0]
	index, _ := strconv.Atoi(strings.TrimRight(parts[1], "]"))

	if currMap[key] == nil {
		currMap[key] = make([]interface{}, index+1)
	} else if len(currMap[key].([]interface{})) <= index {
		currMap[key] = append(currMap[key].([]interface{}), make([]interface{}, index-len(currMap[key].([]interface{}))+1)...)
	}

	if last {
		currMap[key].([]interface{})[index] = value
	} else {
		if currMap[key].([]interface{})[index] == nil {
			currMap[key].([]interface{})[index] = make(map[string]interface{})
		}
		currMap = currMap[key].([]interface{})[index].(map[string]interface{})
	}

	return currMap
}

func handleMapKey(
	currMap map[string]interface{},
	key string,
	last bool,
	value string,
) map[string]interface{} {
	if currMap[key] == nil {
		if last {
			currMap[key] = value
		} else {
			currMap[key] = make(map[string]interface{})
		}
	}

	if !last {
		currMap = currMap[key].(map[string]interface{})
	}

	return currMap
}

// FormatErrors formats ValidationErrors into a map of error strings similar to
// the source data type.
//
// Example DTO -
//
//	type Address struct {
//		Street string `json:"street" validate:"required"`
//	}
//
//	type User struct {
//		Id        string    `json:"id" validate:"required"`
//		UserName  string    `json:"userName" validate:"required"`
//		Email     string    `json:"email" validate:"required,email"`
//		Addresses []Address `json:"addresses" validate:"required,dive"`
//		Books     []string  `json:"books" validate:"required,dive,min=10"`
//	}
//
//	type RequestDTO struct {
//		Data User `json:"data" validate:"required"`
//	}
//
// Example validation error -
//
//	{
//		"data.addresses[0].street":"street is a required field",
//		"data.books[0]":"books[0] must be at least 10 characters in length",
//		"data.books[1]":"books[1] must be at least 10 characters in length",
//		"data.books[2]":"books[2] is a required field",
//		"data.email":"email is a required field",
//		"data.id":"id is a required field",
//		"data.userName":"userName is a required field"
//	}
//
// Corresponding error message after formatting -
//
//	{
//	  "addresses": [
//	    {
//	      "street": "street is a required field"
//	    }
//	  ],
//	  "books": [
//	    "books[0] must be at least 10 characters in length",
//	    "books[1] must be at least 10 characters in length",
//	    "books[2] is a required field"
//	  ],
//	  "email": "email is a required field",
//	  "id": "id is a required field",
//	  "userName": "userName is a required field"
//	}
func FormatErrors(err vd.ValidationErrors) map[string]any {
	trans, _ := uni.GetTranslator("en")

	rv := map[string]any{}

	for _, e := range err {
		setValueInMap(rv, e.Namespace(), e.Translate(trans))
	}

	return rv
}
