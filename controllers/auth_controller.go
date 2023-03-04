package controllers

import (
	"net/http"

	"github.com/notefan-golang/helpers/combh"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/rwh"
	authReqs "github.com/notefan-golang/models/requests/auth_reqs"
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
	input, err := combh.RequestBodyDecodeValidate[authReqs.Login](r.Body)
	errorh.Panic(err)

	authLoginRes, err := controller.service.Login(r.Context(), input)
	errorh.Panic(err)

	// Set token to the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Path:     "/",
		Value:    authLoginRes.AccessToken,
		HttpOnly: true,
	})

	// Send the signed in user and token on the response
	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Login successfully").
			Body("user", authLoginRes),
	)
}

func (controller AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Delete token from cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})
	rwh.WriteResponse(w, responses.NewResponse().Code(http.StatusOK).Success("Logout successfully"))
}

func (controller AuthController) Register(w http.ResponseWriter, r *http.Request) {
	input, err := combh.RequestBodyDecodeValidate[authReqs.Register](r.Body)
	errorh.Panic(err)

	// Register the user
	userRes, err := controller.service.Register(r.Context(), input)
	errorh.LogPanic(err)

	// Send response and the registered user
	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusCreated).
		Success("Successfully registered").
		Body("user", userRes))
}
