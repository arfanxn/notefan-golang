package file_reqs

import (
	"bytes"
	"encoding/binary"
	"mime/multipart"

	"github.com/gabriel-vasile/mimetype"
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
func NewFromFH(fileHeader *multipart.FileHeader) (*File, error) {
	openFile, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer openFile.Close()

	fileBuff := new(bytes.Buffer)
	_, err = fileBuff.ReadFrom(openFile)
	if err != nil {
		return nil, err
	}

	mime := mimetype.Detect(fileBuff.Bytes())

	file := new(File)
	file.Name = fileHeader.Filename
	file.Size = fileHeader.Size
	file.Mime = *mime
	file.Header = fileHeader
	file.Buffer = fileBuff

	return file, nil
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
