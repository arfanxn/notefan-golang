package helper

import (
	"notefan-golang/models/entities"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func JWTGenerate(user entities.User) (string, error) {
	tokenizer := jwt.New(jwt.SigningMethodHS256)
	claims := tokenizer.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix() // expires in N minutes

	signingKey := []byte(os.Getenv("APP_KEY"))
	token, err := tokenizer.SignedString(signingKey)

	LogIfError(err)

	return token, err
}
