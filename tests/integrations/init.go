package integrations

import (
	"os"
	"path"
	"runtime"
)

func init() {
	changeDirToRoot()
}

// changeDirToRoot will change the working directory to the relative path
func changeDirToRoot() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..", "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
