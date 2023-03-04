package fileh

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/sliceh"
)

// FileNamesFromDir returns a list of filenames (with path) from specified directory path
// return example: ["dir/fileone.txt", "dir/filetwo.txt"] and error if error occurs
func FileNamesFromDir(dirpath string) ([]string, error) {
	entries, err := os.ReadDir(dirpath)
	if err != nil {
		errorh.Log(err)
		return []string{}, err
	}

	fileNames := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, filepath.Join(dirpath, entry.Name()))
		}
	}

	return fileNames, err
}

// RandFromDir retrieves a random file from the given directory
func RandFromDir(dirpath string) (*os.File, error) {
	fileNames, err := FileNamesFromDir(dirpath)
	if err != nil {
		errorh.Log(err)
		return nil, err
	}

	if len(fileNames) == 0 {
		return nil, os.ErrNotExist
	}

	fileName := sliceh.Random(fileNames)

	f, err := os.Open(fileName)
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
