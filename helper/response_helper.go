package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notion-golang/models/responses"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func ResponseJSON[T any](w http.ResponseWriter, response responses.Response[T]) (int, error) {
	bytes, err := json.Marshal(response)

	LogFatalIfError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	return w.Write(bytes)
}

func ResponseJSONFromError(w http.ResponseWriter, code int, err error) (int, error) {
	response := responses.Response[any]{
		Code:    code,
		Message: err.Error(),
	}
	return ResponseJSON(w, response)
}

func ResponseJSONFromValidatorErrorsWithTranslator(
	w http.ResponseWriter,
	errs validator.ValidationErrors,
	translator ut.Translator,
) (int, error) {
	data := make(map[string]string)
	for _, err := range errs {
		field := err.Field()
		msg := err.Translate(translator)
		msg = strings.ReplaceAll(msg, field, StrSnakeCaseToCapitalized(field))
		data[field] = msg
	}

	fmt.Println(data)

	response := responses.NewResponse(
		http.StatusUnprocessableEntity,
		responses.MessageError,
		data,
	)

	return ResponseJSON(w, response)
}
