package controllers

import (
	"net/http"

	"github.com/notefan-golang/helpers/combh"
	"github.com/notefan-golang/helpers/decodeh"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/rwh"
	"github.com/notefan-golang/helpers/validationh"
	"github.com/notefan-golang/models/requests/common_reqs"
	"github.com/notefan-golang/models/requests/file_reqs"
	"github.com/notefan-golang/models/requests/media_reqs"
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

// Find finds media by request form data id
func (controller MediaController) Find(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[common_reqs.UUID](r.Form)
	errorh.Panic(err)

	// Find media by id
	mediaRes, err := controller.service.Find(r.Context(), input)
	errorh.Panic(err)

	errorh.LogPanic(err)

	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully retrieve media").
		Body("media", mediaRes),
	)
}

// Update updates media by request form data id
func (controller MediaController) Update(w http.ResponseWriter, r *http.Request) {
	input := decodeh.FormData[media_reqs.Update](r.Form)

	// Get file header from form data
	fileHeader, _ := rwh.RequestFormFileHeader(r, "file")
	if fileHeader != nil {
		input.File = file_reqs.NewFromFH(fileHeader)
	}

	// Validate input
	err := validationh.ValidateStruct(input)
	errorh.Panic(err)

	// Update media by id
	mediaRes, err := controller.service.Update(r.Context(), input)
	errorh.Panic(err)

	errorh.LogPanic(err)

	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully update media").
		Body("media", mediaRes),
	)
}

// Delete deletes media by request form data id
func (controller MediaController) Delete(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[common_reqs.UUID](r.Form)
	errorh.Panic(err)

	// Delete media by id
	err = controller.service.Delete(r.Context(), input)
	errorh.Panic(err)

	errorh.LogPanic(err)

	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully delete media"),
	)
}
