package controllers

import (
	"net/http"

	"github.com/notefan-golang/helpers/contexth"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/rwh"
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
		Success("Successfully retrived current user").
		Body("user", user)
	rwh.WriteResponse(w, response)
}

/* // TODO: Updates user with its avatar */
// Update updates the user
func (controller UserController) Update(w http.ResponseWriter, r *http.Request) {
}
