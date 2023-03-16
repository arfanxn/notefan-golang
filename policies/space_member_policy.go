package policies

import (
	"context"

	"github.com/google/uuid"
	perm_names "github.com/notefan-golang/enums/permission/names"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/contexth"
	"github.com/notefan-golang/models/requests/space_member_reqs"
	"github.com/notefan-golang/repositories"
)

type SpaceMemberPolicy struct {
	permissionRepository *repositories.PermissionRepository
	ursRepository        *repositories.UserRoleSpaceRepository
}

// NewSpaceMemberPolicy instantiates a new SpaceMemberPolicy
func NewSpaceMemberPolicy(
	permissionRepository *repositories.PermissionRepository,
	ursRepository *repositories.UserRoleSpaceRepository,
) *SpaceMemberPolicy {
	return &SpaceMemberPolicy{
		permissionRepository: permissionRepository,
		ursRepository:        ursRepository,
	}
}

// Get policy
func (policy *SpaceMemberPolicy) Get(ctx context.Context, input space_member_reqs.Get) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)

	// return error if not provided
	if (userId == "") || (input.SpaceId == "") {
		return exceptions.HTTPActionUnauthorized
	}

	// Get current logged in user's URS to get info about role
	ursEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.SpaceId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}

	// check if the user has permission to access the Space's Member
	permissionName := perm_names.SpaceMemberView
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
func (policy *SpaceMemberPolicy) Find(ctx context.Context, input space_member_reqs.Action) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)

	// return error if not provided
	if (userId == "") || (input.SpaceId == "") {
		return exceptions.HTTPActionUnauthorized
	}

	// Check if member is member of the given space id
	memberUrsEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, input.MemberId, input.SpaceId)
	if err != nil {
		return
	}
	// return error if member is not the member of the space
	if memberUrsEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}

	// --------------------------------------------------------

	// Get current logged in user's URS to get info about role
	ursEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.SpaceId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}

	// check if the user has permission to access the Space's Member
	permissionName := perm_names.SpaceMemberView
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

// Invite policy
func (policy *SpaceMemberPolicy) Invite(ctx context.Context, input space_member_reqs.Invite) (err error) {
	// current auth user id
	userMap := contexth.GetAuthUser(ctx)
	userId, _ := userMap["id"].(string)
	userEmail, _ := userMap["email"].(string)

	// return error if not provided
	if (userId == "") || (input.SpaceId == "") {
		return exceptions.HTTPActionUnauthorized
	}
	// return error if trying to invite self
	if userEmail == input.Email {
		return exceptions.HTTPActionUnauthorized
	}

	// Get current logged in user's URS to get info about role
	ursEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.SpaceId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}

	// check if the user has permission to Invite a Member into a Space
	permissionName := perm_names.SpaceMemberInvite
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

// UpdateRole policy
func (policy *SpaceMemberPolicy) UpdateRole(ctx context.Context, input space_member_reqs.UpdateRole) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)

	// return error if not provided
	if (userId == "") || (input.SpaceId == "") {
		return exceptions.HTTPActionUnauthorized
	}

	// return error if trying to edit self role
	if userId == input.MemberId {
		return exceptions.HTTPActionUnauthorized
	}

	// Check if member is member of the given space id
	memberUrsEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, input.MemberId, input.SpaceId)
	if err != nil {
		return
	}
	// return error if member is not the member of the space
	if memberUrsEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}

	// --------------------------------------------------------

	// Get current logged in user's URS to get info about role
	ursEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.SpaceId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}

	// check if the user has permission to update Space Member's role
	permissionName := perm_names.SpaceMemberUpdateRole
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

// Remove policy
func (policy *SpaceMemberPolicy) Remove(ctx context.Context, input space_member_reqs.Action) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)

	// return error if not provided
	if (userId == "") || (input.SpaceId == "") {
		return exceptions.HTTPActionUnauthorized
	}
	// return error if trying to remove self
	if userId == input.MemberId {
		return exceptions.HTTPActionUnauthorized
	}

	// Get current logged in user's URS to get info about role
	ursEty, err := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.SpaceId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}

	// check if the user has permission to remove a Member from a Space
	permissionName := perm_names.SpaceMemberRemove
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
