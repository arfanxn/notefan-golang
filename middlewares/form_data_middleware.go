package middlewares

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/notefan-golang/helpers/errorh"
)

// FormDataMiddleware parses request url values to form data if exists and parses FormData if exists or if not exists it will convert request raw json to form data
func FormDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get max memory of multipart from ENV variable
		maxMemory, err := strconv.ParseInt(os.Getenv("MULTIPART_MAX_MEMORY"), 10, 64)
		errorh.LogPanic(err)

		// parse multipart form with specified max memory
		r.ParseMultipartForm(maxMemory)

		// get request wildcards
		wildcards := mux.Vars(r)

		// replace existing form key-value if it exists and key is match in wildcards
		for key, value := range wildcards {
			r.PostForm.Set(key, value)
		}

		next.ServeHTTP(w, r)
	})
}
