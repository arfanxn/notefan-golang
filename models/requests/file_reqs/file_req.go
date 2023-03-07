package file_reqs

import (
	"bytes"
	"encoding/binary"
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

/*
 * ----------------------------------------------------------------
 * Instantiate methods ⬇
 * ----------------------------------------------------------------
 */

// NewFromFH instantiates from FileHeader
func NewFromFH(fileHeader *multipart.FileHeader) *File {
	openFile, err := fileHeader.Open()
	errorh.LogPanic(err)
	defer openFile.Close()

	fileBuff := new(bytes.Buffer)
	fileBuff.ReadFrom(openFile)

	mime := mimetype.Detect(fileBuff.Bytes())

	file := new(File)
	file.Name = fileHeader.Filename
	file.Size = fileHeader.Size
	file.Mime = *mime
	file.Header = fileHeader
	file.Buffer = fileBuff

	return file
}

// NewFromBytes instantiates a new instance from the given bytes
func NewFromBytes(fileBytes []byte) *File {
	mime := mimetype.Detect(fileBytes)

	file := new(File)
	file.Size = int64(binary.Size(fileBytes))
	file.Mime = *mime
	file.Buffer = bytes.NewBuffer(fileBytes)

	return file
}

/*
 * ----------------------------------------------------------------
 *  Struct's methods ⬇
 * ----------------------------------------------------------------
 */

// IsProvided checks whether File is provided
func (file *File) IsProvided() bool {
	switch true {
	case file.Size <= 0:
		return false
	default:
		return true
	}
}

// SetBuffer sets buffer and the buffer related fields
func (file *File) SetBuffer(buffer *bytes.Buffer) {
	fileBytes := buffer.Bytes()
	mime := mimetype.Detect(fileBytes)

	file.Size = int64(binary.Size(fileBytes))
	file.Mime = *mime
	file.Buffer = buffer
}
