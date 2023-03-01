package middlewares

import (
	"net/http"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/rwh"
	"github.com/notefan-golang/models/responses"
)

// Error recover/catcher middleware
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			anyErr := recover()
			if anyErr != nil {
				errorh.Log(anyErr) // log the error to log files

				// check if the error is http error
				httpErr, ok := anyErr.(*exceptions.HTTPError)
				if ok {
					// If error is unprocessable entity (validation failed)
					if httpErr.Code == http.StatusUnprocessableEntity {
						rwh.WriteValidationErrorResponse(w, httpErr.Err)
						return
					}

					rwh.WriteResponse(w, responses.NewResponse().
						Code(httpErr.Code).Error(httpErr.Error()),
					)
					return
				}

				rwh.WriteResponse(w, responses.NewResponse().
					Code(exceptions.HTTPSomethingWentWrong.Code).
					Error(exceptions.HTTPSomethingWentWrong.Error()),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
