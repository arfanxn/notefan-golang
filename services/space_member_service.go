package services

import (
	"context"
	"sync"

	"github.com/google/uuid"
	media_coll_names "github.com/notefan-golang/enums/media/collection_names"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/sliceh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/space_member_reqs"
	"github.com/notefan-golang/models/responses/media_ress"
	"github.com/notefan-golang/models/responses/pagination_ress"
	"github.com/notefan-golang/models/responses/user_ress"
	"github.com/notefan-golang/repositories"
)

type SpaceMemberService struct {
	repository      *repositories.SpaceMemberRepository
	ursRepository   *repositories.UserRoleSpaceRepository
	mediaRepository *repositories.MediaRepository
	waitGroup       *sync.WaitGroup
	mutex           sync.RWMutex
}

func NewSpaceMemberService(
	repository *repositories.SpaceMemberRepository,
	ursRepository *repositories.UserRoleSpaceRepository,
	mediaRepository *repositories.MediaRepository,
) *SpaceMemberService {
	return &SpaceMemberService{
		repository:      repository,
		ursRepository:   ursRepository,
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
	service.repository.Query.AddWith("members")
	memberUserEtys, err := service.repository.GetBySpaceId(ctx, data.SpaceId)
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

	return paginationRes, nil
}

// Find get spaces by user id and parse it to slice of Space Response
func (service *SpaceMemberService) Find(ctx context.Context, data space_member_reqs.Action) (
	memberRes user_ress.User, err error) {
	memberUserEty, err := service.repository.FindByMemberIdAndSpaceId(ctx, data.MemberId, data.Id)
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
