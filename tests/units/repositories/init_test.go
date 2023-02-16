package repositories

import (
	"os"
	"path"
	"runtime"

	"github.com/notefan-golang/config"
)

func init() {
	changeDirToRoot()
	config.LoadTestENV()
}

// changeDirToRoot will change the working directory of test to relative path
func changeDirToRoot() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..", "..", "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
