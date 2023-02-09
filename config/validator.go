package config

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func InitiliazeValidator() *validator.Validate {
	validate := validator.New()

	// This will makes field name same as json tag on the field
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}

func InitiliazeValidatorUniversalTranslator(validate *validator.Validate) *ut.UniversalTranslator {
	enTrans := en.New() // English Translator
	idTrans := id.New() // Indonesian Translator
	universalTranslator := ut.New(enTrans, idTrans)
	return universalTranslator
}

func InitiliazeValidatorAndTranslator(lang string) (*validator.Validate, ut.Translator) {
	validate := InitiliazeValidator()
	universalTranslator := InitiliazeValidatorUniversalTranslator(validate)
	translator, _ := universalTranslator.GetTranslator(lang)
	return validate, translator
}
