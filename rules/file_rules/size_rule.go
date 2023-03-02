package file_rules

import (
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/notefan-golang/exceptions"
)

// Size is rule for validating file size, make sure the value that'll be validated is type of multipart.FileHeader
func Size(min int64, max int64) validation.RuleFunc {
	return func(value interface{}) error {
		fileHeader, _ := value.(multipart.FileHeader)
		fileSize := fileHeader.Size

		if fileSize < min || fileSize > max {
			return exceptions.ValidationFileSize
		}

		return nil
	}
}
