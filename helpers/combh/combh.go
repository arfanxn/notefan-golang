// Combh is Combination helper
package combh

import (
	"net/url"

	"github.com/notefan-golang/helpers/decodeh"
	"github.com/notefan-golang/helpers/validationh"
	reqContracts "github.com/notefan-golang/models/requests/req_contracts"
)

// FormDataDecodeValidate will decode the form data and then validate the decoded form data
// then return decoded form data and error if the validation fails
func FormDataDecodeValidate[T reqContracts.ValidateableContract](
	formData url.Values,
) (input T, err error) {
	input, err = decodeh.FormData[T](formData)
	if err != nil {
		return input, err
	}

	// Validate the input
	err = validationh.ValidateStruct(input)
	if err != nil { // if validation fails then return the valiudation error
		return input, err
	}

	return input, nil
}
