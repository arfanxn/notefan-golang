package middlewares

import (
	"net/http"
)

// AfterMiddleware do something after the request has finished
func AfterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		// Do something after the request has finished

		// CLose and Delete Form data fields
		r.Body.Close()
		for key := range r.Form {
			r.Form.Del(key)
		}
		for key := range r.PostForm {
			r.Form.Del(key)
		}

	})
}
