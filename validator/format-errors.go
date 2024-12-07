package validator

import (
	"strconv"
	"strings"
	"unicode"

	vd "github.com/go-playground/validator/v10"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func setValueInMap(m map[string]interface{}, key string, value string) {
	keys := strings.Split(key, ".")
	currMap := m

	for i := 1; i < len(keys); i++ { // Start from 1 to skip "Data"
		key := toCamelCase(keys[i])
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

func toCamelCase(s string) string {
	words := strings.Split(s, "_")
	for i := 1; i < len(words); i++ {
		words[i] = cases.Title(language.Und).String(words[i])
	}
	return lowerFirst(strings.Join(words, ""))
}

func lowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// FormatErrors formats ValidationErrors into a map of error strings similar to
// the source data type.
//
// Example validation error -
//
//	{
//		"Data.Addresses[0].Street":"Street is a required field",
//		"Data.Books[0]":"Books[0] must be at least 10 characters in length",
//		"Data.Books[1]":"Books[1] must be at least 10 characters in length",
//		"Data.Books[2]":"Books[2] is a required field",
//		"Data.Email":"Email is a required field",
//		"Data.Id":"Id is a required field",
//		"Data.UserName":"UserName is a required field"
//	}
//
// Corresponding error message after formatting -
//
//	{
//	  "addresses": [
//	    {
//	      "street": "Street is a required field"
//	    }
//	  ],
//	  "books": [
//	    "Books[0] must be at least 10 characters in length",
//	    "Books[1] must be at least 10 characters in length",
//	    "Books[2] is a required field"
//	  ],
//	  "email": "Email is a required field",
//	  "id": "Id is a required field",
//	  "userName": "UserName is a required field"
//	}
func FormatErrors(err vd.ValidationErrors) map[string]any {
	trans, _ := uni.GetTranslator("en")
	tErrs := err.Translate(trans)

	rv := map[string]any{}

	for k, v := range tErrs {
		setValueInMap(rv, k, v)
	}

	return rv
}
