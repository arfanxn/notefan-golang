package rwh

import (
	"mime/multipart"
	"net/http"

	"github.com/clarketm/json"
	"github.com/iancoleman/strcase"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/models/responses"
)

func WriteResponse(w http.ResponseWriter, response responses.Response) (int, error) {
	bytes, err := json.Marshal(response.GetBody())

	if err != nil {
		return 0, err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.GetCode())
	return w.Write(bytes)
}

func WriteValidationErrorResponse(w http.ResponseWriter, validationErrs error) (int, error) {
	if validationErrs == nil {
		return 0, nil
	}

	bytes, err := json.Marshal(validationErrs)
	if err != nil {
		return 0, err
	}

	var mapValidationErrs map[string]string
	err = json.Unmarshal(bytes, &mapValidationErrs)
	if err != nil {
		return 0, nil
	}

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

	return WriteResponse(w, responses.NewResponse().
		Code(http.StatusUnprocessableEntity).
		Error(exceptions.HTTPValidationFailed.Error()).
		Body("errors", mapValidationErrs),
	)
}

// WriteSomethingWentWrongRespons writes http something went wrong error to the client
func WriteSomethingWentWrongResponse(w http.ResponseWriter) (int, error) {
	return WriteResponse(w, responses.NewResponse().
		Code(exceptions.HTTPSomethingWentWrong.Code).
		Error(exceptions.HTTPSomethingWentWrong.Error()),
	)
}

// RequestFormFileHeader returns the file header from request form by the given key or nil if error encountered
func RequestFormFileHeader(r *http.Request, key string) (*multipart.FileHeader, error) {
	_, fileHeader, err := r.FormFile(key)

	if err != nil {
		return nil, err
	}

	return fileHeader, nil
}
