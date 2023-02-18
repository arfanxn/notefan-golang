package helper

import (
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
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

// FileSize returns the size of the file, if error occurs it will return -1
func FileSize(f fs.File) int64 {
	fileInfo, err := f.Stat()
	if err != nil { // error happens return -1
		ErrorLog(err)
		return -1
	}
	return fileInfo.Size()
}

func FileSave(fileSrc fs.File, pathDst string) (err error) {
	fileDst, err := os.Create(pathDst)
	if err != nil {
		return
	}
	_, err = io.Copy(fileDst, fileSrc)
	if err != nil {
		return
	}
	return
}

func FileRemoveByPath(paths ...string) (err error) {
	for _, path := range paths {
		err = os.RemoveAll(path)
	}
	return
}
