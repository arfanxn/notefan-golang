package file_reqs

import (
	"bytes"
	"mime/multipart"

	"github.com/gabriel-vasile/mimetype"
	"github.com/notefan-golang/helpers/errorh"
)

type File struct {
	Name   string                `json:"name"` // Name stores path and name of the file, example: "dir/file.txt"
	Size   int64                 `json:"size"`
	Mime   mimetype.MIME         `json:"-"`
	Header *multipart.FileHeader `json:"-"`
	Buffer *bytes.Buffer         `json:"-"`
}

// FillFromFileHeader fills file from file header
func FillFromFileHeader(fileHeader *multipart.FileHeader) File {
	openFile, err := fileHeader.Open()
	errorh.LogPanic(err)
	defer openFile.Close()

	fileBuff := new(bytes.Buffer)
	fileBuff.ReadFrom(openFile)

	mime := mimetype.Detect(fileBuff.Bytes())

	file := File{}
	file.Name = fileHeader.Filename
	file.Size = fileHeader.Size
	file.Mime = *mime
	file.Header = fileHeader
	file.Buffer = fileBuff

	return file
}
