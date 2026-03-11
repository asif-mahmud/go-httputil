package validator

import (
	"strconv"
	"strings"

	vd "github.com/go-playground/validator/v10"
)

func setValueInMap(m map[string]any, key string, value string) {
	keys := strings.Split(key, ".")
	currMap := m

	// If the root doesn't have a struct field part (e.g. "[0]"), we still need to process it
	// In go-playground/validator, for structs it usually looks like "StructName.Field".
	// For slices of structs validated with dive, it might look like "[0].Field"

	startIndex := 0
	if len(keys) > 0 && !strings.HasPrefix(keys[0], "[") && len(keys) > 1 {
		// This skips the top-level struct name which validator attaches by default
		// (e.g., "DTO.Values[0]" -> skip "DTO"). If it is just "[0]", we don't skip.
		startIndex = 1
	}

	for i := startIndex; i < len(keys); i++ {
		k := keys[i]
		last := i == len(keys)-1

		if strings.HasPrefix(k, "[") && startIndex == 1 && i == 1 {
			// Edge case: top-level struct was an array pointer slice like `type DTO []*Data` passed directly
			// The namespace looks like `[0].value`. If startIndex was 1, we shouldn't have skipped it if the first key was `[0]`.
			// Since we check `!strings.HasPrefix(keys[0], "[")`, we covered this.
			// But if it's "DTO[0]", we still process it via handleArrayKey.
		}

		if strings.Contains(k, "[") {
			// If k is just "[0]", parts[0] is "", and parts[1] is "0]".
			// Our handleArrayKey creates an entry under `currMap[""]`. This works but it's nested under a "" key.
			currMap = handleArrayKey(currMap, k, last, value)
		} else {
			currMap = handleMapKey(currMap, k, last, value)
		}
	}
}

func handleArrayKey(
	currMap map[string]any,
	key string,
	last bool,
	value string,
) map[string]any {
	parts := strings.Split(key, "[")
	mapKey := parts[0]
	index, _ := strconv.Atoi(strings.TrimRight(parts[1], "]"))

	if currMap[mapKey] == nil {
		currMap[mapKey] = make([]any, index+1)
	} else if len(currMap[mapKey].([]any)) <= index {
		currMap[mapKey] = append(currMap[mapKey].([]any), make([]any, index-len(currMap[mapKey].([]any))+1)...)
	}

	if last {
		currMap[mapKey].([]any)[index] = value
	} else {
		if currMap[mapKey].([]any)[index] == nil {
			currMap[mapKey].([]any)[index] = make(map[string]any)
		}
		currMap = currMap[mapKey].([]any)[index].(map[string]any)
	}

	return currMap
}

func handleMapKey(
	currMap map[string]any,
	key string,
	last bool,
	value string,
) map[string]any {
	if currMap[key] == nil {
		if last {
			currMap[key] = value
		} else {
			currMap[key] = make(map[string]any)
		}
	}

	if !last {
		currMap = currMap[key].(map[string]any)
	}

	return currMap
}

// FormatErrors formats ValidationErrors into a map or slice structure matching
// the original payload shape.
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
// Example validation error output for User struct -
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
//
// Example validation error output for a slice payload: `[]*User` -
//
//	[
//	  {
//	    "id": "id is a required field"
//	  },
//	  null,
//	  {
//	    "email": "email is a required field"
//	  }
//	]
func FormatErrors(err vd.ValidationErrors) any {
	trans, _ := uni.GetTranslator("en")

	rv := map[string]any{}

	for _, e := range err {
		setValueInMap(rv, e.Namespace(), e.Translate(trans))
	}

	// If the entire payload was a root slice, our parsed map will only
	// have a single empty string key `""`. In this case, we extract the slice
	// and return it directly so the error structure perfectly matches the input DTO.
	if len(rv) == 1 {
		if val, ok := rv[""]; ok {
			return val
		}
	}

	return rv
}
