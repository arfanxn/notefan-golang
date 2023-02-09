package controllers

import (
	"net/http"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/requests"
	"notefan-golang/models/responses"
	"notefan-golang/services"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (controller AuthController) Login(w http.ResponseWriter, r *http.Request) {
	input, err := helper.ParseRequestBodyThenValidateAndWriteResponseIfError[requests.AuthLogin](w, r)
	if err != nil {
		return
	}

	user, err := controller.service.Login(r.Context(), input)
	if err != nil {
		response := responses.NewResponse().
			Code(http.StatusUnauthorized).
			Error(exceptions.AuthFailedRegister.Error())
		helper.ResponseJSON(w, response)
		return
	}

	// Send response and the signed in user
	response := responses.NewResponse().
		Code(http.StatusOK).
		Success("Login successfully").
		Body("user", user)
	helper.ResponseJSON(w, response)
}

func (controller AuthController) Logout(w http.ResponseWriter, r *http.Request) {

}

func (controller AuthController) Register(w http.ResponseWriter, r *http.Request) {
	input, err := helper.ParseRequestBodyThenValidateAndWriteResponseIfError[requests.AuthRegister](w, r)
	if err != nil {
		return
	}

	// Register the user
	user, err := controller.service.Register(r.Context(), input)
	if err != nil {
		response := responses.NewResponse().
			Code(http.StatusInternalServerError).
			Error(exceptions.AuthFailedRegister.Error())
		helper.ResponseJSON(w, response)
		return
	}

	// Send response and the registered user
	response := responses.NewResponse().
		Code(http.StatusCreated).
		Success("Successfully registered").
		Body("user", user)
	helper.ResponseJSON(w, response)
}
