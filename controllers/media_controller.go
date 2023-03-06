package controllers

import (
	"net/http"

	"github.com/notefan-golang/helpers/combh"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/rwh"
	"github.com/notefan-golang/models/requests/common_reqs"
	"github.com/notefan-golang/models/responses"
	"github.com/notefan-golang/services"
)

type MediaController struct {
	service *services.MediaService
}

func NewMediaController(service *services.MediaService) *MediaController {
	return &MediaController{
		service: service,
	}
}

func (controller MediaController) Find(w http.ResponseWriter, r *http.Request) {
	input, err := combh.RequestBodyDecodeValidate[common_reqs.UUID](r.Body)
	errorh.Panic(err)

	// Register the user
	mediaRes, err := controller.service.Find(r.Context(), input)
	errorh.Panic(err)

	errorh.LogPanic(err)

	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully retrieve media").
		Body("media", mediaRes),
	)
}
