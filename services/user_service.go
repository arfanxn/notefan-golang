package services

import (
	"context"

	"github.com/notefan-golang/repositories"
)

type UserService struct {
	Repository *repositories.UserRepository
}

func NewUserService(repository *repositories.UserRepository) *UserService {
	return &UserService{Repository: repository}
}

func (service *UserService) Something(ctx context.Context) {
	// TODO
}
