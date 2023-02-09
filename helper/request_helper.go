package helper

import (
	"net/http"
	"notion-golang/config"
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
		value, err := config.GetENVKey("APP_LANG") // get default language from environment
		LogIfError(err)
		return value
	}
}
