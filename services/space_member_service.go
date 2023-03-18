package services

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	media_coll_names "github.com/notefan-golang/enums/media/collection_names"
	role_names "github.com/notefan-golang/enums/role/names"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/chanh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/sliceh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"
	"github.com/notefan-golang/models/requests/space_member_reqs"
	"github.com/notefan-golang/models/responses/media_ress"
	"github.com/notefan-golang/models/responses/pagination_ress"
	"github.com/notefan-golang/models/responses/role_ress"
	"github.com/notefan-golang/models/responses/user_ress"
	"github.com/notefan-golang/repositories"
)

type SpaceMemberService struct {
	repository      *repositories.SpaceMemberRepository
	userRepository  *repositories.UserRepository
	ursRepository   *repositories.UserRoleSpaceRepository
	roleRepository  *repositories.RoleRepository
	mediaRepository *repositories.MediaRepository
	waitGroup       *sync.WaitGroup
	mutex           sync.RWMutex
}

func NewSpaceMemberService(
	repository *repositories.SpaceMemberRepository,
	userRepository *repositories.UserRepository,
	ursRepository *repositories.UserRoleSpaceRepository,
	roleRepository *repositories.RoleRepository,
	mediaRepository *repositories.MediaRepository,
) *SpaceMemberService {
	return &SpaceMemberService{
		repository:      repository,
		userRepository:  userRepository,
		ursRepository:   ursRepository,
		roleRepository:  roleRepository,
		mediaRepository: mediaRepository,
		waitGroup:       new(sync.WaitGroup),
		mutex:           sync.RWMutex{},
	}
}

// GetByUser get spaces by user id and parse it to slice of Space Response
func (service *SpaceMemberService) Get(ctx context.Context, data space_member_reqs.Get) (
	paginationRes pagination_ress.Pagination[user_ress.User], err error) {
	service.repository.Query.Limit = data.PerPage
	service.repository.Query.Offset = (data.Page - 1) * int64(data.PerPage)
	service.repository.Query.Keyword = data.Keyword
	for _, orderBy := range data.OrderBys {
		keyAndVal := strings.Split(orderBy, "=")
		service.repository.Query.AddOrderBy(keyAndVal[0], keyAndVal[1])
	}
	memberUserEtys, err := service.repository.GetBySpaceId(ctx, data.SpaceId)
	service.repository.Query = query_reqs.Default() // reset query to default after retrieving
	if err != nil {
		return
	}

	// Get Avatar for each Members
	memberAvatarMediaEtys, err := service.mediaRepository.
		GetByModelsAndCollectionNames(ctx,
			sliceh.Map(memberUserEtys, func(userEty entities.User) entities.Media {
				return entities.Media{
					ModelType:      reflecth.GetTypeName(userEty),
					ModelId:        userEty.Id,
					CollectionName: media_coll_names.Avatar,
				}
			})...,
		)

	for index, userEty := range memberUserEtys {
		userRes := user_ress.FillFromEntity(userEty)
		paginationRes.Items = append(paginationRes.Items, userRes)

		service.waitGroup.Add(1)

		go func(userRes *user_ress.User) {
			defer service.waitGroup.Done()
			avatarMediaEtys := sliceh.Filter(memberAvatarMediaEtys, func(media entities.Media) bool {
				return userRes.Id == media.ModelId.String()
			})
			if len(avatarMediaEtys) == 0 {
				return
			}

			service.mutex.Lock()
			defer service.mutex.Unlock()
			userRes.Avatar = media_ress.FillFromEntity(avatarMediaEtys[0])
		}(&paginationRes.Items[index])
	}

	service.waitGroup.Wait()
	paginationRes.Total = len(paginationRes.Items)

	return paginationRes, nil
}

// Find get spaces by user id and parse it to slice of Space Response
func (service *SpaceMemberService) Find(ctx context.Context, data space_member_reqs.Action) (
	memberRes user_ress.User, err error) {
	memberUserEty, err := service.repository.FindByMemberIdAndSpaceId(ctx, data.MemberId, data.SpaceId)
	if err != nil {
		return
	}
	memberRes = user_ress.FillFromEntity(memberUserEty)

	// Load Avatar for member if exists
	memberAvatarMediaEty, err := service.mediaRepository.
		FindByModelAndCollectionName(ctx,
			reflecth.GetTypeName(memberUserEty), memberUserEty.Id.String(), media_coll_names.Avatar)
	if err != nil {
		return
	}
	if memberAvatarMediaEty.Id != uuid.Nil {
		memberRes.Avatar = media_ress.FillFromEntity(memberAvatarMediaEty)
	}

	return
}

