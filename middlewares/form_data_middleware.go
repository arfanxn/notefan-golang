package middlewares

import (
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iancoleman/strcase"
	"github.com/notefan-golang/helpers/errorh"
)

// FormDataMiddleware parses request url values to form data if exists and parses FormData if exists or if not exists it will convert request raw json to form data
func FormDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get max memory of multipart from ENV variable
		maxMemory, err := strconv.ParseInt(os.Getenv("MULTIPART_MAX_MEMORY"), 10, 64)
		errorh.Panic(err)

		// parse multipart form with specified max memory
		err = r.ParseMultipartForm(maxMemory)
		errorh.Panic(err)

		// get request form data and replace existing form key-value with CamelCase key-value
		for key, values := range r.Form {
			nonAlphaRegex := regexp.MustCompile(`[^a-zA-Z ]+`)
			keyCamelCase := strcase.ToCamel(key)
			keyCamelCase = nonAlphaRegex.ReplaceAllString(keyCamelCase, "")
			if key != keyCamelCase {
				continue
			}
			r.Form.Set(keyCamelCase, values[0])
			r.PostForm.Set(keyCamelCase, values[0])
			for index, value := range values {
				keyCamelCaseWithIndex := keyCamelCase + "." + strconv.Itoa(index)
				r.Form.Set(keyCamelCaseWithIndex, value)
				r.PostForm.Set(keyCamelCaseWithIndex, value)
			}
			r.Form.Del(key)
			r.PostForm.Del(key)
		}

		// get request wildcards
		wildcards := mux.Vars(r)
		// replace existing form key-value if it exists and key is match in wildcards
		for key, value := range wildcards {
			key := strcase.ToCamel(key)
			r.Form.Set(key, value)
			r.PostForm.Set(key, value)
		}

		// get request url queries/parameters
		queries := r.URL.Query()
		// replace existing form key-value if it exists and key is match in queries
		for key, values := range queries {
			key := strcase.ToCamel(key)
			r.Form.Set(key, values[0])
			r.PostForm.Set(key, values[0])
			for index, value := range values {
				key := key + "." + strconv.Itoa(index)
				r.Form.Set(key, value)
				r.PostForm.Set(key, value)
			}
		}

		next.ServeHTTP(w, r)
	})
}
