package services

import (
	"context"
	"os"
	"time"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/handlers"
	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests"
	"github.com/notefan-golang/models/responses"
	"github.com/notefan-golang/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository *repositories.UserRepository
}

func NewAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepository: userRepository,
	}
}

func (service *AuthService) Login(ctx context.Context, data requests.AuthLogin) (responses.AuthLogin, error) {
	user, err := service.userRepository.FindByEmail(ctx, data.Email)
	if err != nil { // err not nil == user not found, return exception HTTPAuthLoginFailed
		helper.ErrorLog(err)
		return responses.AuthLogin{}, exceptions.HTTPAuthLoginFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil { // err not nil == password doesnt match, return exception HTTPAuthLoginFailed
		helper.ErrorLog(err)
		return responses.AuthLogin{}, exceptions.HTTPAuthLoginFailed
	}

	// Prepare claims (payload) for the JWT
	claims := map[string]any{}
	claims["authorized"] = true
	claims["id"] = user.Id.String()
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix() // expires in N minutes

	// get app key (signing key)
	signature := os.Getenv("APP_KEY")

	// Encode/Generate JWT token
	token, err := handlers.NewJWTHandler().Encode(signature, claims)
	helper.ErrorPanic(err) // panic if token generation failed

	return responses.AuthLogin{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: token,
	}, nil
}

// Register registers the given user
func (service *AuthService) Register(ctx context.Context, data requests.AuthRegister) (entities.User, error) {
	// Hash the user password
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	helper.ErrorPanic(err) // panic if password hashing failed

	// Save the user into Database
	user, err := service.userRepository.Create(ctx, entities.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: string(password),
	})
	helper.ErrorPanic(err) // panic if save into db failed

	// Return the created user and nil
	return user, nil
}
