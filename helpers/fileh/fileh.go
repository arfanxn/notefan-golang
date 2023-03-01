package fileh

import (
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/notefan-golang/helpers/errorh"
)

// RandFromDir retrieves a random file from the given directory
func RandFromDir(dirpath string) (*os.File, error) {
	var f *os.File
	entries, err := os.ReadDir(dirpath)
	if err != nil {
		errorh.Log(err)
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
		errorh.Log(err)
		return f, err
	}

	return f, nil
}

// FileSize returns the size of the file, if error occurs it will return -1
func GetSize(f fs.File) int64 {
	fileInfo, err := f.Stat()
	if err != nil { // error happens return -1
		errorh.Log(err)
		return -1
	}
	return fileInfo.Size()
}

// RemoveByPath removes recursively from the path
func RemoveByPath(paths ...string) (err error) {
	for _, path := range paths {
		err = os.RemoveAll(path)
	}
	return
}
