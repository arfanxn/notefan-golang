package services

import (
	"context"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/responses"
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
func (service *UserService) Find(ctx context.Context, id string) (responses.User, error) {
	user := responses.User{}

	userEntity, err := service.repository.Find(ctx, id)
	if err != nil { // err not nil return exception HTTPNotFound
		helper.ErrorLog(err)
		return user, exceptions.HTTPNotFound
	}
	user = responses.NewUserFromEntity(userEntity)

	// TODO: load user with avatar (profile_picture)
	// // User avatar
	// avatar, err := service.mediaRepository.FindByModelAndCollectionName(
	// 	ctx,
	// 	helper.ReflectGetTypeName(userEntity),
	// 	userEntity.Id.String(),
	// 	"avatar",
	// )
	// if errors.Is(err, exceptions.HTTPNotFound) { // if user doesn't have avatar return user without avatar
	// 	return user, nil
	// }

	// user.Avatar = responses.NewMediaFromEntity(avatar) // assign user avatar with avatar
	return user, nil
}
