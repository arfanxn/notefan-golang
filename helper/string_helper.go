package helper

import (
	"strings"
	"unicode"
)

func StrSnakeCaseToCapitalized(str string) string {
	str = strings.ReplaceAll(str, "_", " ")
	str = strings.TrimSpace(str)
	firstCharOfStr := []rune(str)[0]
	str = StrReplaceAtIndex(str, 0, unicode.ToUpper(firstCharOfStr))
	return str
}

func StrReplaceAtIndex(in string, i int, r rune) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}
