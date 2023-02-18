package services

import (
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

// TODO: Complete Media service
