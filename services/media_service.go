package services

import (
	"context"
	"database/sql"
	"path/filepath"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/requests/common_reqs"
	"github.com/notefan-golang/models/requests/media_reqs"
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

// Find finds media by the given request id
func (service *MediaService) Find(ctx context.Context, data common_reqs.UUID) (media_ress.Media, error) {
	mediaEty, err := service.repository.Find(ctx, data.Id)
	if err != nil { // err not nil return exception HTTPNotFound
		errorh.Log(err)
		return media_ress.Media{}, exceptions.HTTPNotFound
	}
	return media_ress.FillFromEntity(mediaEty), nil
}

// Update updates media by the given request id
func (service *MediaService) Update(ctx context.Context, data media_reqs.Update) (media_ress.Media, error) {
	mediaEty, err := service.repository.Find(ctx, data.Id)
	errorh.Panic(err) // panic if not found

	if data.Name != "" {
		mediaEty.Name = sql.NullString{String: data.Name, Valid: true}
	}
	if data.FileName != "" {
		fileName := stringh.FileNameWithoutExt(data.FileName) + filepath.Ext(mediaEty.FileName)
		mediaEty.FileName = fileName
	}
	if data.File != nil && data.File.IsProvided() {
		mediaEty.File = data.File
		mediaEty.FileName = stringh.FileNameWithoutExt(data.FileName) + data.File.Mime.Extension()
	}

	_, err = service.repository.UpdateById(ctx, &mediaEty)
	errorh.LogPanic(err)

	return media_ress.FillFromEntity(mediaEty), nil
}

// Delete deletes media by the given request id
func (service *MediaService) Delete(ctx context.Context, data common_reqs.UUID) error {
	mediaEty, err := service.repository.Find(ctx, data.Id)
	errorh.Panic(err) // panic if not found

	_, err = service.repository.DeleteByIds(ctx, mediaEty.Id.String())
	return err
}
