package policies

import (
	"context"

	"github.com/google/uuid"
	perm_names "github.com/notefan-golang/enums/permission/names"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/contexth"
	"github.com/notefan-golang/models/requests/common_reqs"
	"github.com/notefan-golang/models/requests/space_reqs"
	"github.com/notefan-golang/repositories"
)

type SpacePolicy struct {
	permissionRepository *repositories.PermissionRepository
	ursRepository        *repositories.UserRoleSpaceRepository
}

// NewSpacePolicy instantiates a new SpacePolicy
func NewSpacePolicy(
	permissionRepository *repositories.PermissionRepository,
	ursRepository *repositories.UserRoleSpaceRepository,
) *SpacePolicy {
	return &SpacePolicy{
		permissionRepository: permissionRepository,
		ursRepository:        ursRepository,
	}
}

// Find policy
func (policy *SpacePolicy) Find(ctx context.Context, input common_reqs.UUID) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if (userId == "") || (input.Id == "") {
		return exceptions.HTTPActionUnauthorized
	}
	ursEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.Id)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}
	// check if the user has permission to access Space
	permissionName := perm_names.SpaceView
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

// Get policy
func (policy *SpacePolicy) Get(ctx context.Context, input space_reqs.GetByUser) (err error) {
	return nil // every user has permission to get spaces
}

// Create policy
func (policy *SpacePolicy) Create(ctx context.Context, input space_reqs.Create) (err error) {
	return nil // every user has permission to create space
}

// Update policy
func (policy *SpacePolicy) Update(ctx context.Context, input space_reqs.Update) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if (userId == "") || (input.Id == "") {
		return exceptions.HTTPActionUnauthorized
	}
	ursEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.Id)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}
	// check if the user has permission to update Space
	permissionName := perm_names.SpaceUpdate
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
func (policy *SpacePolicy) Delete(ctx context.Context, input common_reqs.UUID) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)

	// return error if no provided
	if (userId == "") || (input.Id == "") {
		return exceptions.HTTPActionUnauthorized
	}

	ursEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.Id)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}

	// check if the user has permission to delete Space
	permissionName := perm_names.SpaceDelete
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
