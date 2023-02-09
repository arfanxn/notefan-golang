package helper

import (
	"notefan-golang/config"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
)

func InitializeValidatorAndDetermineTranslator(lang string) (*validator.Validate, ut.Translator) {
	validate, translator := config.InitiliazeValidatorAndTranslator(lang)

	// This switch statement will determines which translation should be used
	switch lang {
	case "id":
		idTranslations.RegisterDefaultTranslations(validate, translator)
		break
	default: // default lang is English, currently only supports two translations
		enTranslations.RegisterDefaultTranslations(validate, translator)
		break
	}
	return validate, translator
}

func ValidatorTranslateErrors(errs validator.ValidationErrors, translator ut.Translator) map[string]string {
	translatedErrs := make(map[string]string)

	// Format and translate the error messages
	for _, err := range errs {
		field := err.Field()
		msg := err.Translate(translator)
		msg = strings.ReplaceAll(msg, field, StrSnakeCaseToCapitalized(field))
		translatedErrs[field] = msg
	}

	return translatedErrs
}
