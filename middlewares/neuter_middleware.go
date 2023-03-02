package middlewares

import (
	"net/http"
	"strings"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
)

// NeuterMiddleware prevents listing directories from file server
func NeuterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/public") && strings.HasSuffix(r.URL.Path, "/") {
			errorh.Panic(exceptions.HTTPNotFound)
		}

		next.ServeHTTP(w, r)
	})
}
