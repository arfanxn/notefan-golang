package services

import (
	"context"
	"notefan-golang/repositories"
)

type UserService struct {
	repo *repositories.UserRepo
}

func NewUserService(repo *repositories.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (service *UserService) Something(ctx context.Context) {
	// TODO
}
