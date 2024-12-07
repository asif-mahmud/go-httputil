package validator

import (
	"github.com/go-playground/form/v4"
	en "github.com/go-playground/locales/en"
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/mold/v4/scrubbers"
	ut "github.com/go-playground/universal-translator"
	vd "github.com/go-playground/validator/v10"
	ve "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni         *ut.UniversalTranslator
	validate    *vd.Validate
	formDecoder *form.Decoder
	conform     *mold.Transformer
	scrub       *mold.Transformer
)

func init() {
	// setup universal translations
	e := en.New()
	uni = ut.New(e, e)

	// setup validator's translations
	validate = vd.New()
	et, _ := uni.GetTranslator("en")
	ve.RegisterDefaultTranslations(validate, et)

	// setup form (url.Values) decoder
	formDecoder = form.NewDecoder()

	// setup transformers
	conform = modifiers.New()
	scrub = scrubbers.New()
}
