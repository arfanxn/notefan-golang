package helper

import (
	"runtime"
	"strings"
)

func FuncNameFromPC(pc uintptr) string {
	return strings.SplitAfter(runtime.FuncForPC(pc).Name(), ".")[1]
}
