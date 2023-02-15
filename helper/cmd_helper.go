package helper

import (
	"os"
	"strings"
)

func CMDUserFirstArgIs(expetedArg string) bool {
	if len(os.Args) <= 1 { // return false if only one argument or less (one argument is default provided by Go-Lang, so it's not considered as first argument from user input)
		return false
	}

	actualArg := os.Args[1]
	actualArg = strings.TrimSpace(actualArg)

	expetedArg = strings.TrimSpace(expetedArg)

	return actualArg == expetedArg
}
