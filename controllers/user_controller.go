package controllers

import (
	"net/http"
	"notefan-golang/services"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

func (controller UserController) Something(w http.ResponseWriter, r *http.Request) {
	// TODO 
}
