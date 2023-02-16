package controllers

import (
	"net/http"

	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/requests"
	"github.com/notefan-golang/models/responses"
	"github.com/notefan-golang/services"
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
	helper.ErrorPanic(err)

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
	helper.ErrorPanic(err)

	// Send response and the registered user
	helper.ResponseJSON(w, responses.NewResponse().
		Code(http.StatusCreated).
		Success("Successfully registered").
		Body("user", responses.User{
			Id:        user.Id.String(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}))
}
