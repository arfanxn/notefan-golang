package helper

import (
	"net/http"
	"os"
)

func RequestGetLanguage(r http.Request) string {
	acceptLang := r.Header.Get("Accept-Language")
	lang := r.URL.Query().Get("lang")

	switch true {
	case acceptLang != "":
		return acceptLang
	case lang != "":
		return lang
	default:
		return os.Getenv("APP_LANG") // get default language from environment

	}
}
