package services

import (
	"context"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	userRess "github.com/notefan-golang/models/responses/user_ress"
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
func (service *UserService) Find(ctx context.Context, id string) (userRess.User, error) {
	user := userRess.User{}

	userEntity, err := service.repository.Find(ctx, id)
	if err != nil { // err not nil return exception HTTPNotFound
		errorh.Log(err)
		return user, exceptions.HTTPNotFound
	}
	user = userRess.NewUserFromEntity(userEntity)

	// TODO: load user with avatar (profile_picture)
	// // User avatar
	// avatar, err := service.mediaRepository.FindByModelAndCollectionName(
	// 	ctx,
	// 	reflecth.GetTypeName(userEntity),
	// 	userEntity.Id.String(),
	// 	"avatar",
	// )
	// if errors.Is(err, exceptions.HTTPNotFound) { // if user doesn't have avatar return user without avatar
	// 	return user, nil
	// }

	// user.Avatar = responses.NewMediaFromEntity(avatar) // assign user avatar with avatar
	return user, nil
}
