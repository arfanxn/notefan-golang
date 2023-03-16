package controllers

import (
	"net/http"

	"github.com/notefan-golang/helpers/combh"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/helpers/rwh"
	"github.com/notefan-golang/models/requests/space_member_reqs"
	"github.com/notefan-golang/models/responses"
	"github.com/notefan-golang/policies"
	"github.com/notefan-golang/services"
)

type SpaceMemberController struct {
	service *services.SpaceMemberService
	policy  *policies.SpaceMemberPolicy
}

func NewSpaceMemberController(
	service *services.SpaceMemberService,
	policy *policies.SpaceMemberPolicy,
) *SpaceMemberController {
	return &SpaceMemberController{
		service: service,
		policy:  policy,
	}
}

// Get returns members of space
func (controller SpaceMemberController) Get(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[space_member_reqs.Get](r.Form)
	errorh.Panic(err)

	err = controller.policy.Get(r.Context(), input)
	errorh.Panic(err)

	spacePagination, err := controller.service.Get(r.Context(), input)
	errorh.Panic(err)

	spacePagination.SetPage(input.PerPage, input.Page, nullh.NullInt())
	spacePagination.SetURL(r.URL)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully retrieve members of space").
			Body("members", spacePagination),
	)
}

// Find returns a member of space
func (controller SpaceMemberController) Find(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[space_member_reqs.Action](r.Form)
	errorh.Panic(err)

	err = controller.policy.Find(r.Context(), input)
	errorh.Panic(err)

	memberRes, err := controller.service.Find(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully retrieve member of space").
			Body("member", memberRes),
	)
}

// Invite invites a Member into Space
func (controller SpaceMemberController) Invite(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[space_member_reqs.Invite](r.Form)
	errorh.Panic(err)

	err = controller.policy.Invite(r.Context(), input)
	errorh.Panic(err)

	err = controller.service.Invite(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully invite member into space"))
}

// Invite invites a Member into Space
func (controller SpaceMemberController) UpdateRole(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[space_member_reqs.UpdateRole](r.Form)
	errorh.Panic(err)

	err = controller.policy.UpdateRole(r.Context(), input)
	errorh.Panic(err)

	memberRes, err := controller.service.UpdateRole(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully update member role").
			Body("member", memberRes))
}

// Invite invites a Member into Space
func (controller SpaceMemberController) Remove(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[space_member_reqs.Action](r.Form)
	errorh.Panic(err)

	err = controller.policy.Remove(r.Context(), input)
	errorh.Panic(err)

	err = controller.service.Remove(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully remove member from space"))
}
