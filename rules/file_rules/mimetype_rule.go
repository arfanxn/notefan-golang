package file_rules

import (
	"mime/multipart"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
)

func MimeType(mimeTypes ...string) validation.RuleFunc {
	return func(value interface{}) error {
		fileHeader, _ := value.(multipart.FileHeader)
		file, err := fileHeader.Open()
		errorh.Panic(err)
		defer file.Close()
		mime, err := mimetype.DetectReader(file)
		errorh.Panic(err)

		for _, mimeType := range mimeTypes {
			if strings.ToLower(mime.String()) == strings.ToLower(mimeType) {
				return nil
			}
		}

		return exceptions.ValidationFileMimeType
	}
}
