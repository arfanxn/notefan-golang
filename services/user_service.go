package services

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	media_collnames "github.com/notefan-golang/enums/media/collection_names"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/user_reqs"
	"github.com/notefan-golang/models/responses/media_ress"
	"github.com/notefan-golang/models/responses/user_ress"
	"github.com/notefan-golang/repositories"
)

type UserService struct {
	repository      *repositories.UserRepository
	mediaRepository *repositories.MediaRepository
	waitGroup       *sync.WaitGroup
	mutex           sync.Mutex
}

func NewUserService(
	repository *repositories.UserRepository,
	mediaRepository *repositories.MediaRepository,
) *UserService {
	return &UserService{
		repository:      repository,
		mediaRepository: mediaRepository,
		waitGroup:       new(sync.WaitGroup),
		mutex:           sync.Mutex{},
	}
}

// Find Finds user by id and returns user with its relations
func (service *UserService) Find(ctx context.Context, id string) (
	userRes user_ress.User, err error) {
	userEty, err := service.repository.Find(ctx, id)
	if err != nil {
		return
	}
	if userEty.Id == uuid.Nil { // err not nil return exception HTTPNotFound
		err = exceptions.HTTPNotFound
		return
	}
	userRes = user_ress.FillFromEntity(userEty)

	// Get User avatar
	avatarMediaEty, err := service.mediaRepository.FindByModelAndCollectionName(
		ctx,
		reflecth.GetTypeName(userEty),
		userEty.Id.String(),
		media_collnames.Avatar,
	)
	if avatarMediaEty.Id == uuid.Nil { // if user doesn't have avatar return user without avatar
		return userRes, nil
	}

	userRes.Avatar = media_ress.FillFromEntity(avatarMediaEty) // assign user avatar with avatar
	return userRes, nil
}

// UpdateProfile updates the given user update request
func (service *UserService) UpdateProfile(ctx context.Context, data user_reqs.UpdateProfile) (
	userRes user_ress.User, err error) {
	var (
		userEty        entities.User
		avatarMediaEty entities.Media
	)

	// Get User entity
	userEty, err = service.repository.Find(ctx, data.Id)
	if err != nil {
		return
	}
	if userEty.Id == uuid.Nil { // return exceptions.HTTPNotFound ig userEty not found
		err = exceptions.HTTPNotFound
		return
	}
	// Convert request to entity
	userEty.Name = data.Name

	// Update User Entity
	_, err = service.repository.UpdateById(ctx, &userEty)
	if err != nil {
		return
	}
	userRes = user_ress.FillFromEntity(userEty)

	// if avatar/file is present and provided save avatar
	if data.Avatar != nil && data.Avatar.IsProvided() {
		// Get Media entity (user's avatar) for update operation
		avatarMediaEty, err = service.mediaRepository.FindByModelAndCollectionName(
			ctx,
			reflecth.GetTypeName(userEty),
			data.Id,
			media_collnames.Avatar,
		)
		if err != nil {
			return
		}

		// Assign new file to the media file
		avatarMediaEty.FileName = filepath.Base(data.Avatar.Name)
		avatarMediaEty.File = data.Avatar

		// Do update
		_, err = service.mediaRepository.UpdateById(ctx, &avatarMediaEty)
		if err != nil {
			return
		}

		// Assign the media response to user avatar response
		userRes.Avatar = media_ress.FillFromEntity(avatarMediaEty)
	}

	return userRes, nil
}
