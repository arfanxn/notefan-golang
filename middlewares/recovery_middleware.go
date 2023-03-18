package middlewares

import (
	"fmt"
	"net/http"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/rwh"
)

// Error recover/catcher middleware
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var (
				anyErr  = recover() // recover any error
				written int         // written represents the number of bytes written to the client response
				err     error       // error state
			)
			_ = err

			if anyErr != nil {
				errorh.Log(anyErr) // log the error to log files

				// check if the error is http error
				httpErr, ok := anyErr.(*exceptions.HTTPError)
				if ok {
					fmt.Println("recovery middleware, error is http error", httpErr.Code)
					// If error is unprocessable entity (validation failed)
					if httpErr.Code == http.StatusUnprocessableEntity {
						written, err = rwh.WriteValidationErrorResponse(w, httpErr.Err)
					}
				}

				// if no bytes were written to the client response then write something went wrong error response
				if written == 0 {
					rwh.WriteSomethingWentWrongResponse(w)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
