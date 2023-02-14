package middlewares

import (
	"net/http"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/responses"
)

// Error recover/catcher middleware
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			anyErr := recover()
			if anyErr != nil {
				helper.ErrorLog(anyErr) // log the error to log files

				httpErr, ok := anyErr.(*exceptions.HTTPError)
				if ok {
					helper.ResponseJSON(w, responses.NewResponse().
						Code(httpErr.Code).Error(httpErr.Error()),
					)
					return
				}

				helper.ResponseJSON(w, responses.NewResponse().
					Code(exceptions.HTTPSomethingWentWrong.Code).
					Error(exceptions.HTTPSomethingWentWrong.Error()),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
