package services

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/notefan-golang/config"
	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"
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

/* // TODO: Refactor inside of this function and implement goroutines  */
// Insert insert medias into database and stores media's file to selected disk (filesystem disk)
func (service *MediaService) Insert(ctx context.Context, medias ...entities.Media) ([]entities.Media, error) {
	storedFilePaths := []string{}

	// store media to storage
	for _, media := range medias {
		if media.CollectionName == "" {
			media.CollectionName = "default"
		}
		if media.MimeType == "" {
			mimeType, err := helper.FileContentType(media.File)
			if err != nil {
				helper.ErrorLog(err)
				for _, storedFilePath := range storedFilePaths { // once it fails, delete all previuosly stored files
					os.RemoveAll(storedFilePath)
				}
				return medias, exceptions.FileInvalidType
			}
			media.MimeType = mimeType
		}

		// If file exists do write operation
		if helper.FileSize(media.File) > 0 {
			fileSrc := media.File
			defer fileSrc.Close()

			root := config.FSDisks[media.Disk].Root
			path := filepath.Join(root, "medias", media.Id.String(), filepath.Base(fileSrc.Name()))

			err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
			fileDst, err := os.Create(path)
			defer fileDst.Close()
			if err != nil {
				helper.ErrorLog(err)
				for _, storedFilePath := range storedFilePaths { // once it fails, delete all previuosly stored files
					os.RemoveAll(storedFilePath)
				}
				return medias, err
			}

			_, err = io.Copy(fileDst, fileSrc)
			if err != nil {
				helper.ErrorLog(err)
				for _, storedFilePath := range storedFilePaths { // once it fails, delete all previuosly stored files
					os.RemoveAll(storedFilePath)
				}
				return medias, err
			}

			storedFilePaths = append(storedFilePaths, path)
		} else { // otherwise returns error
			err := exceptions.FileNotProvided
			helper.ErrorLog(err)
			for _, storedFilePath := range storedFilePaths { // once it fails, delete all previuosly stored files
				os.RemoveAll(storedFilePath)
			}
			return medias, err
		}
	}

	// insert media to database
	medias, err := service.repository.Insert(ctx, medias...)
	if err != nil {
		for _, storedFilePath := range storedFilePaths { // once it fails, delete all previuosly stored files
			os.RemoveAll(storedFilePath)
		}
		return medias, err
	}

	return medias, nil
}
