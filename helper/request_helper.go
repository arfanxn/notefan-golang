package helper

import (
	"net/http"
	"notefan-golang/exceptions"
	"notefan-golang/models/responses"
	"os"

	"github.com/go-playground/validator/v10"
)

func RequestGetLanguage(r http.Request) string {
	acceptLang := r.Header.Get("Accept-Language")
	lang := r.URL.Query().Get("lang")

	switch true {
	case acceptLang != "":
		return acceptLang
	case lang != "":
		return lang
	default:
		return os.Getenv("APP_LANG") // get default language from environment
	}
}

func RequestParseBodyThenValidateAndWriteResponseIfError[T any](
	w http.ResponseWriter,
	r *http.Request,
) (T, error) {
	// Parse request body to golang struct
	parsedReqBody, err := JSONDecodeFromReader[T](r.Body)
	ErrorPanic(err)
	defer r.Body.Close()

	lang := RequestGetLanguage(*r)
	validate, trans := ValidatorInitAndDetermineTranslator(lang)

	// Validate parsed request body, if validation error write error to response
	if err := validate.Struct(parsedReqBody); err != nil {
		response := responses.NewResponse().
			Code(http.StatusUnprocessableEntity).
			Error(exceptions.HTTPValidationFailed.Error()).
			Body("errors", ValidatorTranslateErrors(err.(validator.ValidationErrors), trans))
		ResponseJSON(w, response)
		return parsedReqBody, err
	}

	return parsedReqBody, nil
}
