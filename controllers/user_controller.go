package controllers

import (
	"net/http"

	"github.com/notefan-golang/helper"
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
	userId := r.Context().Value("user").(map[string]any)["id"].(string)
	user, err := controller.service.Repository.Find(r.Context(), userId)
	helper.ErrorPanic(err)

	response := responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully retrived current user").
		Body("user", responses.User{
			Id:        user.Id.String(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt.Time,
		})
	helper.ResponseJSON(w, response)
}
