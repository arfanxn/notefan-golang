package services

import (
	"context"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests"
	"github.com/notefan-golang/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository *repositories.UserRepository
}

func NewAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{userRepository: userRepository}
}

func (service *AuthService) Login(ctx context.Context, data requests.AuthLogin) (
	entities.User, string, error) {
	user, err := service.userRepository.FindByEmail(ctx, data.Email)
	if err != nil { // err not nil == user not found, return exception HTTPAuthLoginFailed
		helper.ErrorLog(err)
		return user, "", exceptions.HTTPAuthLoginFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil { // err not nil == password doesnt match, return exception HTTPAuthLoginFailed
		helper.ErrorLog(err)
		return user, "", exceptions.HTTPAuthLoginFailed
	}

	// Generate JWT token for user authentication
	token, err := helper.JWTGenerate(user)
	helper.ErrorPanic(err) // panic if token generation failed

	return user, token, nil
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
