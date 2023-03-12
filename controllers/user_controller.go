package controllers

import (
	"net/http"

	media_collnames "github.com/notefan-golang/enums/media/collection_names"
	"github.com/notefan-golang/helpers/contexth"
	"github.com/notefan-golang/helpers/decodeh"
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
	errorh.LogPanic(err)

	response := responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully retrive current user").
		Body("user", user)
	rwh.WriteResponse(w, response)
}

// UpdateProfileSelf updates current logged in user profile
func (controller UserController) UpdateProfileSelf(w http.ResponseWriter, r *http.Request) {
	// Get user avatar file haeder
	avatarFH, _ := rwh.RequestFormFileHeader(r, media_collnames.Avatar)

	// Decode request form data
	input, err := decodeh.FormData[user_reqs.UpdateProfile](r.MultipartForm.Value)
	errorh.LogPanic(err)
	input.Id = contexth.GetAuthUserId(r.Context())
	if avatarFH != nil {
		input.Avatar = file_reqs.NewFromFH(avatarFH)
	}

	// Validate input
	err = validationh.ValidateStruct(input)
	errorh.LogPanic(err)

	userRes, err := controller.service.UpdateProfile(r.Context(), input)
	errorh.LogPanic(err)

	response := responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully update current user").
		Body("user", userRes)
	rwh.WriteResponse(w, response)
}
