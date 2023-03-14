package controllers

import (
	"net/http"

	"github.com/notefan-golang/helpers/combh"
	"github.com/notefan-golang/helpers/contexth"
	"github.com/notefan-golang/helpers/decodeh"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/helpers/rwh"
	"github.com/notefan-golang/helpers/validationh"
	"github.com/notefan-golang/models/requests/common_reqs"
	"github.com/notefan-golang/models/requests/file_reqs"
	"github.com/notefan-golang/models/requests/space_reqs"
	"github.com/notefan-golang/models/responses"
	"github.com/notefan-golang/policies"
	"github.com/notefan-golang/services"
)

type SpaceController struct {
	service *services.SpaceService
	policy  *policies.SpacePolicy
}

func NewSpaceController(service *services.SpaceService, policy *policies.SpacePolicy) *SpaceController {
	return &SpaceController{
		service: service,
		policy:  policy,
	}
}

func (controller SpaceController) Get(w http.ResponseWriter, r *http.Request) {
	input, err := decodeh.FormData[space_reqs.GetByUser](r.Form)
	errorh.Panic(err)
	input.UserId = contexth.GetAuthUserId(r.Context())

	// Validate input
	err = validationh.ValidateStruct(input)
	errorh.Panic(err)

	err = controller.policy.Get(r.Context(), input)
	errorh.Panic(err)

	spacePagination, err := controller.service.GetByUser(r.Context(), input)
	errorh.Panic(err)

	spacePagination.SetPage(input.PerPage, input.Page, nullh.NullInt())
	spacePagination.SetURL(r.URL)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully retrieve spaces of user").
			Body("spaces", spacePagination),
	)
}

// Find finds a space by request form data id
func (controller SpaceController) Find(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[common_reqs.UUID](r.Form)
	errorh.Panic(err)

	err = controller.policy.Find(r.Context(), input)
	errorh.Panic(err)

	spaceRes, err := controller.service.Find(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully retrieve space").
			Body("space", spaceRes),
	)
}

// Create creates space from request data
func (controller SpaceController) Create(w http.ResponseWriter, r *http.Request) {
	input, err := decodeh.FormData[space_reqs.Create](r.Form)
	errorh.Panic(err)
	input.UserId = contexth.GetAuthUserId(r.Context())

	// Get icon file header from form data
	iconFH, _ := rwh.RequestFormFileHeader(r, "icon")
	if iconFH != nil {
		fileReq, err := file_reqs.NewFromFH(iconFH)
		errorh.Panic(err)
		input.Icon = fileReq
	}

	// Validate input
	err = validationh.ValidateStruct(input)
	errorh.Panic(err)

	err = controller.policy.Create(r.Context(), input)
	errorh.Panic(err)

	// Create space
	spaceRes, err := controller.service.Create(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully create space").
		Body("space", spaceRes),
	)
}

// Update updates space by request form data id
func (controller SpaceController) Update(w http.ResponseWriter, r *http.Request) {
	input, err := decodeh.FormData[space_reqs.Update](r.Form)
	errorh.Panic(err)
	// Get icon file header from form data
	iconFH, _ := rwh.RequestFormFileHeader(r, "icon")
	if iconFH != nil {
		fileReq, err := file_reqs.NewFromFH(iconFH)
		errorh.Panic(err)
		input.Icon = fileReq
	}

	// Validate input
	err = validationh.ValidateStruct(input)
	errorh.Panic(err)

	err = controller.policy.Update(r.Context(), input)
	errorh.Panic(err)

	// Update space
	spaceRes, err := controller.service.Update(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully update space").
		Body("space", spaceRes),
	)
}

// Delete deletes media by request form data id
func (controller SpaceController) Delete(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[common_reqs.UUID](r.Form)
	errorh.Panic(err)

	err = controller.policy.Delete(r.Context(), input)
	errorh.Panic(err)

	// Delete space by id
	err = controller.service.Delete(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully delete space"),
	)
}