// Invite invites a member to join a Space
func (service *SpaceMemberService) Invite(ctx context.Context, data space_member_reqs.Invite) (
	err error) {
	// Member related
	memberUserEty, err := service.userRepository.FindByEmail(ctx, data.Email)
	if err != nil {
		return
	}
	if memberUserEty.Id == uuid.Nil {
		err = exceptions.HTTPNotFound
		return
	}

	// Check if Member/User has already in Space or not
	ursEty, err := service.ursRepository.FindByUserIdAndSpaceId(ctx, memberUserEty.Id.String(), data.SpaceId)
	if ursEty.UserId == memberUserEty.Id { // if member is already in space return validation error
		err = exceptions.NewHTTPError(http.StatusUnprocessableEntity,
			exceptions.NewValidationError("email", data.Email+" has already in Space"))
		return
	}

	// Role related
	roleEty, err := service.roleRepository.FindByName(ctx, role_names.SpaceMember)
	if err != nil {
		return
	}
	if roleEty.Id == uuid.Nil {
		err = exceptions.HTTPNotFound
		return
	}

	// Create URS to associate between User/Member and Role and Space
	ursEty = entities.UserRoleSpace{
		UserId:  memberUserEty.Id,
		SpaceId: uuid.MustParse(data.SpaceId),
		RoleId:  roleEty.Id,
	}
	_, err = service.ursRepository.Create(ctx, &ursEty)
	if err != nil {
		return
	}

	err = nil
	return
}

// UpdateRole updates Member's role from the given Space
func (service *SpaceMemberService) UpdateRole(ctx context.Context, data space_member_reqs.UpdateRole) (
	memberUserRes user_ress.User, err error) {
	var (
		memberUserEty  entities.User
		ursEty         entities.UserRoleSpace
		roleEty        entities.Role
		avatarMediaEty entities.Media
		errChan        = chanh.Make[error](nil, 1)
	)
	defer close(errChan)
	service.waitGroup.Add(4)
	go func() { // goroutine for load User
		defer service.waitGroup.Done()
		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		ety, errChanVal := service.userRepository.Find(ctx, data.MemberId)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		if ety.Id == uuid.Nil {
			errChan <- exceptions.HTTPNotFound
			return
		}
		service.mutex.Lock()
		defer service.mutex.Unlock()
		memberUserEty = ety
	}()
	go func() { // goroutine for load UserRoleSpace
		defer service.waitGroup.Done()
		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		ety, errChanVal := service.ursRepository.FindByUserIdAndSpaceId(ctx, data.MemberId, data.SpaceId)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		if ety.Id == uuid.Nil {
			errChan <- exceptions.HTTPNotFound
			return
		}
		service.mutex.Lock()
		defer service.mutex.Unlock()
		ursEty = ety
	}()
	go func() { // goroutine for load Role
		defer service.waitGroup.Done()
		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		ety, errChanVal := service.roleRepository.FindByName(ctx, data.RoleName)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		if ety.Id == uuid.Nil {
			errChan <- exceptions.HTTPNotFound
			return
		}
		service.mutex.Lock()
		defer service.mutex.Unlock()
		roleEty = ety
	}()
	go func() { // goroutine for load user/member avatar
		defer service.waitGroup.Done()
		errChanVal := chanh.GetValAndKeep(errChan)
		if errChanVal != nil {
			return
		}
		ety, errChanVal := service.mediaRepository.FindByModelAndCollectionName(
			ctx,
			reflecth.GetTypeName(entities.User{}),
			data.MemberId,
			media_coll_names.Avatar,
		)
		if errChanVal != nil {
			chanh.ReplaceVal(errChan, errChanVal)
			return
		}
		if ety.Id != uuid.Nil {
			service.mutex.Lock()
			defer service.mutex.Unlock()
			avatarMediaEty = ety
		}
	}()

	service.waitGroup.Wait()

	err = <-errChan
	if err != nil {
		return
	}

	// Update the member role
	ursEty.RoleId = roleEty.Id
	_, err = service.ursRepository.UpdateById(ctx, &ursEty)
	if err != nil {
		return
	}

	// Prepare response
	memberUserRes = user_ress.FillFromEntity(memberUserEty)
	memberUserRes.Role = role_ress.FillFromEntity(roleEty)
	if avatarMediaEty.Id != uuid.Nil {
		memberUserRes.Avatar = media_ress.FillFromEntity(avatarMediaEty)
	}
	err = nil
	return
}

// Remove removes a Member from the given Space
func (service *SpaceMemberService) Remove(ctx context.Context, data space_member_reqs.Action) (
	err error) {
	ursEty, err := service.ursRepository.FindByUserIdAndSpaceId(ctx, data.MemberId, data.SpaceId)
	if err != nil {
		return
	}
	if ursEty.UserId == uuid.Nil {
		err = exceptions.HTTPNotFound
		return
	}

	// Remove User/Member from Space
	_, err = service.ursRepository.DeleteByIds(ctx, ursEty.Id.String())
	if err != nil {
		return
	}

	err = nil
	return
}
