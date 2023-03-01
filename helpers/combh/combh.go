// Combh is Combination helper
package combh

import (
	"io"
	"net/http"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/jsonh"
	reqContracts "github.com/notefan-golang/models/requests/req_contracts"
)

// RequestBodyDecodeValidate will decode the request body and then validate the decoded request body
// then return decoded request body and error if it fails
func RequestBodyDecodeValidate[T reqContracts.ValidateableContract](requestBody io.ReadCloser) (
	T, error) {
	// Decode the request body
	input, err := jsonh.DecodeFromReader[T](requestBody)
	errorh.LogPanic(err)

	// Validate the input
	validationErrs := input.Validate()
	if validationErrs != nil {
		return input, exceptions.NewHTTPError(http.StatusUnprocessableEntity, validationErrs)
	}

	return input, nil
}