package fileh

import (
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"sync"

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

// BatchRemove removes many files/dirs in one
func BatchRemove(paths ...string) (err error) {
	wg := new(sync.WaitGroup)
	for _, path := range paths {
		if err != nil {
			return
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, path string) {
			defer wg.Done()
			err = os.RemoveAll(path)
		}(wg, path)
	}
	wg.Wait()

	return
}

// Save saves a file by given path
func Save(dstPath string, file io.Reader) (err error) {
	errMkdir := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	if errMkdir != nil {
		err = errMkdir
		return
	}
	fileDst, errCreate := os.Create(dstPath)
	if errCreate != nil {
		err = errCreate
		return
	}
	defer fileDst.Close()

	_, errCopy := io.Copy(fileDst, file)
	if errCopy != nil {
		err = errCopy
		return
	}

	return
}
