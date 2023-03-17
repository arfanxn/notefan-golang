package policies

import (
	"context"
	"sync"

	"github.com/google/uuid"
	perm_names "github.com/notefan-golang/enums/permission/names"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/chanh"
	"github.com/notefan-golang/helpers/contexth"
	"github.com/notefan-golang/helpers/synch"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/page_reqs"
	"github.com/notefan-golang/repositories"
)

type PagePolicy struct {
	repository           *repositories.PageRepository
	permissionRepository *repositories.PermissionRepository
	ursRepository        *repositories.UserRoleSpaceRepository
	waitGroup            *sync.WaitGroup
	mutex                sync.Mutex
}

// NewPagePolicy instantiates a new PagePolicy
func NewPagePolicy(
	repository *repositories.PageRepository,
	permissionRepository *repositories.PermissionRepository,
	ursRepository *repositories.UserRoleSpaceRepository,
) *PagePolicy {
	return &PagePolicy{
		repository:           repository,
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
	var (
		userId  = contexth.GetAuthUserId(ctx)
		ursEty  entities.UserRoleSpace
		errChan = synch.MakeChanWithValue[error](nil, 1)
	)
	defer close(errChan)
	// return error if no provided
	if (userId == "") || (input.SpaceId == "") || (input.PageId == "") {
		return exceptions.HTTPActionUnauthorized
	}

	// URS related operations
	ursEty, err = policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.SpaceId)
	if err != nil {
		return
	}
	// return error not found if not found
	if ursEty.UserId == uuid.Nil {
		err = exceptions.HTTPNotFound
		return
	}

	// Page related operations
	pageEty, err := policy.repository.Find(ctx, input.PageId)
	if err != nil {
		return
	}
	// return error not found if not found
	if pageEty.Id == uuid.Nil {
		err = exceptions.HTTPNotFound
		return
	}
	// return error if Page is not child of Space
	if pageEty.SpaceId.String() != input.SpaceId {
		err = exceptions.HTTPActionUnauthorized
		return
	}

	// Permission related operations
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
	var (
		userId  = contexth.GetAuthUserId(ctx)
		ursEty  entities.UserRoleSpace
		errChan = synch.MakeChanWithValue[error](nil, 1)
	)
	defer close(errChan)
	// return error if no provided
	if (userId == "") || (input.SpaceId == "") || (input.PageId == "") {
		return exceptions.HTTPActionUnauthorized
	}
	policy.waitGroup.Add(2)
	go func() { // goroutine for get UserRoleSpace
		defer policy.waitGroup.Done()
		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		ety, errChanVal := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.SpaceId)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		// return error not found if not found
		if ety.UserId == uuid.Nil {
			chanh.ReplaceVal(errChan, error(exceptions.HTTPNotFound))
			return
		}
		policy.mutex.Lock()
		defer policy.mutex.Unlock()
		ursEty = ety
	}()
	go func() { // goroutine for get Page
		defer policy.waitGroup.Done()
		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		// get page by page id only (not with space id)
		ety, errChanVal := policy.repository.Find(ctx, input.PageId)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		// return error not found if not found
		if ety.Id == uuid.Nil {
			chanh.ReplaceVal(errChan, error(exceptions.HTTPNotFound))
			return
		}
		// return error if Page is not child of Space
		if ety.SpaceId.String() != input.SpaceId {
			chanh.ReplaceVal(errChan, error(exceptions.HTTPActionUnauthorized))
			return
		}
	}()
	policy.waitGroup.Wait()
	// check if the user has permission to update the Page
	permissionName := perm_names.PageUpdate
	permission, err := policy.permissionRepository.
		FindByNameAndRoleId(ctx, permissionName, ursEty.RoleId.String())
	if err != nil {
		return
	}
	err = <-errChan
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
	var (
		userId  = contexth.GetAuthUserId(ctx)
		ursEty  entities.UserRoleSpace
		errChan = synch.MakeChanWithValue[error](nil, 1)
	)
	defer close(errChan)
	// return error if no provided
	if (userId == "") || (input.SpaceId == "") || (input.PageId == "") {
		return exceptions.HTTPActionUnauthorized
	}
	policy.waitGroup.Add(2)
	go func() { // goroutine for get UserRoleSpace
		defer policy.waitGroup.Done()
		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		ety, errChanVal := policy.ursRepository.FindByUserIdAndSpaceId(ctx, userId, input.SpaceId)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		// return error not found if not found
		if ety.UserId == uuid.Nil {
			errChan <- exceptions.HTTPNotFound
			return
		}
		policy.mutex.Lock()
		defer policy.mutex.Unlock()
		ursEty = ety
	}()
	go func() { // goroutine for get Page
		defer policy.waitGroup.Done()
		errChanVal := synch.GetChanValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		// get page by page id only (not with space id)
		ety, errChanVal := policy.repository.Find(ctx, input.PageId)
		if errChanVal != nil {
			errChan <- errChanVal
			return
		}
		// return error not found if not found
		if ety.Id == uuid.Nil {
			errChan <- exceptions.HTTPNotFound
			return
		}
		// return error if Page is not child of Space
		if ety.SpaceId.String() != input.SpaceId {
			errChan <- exceptions.HTTPActionUnauthorized
			return
		}
	}()
	policy.waitGroup.Wait()
	// check if the user has permission to update the Page
	permissionName := perm_names.PageDelete
	permission, err := policy.permissionRepository.
		FindByNameAndRoleId(ctx, permissionName, ursEty.RoleId.String())
	if err != nil {
		return
	}
	err = <-errChan
	if err != nil {
		return
	}
	if permission.Id == uuid.Nil {
		return exceptions.HTTPActionUnauthorized
	}
	return nil
}
