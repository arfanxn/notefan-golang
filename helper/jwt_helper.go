package helper

import (
	"os"
	"time"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/models/entities"

	"github.com/golang-jwt/jwt/v4"
)

func JWTGenerate(user entities.User) (string, error) {
	tokenizer := jwt.New(jwt.SigningMethodHS256)
	claims := tokenizer.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["id"] = user.Id.String()
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix() // expires in N minutes
	tokenizer.Claims = claims

	signingKey := []byte(os.Getenv("APP_KEY"))
	token, err := tokenizer.SignedString(signingKey)

	ErrorLog(err)

	return token, err
}

func JWTParse(accessToken string) (*jwt.Token, error) {
	tokenizer, err := jwt.Parse(accessToken, func(tokenizer *jwt.Token) (any, error) {
		_, ok := tokenizer.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			err := exceptions.JWTInvalidSigningMethod
			ErrorLog(err)
			return nil, err
		}

		signingKey := []byte(os.Getenv("APP_KEY"))
		return signingKey, nil
	})
	ErrorLog(err)
	return tokenizer, err
}
