package helper

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
)

func NewValidatorAndTranslator(locale string) (*validator.Validate, ut.Translator) {
	validate := validator.New()

	// This will makes field name same as json tag on the field
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	englishTranslator := en.New()
	indonesiaTranslator := id.New()
	universalTranslator := ut.New(englishTranslator, indonesiaTranslator)
	translator, _ := universalTranslator.GetTranslator(locale)

	// currently only supports two translations
	switch locale {
	case "id":
		idTranslations.RegisterDefaultTranslations(validate, translator)
		break
	default: // default locale is English
		enTranslations.RegisterDefaultTranslations(validate, translator)
		break
	}
	return validate, translator
}