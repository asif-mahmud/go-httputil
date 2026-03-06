package validator

import (
	"reflect"
	"strings"
)

// ExtractTagName figures out the field name from various tags of struct field.
//
// The selection precedence is as follows -
//
// 1. form tag
// 2. path tag
// 3. json tag
// 4. struct field name
func ExtractTagName(fld reflect.StructField) string {
	for _, tag := range []string{"form", "path", "json"} {
		tagName := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]
		if len(tagName) > 0 {
			return tagName
		}
		if tagName == "-" {
			return ""
		}
	}

	return fld.Name
}
