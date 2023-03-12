package file_rules

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/sliceh"
	fileReqs "github.com/notefan-golang/models/requests/file_reqs"
)

// Size is rule for validating file size, make sure the value that'll be validated is type of multipart.FileHeader
func File(required bool, min int64, max int64, mimeTypes []string) validation.RuleFunc {
	return func(value interface{}) error {
		filePtr, err := value.(*fileReqs.File) // parse value to pointer fileReqs.File
		errorh.Panic(err)

		// if file is not required in rule and the file is not provided then immediately return nil error
		if required == false {
			return nil
		}

		// check if file is required but not exists
		if (required) && ((filePtr == nil) || (filePtr.Size == 0)) {
			return exceptions.ValidationFileNotProvided
		}

		// check if file size is between the specified min and max size
		if filePtr.Size < min || filePtr.Size > max {
			return exceptions.ValidationFileSize
		}

		// check if mimetype rules are available
		if len(mimeTypes) > 0 {
			matchedMimeTypes := sliceh.Filter(mimeTypes, func(mimeType string) bool {
				return strings.ToLower(filePtr.Mime.String()) == strings.ToLower(mimeType)
			})

			// if no matching mime types found return an ValidationFileMimeType error
			if len(matchedMimeTypes) == 0 {
				return exceptions.ValidationFileMimeType
			}
		}

		return nil
	}
}
