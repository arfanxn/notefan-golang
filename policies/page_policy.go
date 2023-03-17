package policies

import (
	"context"
	"sync"

	"github.com/google/uuid"
	perm_names "github.com/notefan-golang/enums/permission/names"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/contexth"
	"github.com/notefan-golang/models/requests/page_reqs"
	"github.com/notefan-golang/repositories"
)

type PagePolicy struct {
	permissionRepository *repositories.PermissionRepository
	ursRepository        *repositories.UserRoleSpaceRepository
	waitGroup            *sync.WaitGroup
	mutex                sync.Mutex
}

// NewPagePolicy instantiates a new PagePolicy
func NewPagePolicy(
	permissionRepository *repositories.PermissionRepository,
	ursRepository *repositories.UserRoleSpaceRepository,
) *PagePolicy {
	return &PagePolicy{
		permissionRepository: permissionRepository,
		ursRepository:        ursRepository,
		waitGroup:            new(sync.WaitGroup),
		mutex:                sync.Mutex{},
	}
}

// Get policy
func (policy *PagePolicy) Get(ctx context.Context, input page_reqs.GetBySpace) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if (userId == "") || (input.SpaceId == "") {
		return exceptions.HTTPActionUnauthorized
	}
	ursEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.SpaceId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}
	// check if the user has permission to access the Page
	permissionName := perm_names.PageView
	permission, err := policy.permissionRepository.
		FindByNameAndRoleId(ctx, permissionName, ursEty.RoleId.String())
	if err != nil {
		return
	}
	if permission.Id == uuid.Nil {
		return exceptions.HTTPActionUnauthorized
	}
	return nil
}

// Find policy
func (policy *PagePolicy) Find(ctx context.Context, input page_reqs.Action) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if (userId == "") || (input.SpaceId == "") || (input.PageId == "") {
		return exceptions.HTTPActionUnauthorized
	}
	// URS related operations
	ursEty, err := policy.ursRepository.FindByUserIdAndPageId(ctx, userId, input.PageId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		err = exceptions.HTTPNotFound
		return
	}
	// check if the user has permission to access the Page
	permissionName := perm_names.PageView
	permission, err := policy.permissionRepository.
		FindByNameAndRoleId(ctx, permissionName, ursEty.RoleId.String())
	if err != nil {
		return
	}
	if permission.Id == uuid.Nil {
		return exceptions.HTTPActionUnauthorized
	}
	return nil
}

// Create policy
func (policy *PagePolicy) Create(ctx context.Context, input page_reqs.Create) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if (userId == "") || (input.SpaceId == "") {
		return exceptions.HTTPActionUnauthorized
	}
	ursEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.SpaceId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}
	// check if the user has permission to create Page
	permissionName := perm_names.PageCreate
	permission, err := policy.permissionRepository.
		FindByNameAndRoleId(ctx, permissionName, ursEty.RoleId.String())
	if err != nil {
		return
	}
	if permission.Id == uuid.Nil {
		return exceptions.HTTPActionUnauthorized
	}
	return nil
}

// Update policy
func (policy *PagePolicy) Update(ctx context.Context, input page_reqs.Update) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if (userId == "") || (input.SpaceId == "") || (input.PageId == "") {
		return exceptions.HTTPActionUnauthorized
	}
	// URS related operations
	ursEty, err := policy.ursRepository.FindByUserIdAndPageId(ctx, userId, input.PageId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		err = exceptions.HTTPNotFound
		return
	}
	// check if the user has permission to update Page
	permissionName := perm_names.PageUpdate
	permission, err := policy.permissionRepository.
		FindByNameAndRoleId(ctx, permissionName, ursEty.RoleId.String())
	if err != nil {
		return
	}
	if permission.Id == uuid.Nil {
		return exceptions.HTTPActionUnauthorized
	}
	return nil
}

// Delete policy
func (policy *PagePolicy) Delete(ctx context.Context, input page_reqs.Action) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if (userId == "") || (input.SpaceId == "") || (input.PageId == "") {
		return exceptions.HTTPActionUnauthorized
	}
	// URS related operations
	ursEty, err := policy.ursRepository.FindByUserIdAndPageId(ctx, userId, input.PageId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		err = exceptions.HTTPNotFound
		return
	}
	// check if the user has permission to delete Page
	permissionName := perm_names.PageDelete
	permission, err := policy.permissionRepository.
		FindByNameAndRoleId(ctx, permissionName, ursEty.RoleId.String())
	if err != nil {
		return
	}
	if permission.Id == uuid.Nil {
		return exceptions.HTTPActionUnauthorized
	}
	return nil
}
