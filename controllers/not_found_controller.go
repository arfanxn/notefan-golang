package controllers

import (
	"net/http"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/rwh"
	"github.com/notefan-golang/models/responses"
)

type NotFoundController struct {
}

func NewNotFoundController() *NotFoundController {
	return &NotFoundController{}
}

// Handle handles the not found http request
func (controller NotFoundController) Handle(w http.ResponseWriter, r *http.Request) {
	response := responses.NewResponse().
		Code(http.StatusOK).
		Error(exceptions.HTTPNotFound.Error())
	rwh.WriteResponse(w, response)
}
