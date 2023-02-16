package middlewares

import (
	"context"
	"net/http"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/responses"

	"github.com/golang-jwt/jwt/v4"
)

func AuthenticateMiddleware(next http.Handler) http.Handler {
	responseUnauthorized := func(w http.ResponseWriter) (int, error) {
		return helper.ResponseJSON(w, responses.NewResponse().
			Code(exceptions.HTTPAuthNotSignIn.Code).
			Error(exceptions.HTTPAuthNotSignIn.Error()))
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
					Error(exceptions.HTTPSomethingWentWrong.Error()))
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

		ctx := context.WithValue(context.Background(), "user", map[string]any{
			"id":    claims["id"],
			"name":  claims["name"],
			"email": claims["email"],
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
