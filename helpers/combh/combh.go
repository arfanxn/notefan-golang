// Combh is Combination helper
package combh

import (
	"net/url"

	"github.com/gorilla/schema"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/validationh"
	reqContracts "github.com/notefan-golang/models/requests/req_contracts"
)

// FormDataDecodeValidate will decode the form data and then validate the decoded form data
// then return decoded form data and error if it fails
func FormDataDecodeValidate[T reqContracts.ValidateableContract](
	formData url.Values,
) (
	T, error) {
	var input T

	// Decode the request body
	err := schema.NewDecoder().Decode(&input, formData)
	errorh.LogPanic(err)

	// Validate the input
	err = validationh.ValidateStruct(input)
	if err != nil { // if validation fails then return the valiudation error
		return input, err
	}

	return input, nil
}
