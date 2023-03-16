package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/notefan-golang/helpers/combh"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/rwh"
	"github.com/notefan-golang/models/requests/auth_reqs"
	"github.com/notefan-golang/models/responses"
	"github.com/notefan-golang/services"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// Login will signing a user with the given credentials
func (controller AuthController) Login(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[auth_reqs.Login](r.Form)
	errorh.Panic(err)

	authLoginRes, err := controller.service.Login(r.Context(), input)
	errorh.Panic(err)

	// Get cookie max age from config env variable
	maxAge, err := strconv.ParseInt(os.Getenv("AUTH_MAX_AGE"), 10, 64)
	errorh.Panic(err)
	// Set token to the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Path:     "/",
		Value:    authLoginRes.AccessToken,
		HttpOnly: true,
		MaxAge:   int(maxAge),
	})

	// Send the signed in user and token on the response
	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Login successfully").
			Body("user", authLoginRes),
	)
}

// Logout will signing out the signed in if user
func (controller AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Delete token from cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1, // immediately mark as expired
	})
	rwh.WriteResponse(w, responses.NewResponse().Code(http.StatusOK).Success("Logout successfully"))
}

// Register registers a new user with the given username and password and other parameters
func (controller AuthController) Register(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[auth_reqs.Register](r.Form)
	errorh.Panic(err)

	// Register the user
	userRes, err := controller.service.Register(r.Context(), input)
	errorh.Panic(err)

	// Send response with registered user data
	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusCreated).
		Success("Successfully registered").
		Body("user", userRes))

}

// ForgotPassword sends a reset password token to the given user email from request
func (controller AuthController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[auth_reqs.ForgotPassword](r.Form)
	errorh.Panic(err)

	err = controller.service.ForgotPassword(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusCreated).
		Success("Successfully sent reset password token to "+input.Email))
}

// ResetPassword sends a reset password token to the given user email from request
func (controller AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[auth_reqs.ResetPassword](r.Form)
	errorh.Panic(err)

	err = controller.service.ResetPassword(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully reset password"))
}
