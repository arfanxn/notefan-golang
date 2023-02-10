package middlewares

import (
	"context"
	"net/http"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/responses"

	"github.com/golang-jwt/jwt/v4"
)

func AuthenticateMiddleware(next http.Handler) http.Handler {
	responseUnauthorized := func(w http.ResponseWriter) (int, error) {
		message := "Unauthorized action, please sign in and try again"
		return helper.ResponseJSON(w, responses.NewResponse().
			Code(http.StatusUnauthorized).
			Error(message))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieAccessToken, err := r.Cookie("Access-Token")
		if err != nil {
			switch err {
			case http.ErrNoCookie:
				responseUnauthorized(w)
			default:
				helper.ResponseJSON(w, responses.NewResponse().
					Code(http.StatusInternalServerError).
					Error(exceptions.SomethingWentWrongError.Error()))
			}
			return
		}

		tokenizer, err := helper.JWTParse(cookieAccessToken.Value) // parse jwt token from cookie access token
		claims, ok := tokenizer.Claims.(jwt.MapClaims)             // get jwt claims

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				responseUnauthorized(w)
				return
			case jwt.ValidationErrorExpired:
				helper.ResponseJSON(w, responses.NewResponse().
					Code(http.StatusUnauthorized).
					Error("Token expired, please sign in again"))
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

		ctx := context.WithValue(context.Background(), "user", claims["user"])
		r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
