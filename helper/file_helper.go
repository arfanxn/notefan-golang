package helper

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/notefan-golang/config"
	"github.com/notefan-golang/models/entities"
)

func FileRandFromDir(dirpath string) (*os.File, error) {
	var f *os.File
	entries, err := os.ReadDir(dirpath)
	if err != nil {
		ErrorLog(err)
		return f, err
	}

	fileNames := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}
	fileName := fileNames[rand.Intn(len(fileNames)-1)]

	path := filepath.Join(dirpath, fileName)
	f, err = os.Open(path)
	if err != nil {
		ErrorLog(err)
		return f, err
	}

	return f, nil
}

func FileContentType(f *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := f.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	// Reset the offset for the next read/write operation
	f.Seek(0, io.SeekStart)

	return contentType, nil
}

// FileSize returns the size of the file, if error occurs it will return -1
func FileSize(f *os.File) int64 {
	fileInfo, err := f.Stat()
	if err != nil { // error happens return -1
		ErrorLog(err)
		return -1
	}
	return fileInfo.Size()
}

// FileURLFromMedia returns a file URL from the given media
func FileURLFromMedia(media entities.Media) string {
	disk := config.FSDisks[media.Disk]
	return disk.Root + "/medias/" + media.FileName
}
