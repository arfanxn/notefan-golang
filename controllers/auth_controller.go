package controllers

import (
	"net/http"
	"notion-golang/exceptions"
	"notion-golang/helper"
	"notion-golang/models/requests"
	"notion-golang/models/responses"
	"notion-golang/services"

	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (controller AuthController) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login Page"))
}

func (controller AuthController) Logout(w http.ResponseWriter, r *http.Request) {

}

func (controller AuthController) Register(w http.ResponseWriter, r *http.Request) {
	// Parse request body to golang struct
	authRegisterRequest, err := helper.JSONDecodeFromReader[requests.AuthRegisterReq](r.Body)
	if err != nil {
		helper.ResponseJSONFromError(w, http.StatusInternalServerError, exceptions.DecodingError)
		return
	}
	defer r.Body.Close()

	// Validate parsed request body
	validate, trans := helper.InitializeValidatorAndDetermineTranslator(helper.RequestGetLanguage(*r))
	if err := validate.Struct(authRegisterRequest); err != nil {
		helper.ResponseJSONFromValidatorErrorsWithTranslation(w, err.(validator.ValidationErrors), trans)
		return
	}

	// Register the user
	user, err := controller.service.Register(r.Context(), authRegisterRequest)
	if err != nil {
		helper.ResponseJSONFromError(w, http.StatusInternalServerError, exceptions.AuthFailedToRegister)
		return
	}

	// Send response and the registered user
	response := responses.NewResponse(http.StatusCreated, responses.MessageSuccess, user)
	helper.ResponseJSON(w, response)
}
