package middlewares

import (
	"net/http"
	"notefan-golang/helper"
	"notefan-golang/models/responses"
)

// Error recover/catcher middleware
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				helper.ErrorLog(err) // log the error to log files

				helper.ResponseJSON(w, responses.NewResponse().
					Code(http.StatusInternalServerError).Error("Something went wrong"),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
