package validator

import (
	"fmt"
	"log/slog"

	ut "github.com/go-playground/universal-translator"
	vd "github.com/go-playground/validator/v10"
)

func registrationFunc(
	tag string,
	translation string,
	override bool,
) vd.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return
	}
}

func translateFunc(ut ut.Translator, fe vd.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		slog.Warn(fmt.Sprintf("error translating FieldError: %#v", fe))
		return fe.(error).Error()
	}

	return t
}

// Translation defines a custom translation for a specific tag
type Translation struct {
	// Tag is the validator's tag
	Tag string

	// Translation is the translation template.
	// i.e "{0} is a required field"
	Translation string

	// Override will override any existing translation for the tag.
	Override bool

	// CustomRegisFunc custom registration function for validate.
	//
	// See the implementations in here -
	// https://github.com/go-playground/validator/blob/v10.23.0/translations/en/en.go
	// It can be left empty/nil if there's no variation of a validation tag or
	// no need to specify different registration functions.
	CustomRegisFunc vd.RegisterTranslationsFunc

	// CustomTransFunc is custiom translator function.
	//
	// See the implementations in here -
	// https://github.com/go-playground/validator/blob/v10.23.0/translations/en/en.go
	// It can be left empty/nil if there's no variation of a validation tag or
	// no need to specify different translation functions.
	CustomTransFunc vd.TranslationFunc
}

// RegisterTranslation registers a translation or overrides an existing one.
//
// Example of such translation -
//
//	{
//		Tag:         "required",
//		Translation: "{0} is a required field",
//		Override:    false,
//	}
func RegisterTranslation(t Translation) {
	var err error
	trans, _ := uni.GetTranslator("en")
	if t.CustomTransFunc != nil && t.CustomRegisFunc != nil {
		err = validate.RegisterTranslation(t.Tag, trans, t.CustomRegisFunc, t.CustomTransFunc)
	} else if t.CustomTransFunc != nil && t.CustomRegisFunc == nil {
		err = validate.RegisterTranslation(t.Tag, trans, registrationFunc(t.Tag, t.Translation, t.Override), t.CustomTransFunc)
	} else if t.CustomTransFunc == nil && t.CustomRegisFunc != nil {
		err = validate.RegisterTranslation(t.Tag, trans, t.CustomRegisFunc, translateFunc)
	} else {
		err = validate.RegisterTranslation(t.Tag, trans, registrationFunc(t.Tag, t.Translation, t.Override), translateFunc)
	}

	if err != nil {
		slog.Warn(fmt.Sprintf("error registerring new translation: %s", err.Error()))
		return
	}
}
