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
	input, err := helper.RequestParseBodyThenValidateAndWriteResponseIfError[requests.AuthLogin](w, r)
	if err != nil {
		return
	}

	user, token, err := controller.service.Login(r.Context(), input)
	if err != nil {
		response := responses.NewResponse().
			Code(http.StatusUnauthorized).
			Error(exceptions.AuthFailedLogin.Error())
		helper.ResponseJSON(w, response)
		return
	}

	// Set token to the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Access-Token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	// Send the signed in user and token on the response
	response := responses.NewResponse().
		Code(http.StatusOK).
		Success("Login successfully").
		Body("user", responses.AuthLogin{
			Id:          user.Id.String(),
			Name:        user.Name,
			Email:       user.Email,
			AccessToken: token,
		})
	helper.ResponseJSON(w, response)
}

func (controller AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Delete token from cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Access-Token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})
	helper.ResponseJSON(w, responses.NewResponse().Code(http.StatusOK).Success("Logout successfully"))
}

func (controller AuthController) Register(w http.ResponseWriter, r *http.Request) {
	input, err := helper.RequestParseBodyThenValidateAndWriteResponseIfError[requests.AuthRegister](w, r)
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
