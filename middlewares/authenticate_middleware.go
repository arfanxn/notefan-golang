package middlewares

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/jwth"
	"github.com/notefan-golang/helpers/rwh"
	"github.com/notefan-golang/models/responses"

	"github.com/golang-jwt/jwt/v4"
)

func AuthenticateMiddleware(next http.Handler) http.Handler {
	responseUnauthorized := func(w http.ResponseWriter) (int, error) {
		return rwh.WriteResponse(w, responses.NewResponse().
			Code(exceptions.HTTPAuthNotSignIn.Code).
			Error(exceptions.HTTPAuthNotSignIn.Error()))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieAccessToken, err := r.Cookie("Authorization")
		if err != nil {
			switch err {
			case http.ErrNoCookie:
				responseUnauthorized(w)
			default:
				rwh.WriteResponse(w, responses.NewResponse().
					Code(http.StatusInternalServerError).
					Error(exceptions.HTTPSomethingWentWrong.Error()))
			}
			return
		}

		signature := os.Getenv("APP_KEY")
		tokenizer, err := jwth.Decode(signature, cookieAccessToken.Value) // parse jwt token from cookie access token
		claims, ok := tokenizer.Claims.(jwt.MapClaims)                    // get jwt claims

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				responseUnauthorized(w)
				return
			case jwt.ValidationErrorExpired:
				// response auth token expired
				rwh.WriteResponse(w, responses.NewResponse().
					Code(exceptions.HTTPAuthTokenExpired.Code).
					Error(exceptions.HTTPAuthTokenExpired.Error()))
				return
			default:
				responseUnauthorized(w)
				return
			}
		}

		if !ok || !tokenizer.Valid {
			responseUnauthorized(w)
			return
		}

		// Refresh Authorization coookie max age ⬇
		authorizationCookieName := "Authorization"                         // the cookie name
		authorizationCookie, err := r.Cookie(authorizationCookieName)      // get authorization cookie
		maxAge, err := strconv.ParseInt(os.Getenv("AUTH_MAX_AGE"), 10, 64) // get authorization cookie max age
		errorh.LogPanic(err)                                               // log and panic if error
		http.SetCookie(w, &http.Cookie{                                    // set authorization cookie with new max age and expiration
			Name:     authorizationCookieName,
			Path:     "/",
			Value:    authorizationCookie.Value,
			HttpOnly: true,
			MaxAge:   int(maxAge),
			Expires:  time.Now().Add(time.Duration(maxAge)),
		})

		// extract jwt claims to context
		ctx := context.WithValue(context.Background(), "user", map[string]any{
			"id":    claims["id"],
			"name":  claims["name"],
			"email": claims["email"],
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
