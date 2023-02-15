package services

import (
	"context"

	"github.com/notefan-golang/repositories"
)

type UserService struct {
	repository *repositories.UserRepository
}

func NewUserService(repository *repositories.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (service *UserService) Something(ctx context.Context) {
	// TODO
}
