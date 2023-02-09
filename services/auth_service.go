package services

import (
	"context"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/models/requests"
	"notefan-golang/models/responses"
	"notefan-golang/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepo
}

func NewAuthService(userRepo *repositories.UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (service *AuthService) Login(ctx context.Context, data requests.AuthLogin) (responses.AuthLogin, error) {
	user, err := service.userRepo.FindByEmail(ctx, data.Email)
	if err != nil {
		return responses.AuthLogin{}, exceptions.AuthFailedLogin
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return responses.AuthLogin{}, exceptions.AuthFailedLogin
	}

	// Generate JWT token for user authentication
	token, err := helper.JWTGenerate(user)
	if err != nil {
		return responses.AuthLogin{}, exceptions.AuthFailedLogin
	}

	return responses.AuthLogin{
		Id:    user.Id.String(),
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}, nil
}

// Register
// Register the given user
func (service *AuthService) Register(ctx context.Context, data requests.AuthRegister) (entities.User, error) {
	// Hash the user password
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		helper.LogIfError(err)
		return entities.User{}, err
	}

	// Prepare the User struct
	user := entities.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: string(password),
	}

	// Save the user into Database
	user, err = service.userRepo.Create(ctx, user)
	if err != nil {
		return entities.User{}, exceptions.AuthFailedRegister
	}

	// Return the created user and nil
	return user, nil
}
