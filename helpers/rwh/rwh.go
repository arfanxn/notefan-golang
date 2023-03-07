// RWH stands for Reader Writer Helper
package rwh

import (
	"net/http"

	"github.com/clarketm/json"
	"github.com/iancoleman/strcase"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/models/responses"
)

func WriteResponse(w http.ResponseWriter, response responses.Response) (int, error) {
	bytes, err := json.Marshal(response.GetBody())

	if err != nil {
		errorh.Log(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.GetCode())
	return w.Write(bytes)
}

func WriteValidationErrorResponse(w http.ResponseWriter, validationErrs error) {
	if validationErrs == nil {
		return
	}

	bytes, err := json.Marshal(validationErrs)
	errorh.LogPanic(err)

	var mapValidationErrs map[string]string
	err = json.Unmarshal(bytes, &mapValidationErrs)
	errorh.LogPanic(err)

	for key, value := range mapValidationErrs {
		// string to snake case
		keySnakeCase := strcase.ToSnake(key)

		if keySnakeCase == key { // if key is already snake_case then continue the loop
			continue
		}

		// set/add map with snake_case key
		mapValidationErrs[keySnakeCase] = value

		// delete/remove map with key other than snake_case
		delete(mapValidationErrs, key)
	}

	WriteResponse(w, responses.NewResponse().
		Code(http.StatusUnprocessableEntity).
		Error(exceptions.HTTPValidationFailed.Error()).
		Body("errors", mapValidationErrs),
	)
}
