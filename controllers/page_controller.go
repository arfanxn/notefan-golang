package controllers

import (
	"fmt"
	"net/http"

	"github.com/notefan-golang/helpers/combh"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/helpers/rwh"
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

	fmt.Println("pagePagination", pagePagination.Items[0].Title)

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
