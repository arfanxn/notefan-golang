package services

import (
	"context"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/models/requests"
	"notefan-golang/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepo
}

func NewAuthService(userRepo *repositories.UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (service *AuthService) Login() {

}

// Register
// Register the given user
func (service *AuthService) Register(ctx context.Context, data requests.AuthRegisterReq) (entities.User, error) {
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
		helper.LogIfError(err)
		return entities.User{}, err
	}

	// Return the created user and error if it was not created
	return user, err
}
