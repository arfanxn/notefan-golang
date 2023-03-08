package services

import (
	"context"
	"errors"
	"path/filepath"
	"sync"

	media_collnames "github.com/notefan-golang/enums/media/collection_names"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
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
func (service *UserService) Find(ctx context.Context, id string) (user_ress.User, error) {
	user := user_ress.User{}

	userEntity, err := service.repository.Find(ctx, id)
	if err != nil { // err not nil return exception HTTPNotFound
		errorh.Log(err)
		return user, exceptions.HTTPNotFound
	}
	user = user_ress.FillFromEntity(userEntity)

	// Get User avatar
	avatar, err := service.mediaRepository.FindByModelAndCollectionName(
		ctx,
		reflecth.GetTypeName(userEntity),
		userEntity.Id.String(),
		media_collnames.Avatar,
	)
	if errors.Is(err, exceptions.HTTPNotFound) { // if user doesn't have avatar return user without avatar
		return user, nil
	}

	user.Avatar = media_ress.FillFromEntity(avatar) // assign user avatar with avatar
	return user, nil
}

// UpdateProfile updates the given user update request
func (service *UserService) UpdateProfile(ctx context.Context, data user_reqs.UpdateProfile) (user_ress.User, error) {
	var (
		userEty  entities.User
		mediaEty entities.Media
		userRes  user_ress.User
		err      error
	)

	service.waitGroup.Add(2)

	go func() { // goroutine for updating user
		defer service.waitGroup.Done()

		service.mutex.Lock()
		// Get User entity
		userEty, err = service.repository.Find(ctx, data.Id)
		errorh.Panic(err) // panic if not found
		// Convert request to entity
		userEty.Name = data.Name

		// Update User Entity
		_, err = service.repository.UpdateById(ctx, &userEty)
		errorh.LogPanic(err) // panic and log the error if error at update
		userRes = user_ress.FillFromEntity(userEty)
		service.mutex.Unlock()
	}()

	go func() { // goroutine for updating user's avatar if provided
		defer service.waitGroup.Done()

		// if avatar/file is nil or not provided return immediately
		if data.Avatar == nil || !data.Avatar.IsProvided() {
			return
		}

		service.mutex.Lock()
		// Get Media entity (user's avatar) for update operation
		mediaEty, err = service.mediaRepository.FindByModelAndCollectionName(
			ctx,
			reflecth.GetTypeName(userEty),
			data.Id,
			media_collnames.Avatar,
		)
		errorh.LogPanic(err) // panic if find error

		// Assign new file to the media file
		mediaEty.FileName = filepath.Base(data.Avatar.Name)
		mediaEty.File = data.Avatar

		// Do update
		_, err = service.mediaRepository.UpdateById(ctx, &mediaEty)
		errorh.LogPanic(err) // panic and log the error if error at update

		// Assign the media response to user avatar response
		userRes.Avatar = media_ress.FillFromEntity(mediaEty)

		service.mutex.Unlock()
	}()

	service.waitGroup.Wait()

	return userRes, nil
}
