package services

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/requests/user_reqs"
	"github.com/notefan-golang/models/responses/media_ress"
	"github.com/notefan-golang/models/responses/user_ress"
	"github.com/notefan-golang/repositories"
)

type UserService struct {
	repository      *repositories.UserRepository
	mediaRepository *repositories.MediaRepository
}

func NewUserService(
	repository *repositories.UserRepository,
	mediaRepository *repositories.MediaRepository,
) *UserService {
	return &UserService{repository: repository, mediaRepository: mediaRepository}
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
		"avatar",
	)
	if errors.Is(err, exceptions.HTTPNotFound) { // if user doesn't have avatar return user without avatar
		return user, nil
	}

	user.Avatar = media_ress.FillFromEntity(avatar) // assign user avatar with avatar
	return user, nil
}

// UpdateProfile updates the given user update request
func (service *UserService) UpdateProfile(ctx context.Context, data user_reqs.UpdateProfile) (user_ress.User, error) {
	userRes := user_ress.User{}

	// Get User entity
	userEty, err := service.repository.Find(ctx, data.Id)
	errorh.Panic(err) // panic if not found

	// Convert request to entity
	userEty.Name = data.Name

	// Update User Entity
	_, err = service.repository.UpdateById(ctx, &userEty)
	errorh.LogPanic(err) // panic and log the error if error at update
	userRes = user_ress.FillFromEntity(userEty)

	// check if avatar exists, if exists do update user avatar operation
	if data.Avatar.Size > 0 {
		// Get Media entity (user's avatar) for update operation
		mediaEty, _ := service.mediaRepository.FindByModelAndCollectionName(
			ctx,
			reflecth.GetTypeName(userEty),
			userEty.Id.String(),
			"avatar",
		)

		// Assign new file to the media file
		mediaEty.FileName = filepath.Base(data.Avatar.Name)
		mediaEty.File = bytes.NewBuffer(data.Avatar.Buffer.Bytes())

		// Do update
		_, err = service.mediaRepository.UpdateById(ctx, &mediaEty)
		errorh.LogPanic(err) // panic and log the error if error at update

		// Assign the media response to user avatar response
		userRes.Avatar = media_ress.FillFromEntity(mediaEty)
	}

	return userRes, nil
}
