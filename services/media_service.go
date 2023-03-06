package services

import (
	"context"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/models/requests/common_reqs"
	"github.com/notefan-golang/models/responses/media_ress"
	"github.com/notefan-golang/repositories"
)

type MediaService struct {
	repository *repositories.MediaRepository
}

func NewMediaService(
	repository *repositories.MediaRepository,
) *MediaService {
	return &MediaService{repository: repository}
}

func (service *MediaService) Find(ctx context.Context, data common_reqs.UUID) (media_ress.Media, error) {
	mediaEty, err := service.repository.Find(ctx, data.Id)
	if err != nil { // err not nil return exception HTTPNotFound
		errorh.Log(err)
		return media_ress.Media{}, exceptions.HTTPNotFound
	}
	return media_ress.FillFromEntity(mediaEty), nil
}
