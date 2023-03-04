package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
	"github.com/notefan-golang/helpers/contexth"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/rwh"
	"github.com/notefan-golang/helpers/validationh"
	"github.com/notefan-golang/models/requests/file_reqs"
	"github.com/notefan-golang/models/requests/user_reqs"
	"github.com/notefan-golang/models/responses"
	"github.com/notefan-golang/services"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

// Self gets the current logged in user
func (controller UserController) Self(w http.ResponseWriter, r *http.Request) {
	userId := contexth.GetAuthUserId(r.Context())
	user, err := controller.service.Find(r.Context(), userId)
	errorh.Panic(err)

	response := responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully retrive current user").
		Body("user", user)
	rwh.WriteResponse(w, response)
}

// UpdateProfileSelf updates current logged in user profile
func (controller UserController) UpdateProfileSelf(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 10MB memory limit
	err := r.ParseMultipartForm(2 << 20)
	errorh.Panic(err)

	// Get user avatar file haeder
	_, avatarFH, _ := r.FormFile("avatar")

	// Operation before validation
	var input user_reqs.UpdateProfile
	err = schema.NewDecoder().Decode(&input, r.MultipartForm.Value)
	errorh.Panic(err)

	input.Id = contexth.GetAuthUserId(r.Context())
	if avatarFH != nil { // if avatarFH is present then fill input avatar
		input.Avatar = file_reqs.FillFromFileHeader(avatarFH)
	}

	// Validate
	err = validationh.ValidateStruct(input)
	errorh.Panic(err)

	userRes, err := controller.service.UpdateProfile(r.Context(), input)
	errorh.Panic(err)

	response := responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully update current user").
		Body("user", userRes)
	rwh.WriteResponse(w, response)
}
