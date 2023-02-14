package services

import (
	"context"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/models/requests"
	"notefan-golang/repositories"

	"github.com/google/uuid"
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

	// Prepare the User struct
	user := entities.User{
		Id:       uuid.New(),
		Name:     data.Name,
		Email:    data.Email,
		Password: string(password),
	}

	// Save the user into Database
	user, err = service.userRepository.Create(ctx, user)
	helper.ErrorPanic(err) // panic if save into db failed

	// Return the created user and nil
	return user, nil
}
