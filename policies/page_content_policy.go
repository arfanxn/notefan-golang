package policies

import (
	"context"
	"sync"

	"github.com/google/uuid"
	perm_names "github.com/notefan-golang/enums/permission/names"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/contexth"
	pc_reqs "github.com/notefan-golang/models/requests/page_content_reqs"
	"github.com/notefan-golang/repositories"
)

type PageContentPolicy struct {
	permissionRepository *repositories.PermissionRepository
	ursRepository        *repositories.UserRoleSpaceRepository
	waitGroup            *sync.WaitGroup
	mutex                sync.Mutex
}

// NewPageContentPolicy instantiates a new PageContentPolicy
func NewPageContentPolicy(
	permissionRepository *repositories.PermissionRepository,
	ursRepository *repositories.UserRoleSpaceRepository,
) *PageContentPolicy {
	return &PageContentPolicy{
		permissionRepository: permissionRepository,
		ursRepository:        ursRepository,
		waitGroup:            new(sync.WaitGroup),
		mutex:                sync.Mutex{},
	}
}

/*
 * ----------------------------------------------------------------
 * Policy utility methods ⬇
 * ----------------------------------------------------------------
 */

//

/*
 * ----------------------------------------------------------------
 * Policy authorization methods ⬇
 * ----------------------------------------------------------------
 */

// Get policy
func (policy *PageContentPolicy) Get(ctx context.Context, input pc_reqs.GetByPage) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if input.PageId == "" {
		return exceptions.HTTPActionUnauthorized
	}
	// get URS
	ursEty, err := policy.ursRepository.FindByUserIdAndPageId(ctx, userId, input.PageId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}
	// check if the user has permission to access the PageContent
	permissionName := perm_names.PageContentView
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
func (policy *PageContentPolicy) Find(ctx context.Context, input pc_reqs.Action) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if input.PageId == "" {
		return exceptions.HTTPActionUnauthorized
	}
	// get URS
	ursEty, err := policy.ursRepository.FindByUserIdAndPageContentId(ctx, userId, input.PageContentId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}
	// check if the user has permission to delete PageContent
	permissionName := perm_names.PageContentDelete
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
func (policy *PageContentPolicy) Create(ctx context.Context, input pc_reqs.Create) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if input.PageId == "" {
		return exceptions.HTTPActionUnauthorized
	}
	// get URS
	ursEty, err := policy.ursRepository.FindByUserIdAndPageId(ctx, userId, input.PageId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}
	// check if the user has permission to create PageContent
	permissionName := perm_names.PageContentCreate
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
func (policy *PageContentPolicy) Update(ctx context.Context, input pc_reqs.Update) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if input.PageId == "" {
		return exceptions.HTTPActionUnauthorized
	}
	// get URS
	ursEty, err := policy.ursRepository.FindByUserIdAndPageContentId(ctx, userId, input.PageContentId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}
	// check if the user has permission to update PageContent
	permissionName := perm_names.PageContentUpdate
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
func (policy *PageContentPolicy) Delete(ctx context.Context, input pc_reqs.Action) (err error) {
	// current auth user id
	userId := contexth.GetAuthUserId(ctx)
	// return error if no provided
	if input.PageId == "" {
		return exceptions.HTTPActionUnauthorized
	}
	// get URS
	ursEty, err := policy.ursRepository.FindByUserIdAndPageContentId(ctx, userId, input.PageContentId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		return exceptions.HTTPNotFound
	}
	// check if the user has permission to delete PageContent
	permissionName := perm_names.PageContentDelete
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
