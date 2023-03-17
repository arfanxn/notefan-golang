package controllers

import (
	"net/http"

	media_collnames "github.com/notefan-golang/enums/media/collection_names"
	"github.com/notefan-golang/helpers/combh"
	"github.com/notefan-golang/helpers/decodeh"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/helpers/rwh"
	"github.com/notefan-golang/helpers/validationh"
	"github.com/notefan-golang/models/requests/file_reqs"
	"github.com/notefan-golang/models/requests/page_reqs"
	"github.com/notefan-golang/models/responses"
	"github.com/notefan-golang/policies"
	"github.com/notefan-golang/services"
)

type PageController struct {
	service *services.PageService
	policy  *policies.PagePolicy
}

func NewPageController(
	service *services.PageService,
	policy *policies.PagePolicy,
) *PageController {
	return &PageController{
		service: service,
		policy:  policy,
	}
}

// Get gets Pages by request form data
func (controller PageController) Get(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[page_reqs.GetBySpace](r.Form)
	errorh.Panic(err)

	err = controller.policy.Get(r.Context(), input)
	errorh.Panic(err)

	pagePagination, err := controller.service.GetBySpace(r.Context(), input)
	errorh.Panic(err)

	pagePagination.SetPage(input.PerPage, input.Page, nullh.NullInt())
	pagePagination.SetURL(r.URL)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully retrieve pages of space").
			Body("pages", pagePagination),
	)
}

// Find finds a Page by request form data
func (controller PageController) Find(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[page_reqs.Action](r.Form)
	errorh.Panic(err)

	err = controller.policy.Find(r.Context(), input)
	errorh.Panic(err)

	pageRes, err := controller.service.Find(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully retrieve page").
			Body("page", pageRes),
	)
}

// Create creates a Page from given request form data
func (controller PageController) Create(w http.ResponseWriter, r *http.Request) {
	input, err := decodeh.FormData[page_reqs.Create](r.Form)
	errorh.Panic(err)

	if iconFH, _ := rwh.RequestFormFileHeader(r, media_collnames.Icon); iconFH != nil {
		fileReq, err := file_reqs.NewFromFH(iconFH)
		errorh.Panic(err)
		input.Icon = fileReq
	}
	if coverFH, _ := rwh.RequestFormFileHeader(r, media_collnames.Cover); coverFH != nil {
		fileReq, err := file_reqs.NewFromFH(coverFH)
		errorh.Panic(err)
		input.Cover = fileReq
	}

	err = validationh.ValidateStruct(input)
	errorh.Panic(err)

	err = controller.policy.Create(r.Context(), input)
	errorh.Panic(err)

	pageRes, err := controller.service.Create(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully create page").
			Body("page", pageRes),
	)
}

// Update updates a Page by request form data
func (controller PageController) Update(w http.ResponseWriter, r *http.Request) {
	input, err := decodeh.FormData[page_reqs.Update](r.Form)
	errorh.Panic(err)

	if iconFH, _ := rwh.RequestFormFileHeader(r, media_collnames.Icon); iconFH != nil {
		fileReq, err := file_reqs.NewFromFH(iconFH)
		errorh.Panic(err)
		input.Icon = fileReq
	}
	if coverFH, _ := rwh.RequestFormFileHeader(r, media_collnames.Cover); coverFH != nil {
		fileReq, err := file_reqs.NewFromFH(coverFH)
		errorh.Panic(err)
		input.Cover = fileReq
	}

	err = validationh.ValidateStruct(input)
	errorh.Panic(err)

	err = controller.policy.Update(r.Context(), input)
	errorh.Panic(err)

	pageRes, err := controller.service.Update(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully update page").
			Body("page", pageRes),
	)
}

// Delete deletes a Page by space id and page id
func (controller PageController) Delete(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[page_reqs.Action](r.Form)
	errorh.Panic(err)

	err = controller.policy.Delete(r.Context(), input)
	errorh.Panic(err)

	// Delete page by id
	err = controller.service.Delete(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully delete page"),
	)
}
