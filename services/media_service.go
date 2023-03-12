package services

import (
	"context"
	"database/sql"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/notefan-golang/exceptions"
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
func (service *MediaService) Find(ctx context.Context, data common_reqs.UUID) (mediaRes media_ress.Media, err error) {
	mediaEty, err := service.repository.Find(ctx, data.Id)
	if err != nil {
		return
	}
	if mediaEty.Id == uuid.Nil { // if media not found return exception HTTPNotFound
		return mediaRes, exceptions.HTTPNotFound
	}
	mediaRes = media_ress.FillFromEntity(mediaEty)
	return mediaRes, nil
}

// Update updates media by the given request id
func (service *MediaService) Update(ctx context.Context, data media_reqs.Update) (mediaRes media_ress.Media, err error) {
	mediaEty, err := service.repository.Find(ctx, data.Id)
	if err != nil {
		return
	}
	if mediaEty.Id == uuid.Nil { // if media not found return exception HTTPNotFound
		return mediaRes, exceptions.HTTPNotFound
	}

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
	if err != nil {
		return
	}

	mediaRes = media_ress.FillFromEntity(mediaEty)
	return mediaRes, nil
}

// Delete deletes media by the given request id
func (service *MediaService) Delete(ctx context.Context, data common_reqs.UUID) (err error) {
	mediaEty, err := service.repository.Find(ctx, data.Id)
	if err != nil {
		return
	}
	if mediaEty.Id == uuid.Nil { // if media not found return exception HTTPNotFound
		err = exceptions.HTTPNotFound
		return
	}

	_, err = service.repository.DeleteByIds(ctx, mediaEty.Id.String())
	return err
}
