package factories

import (
	"database/sql"
	"time"

	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

// var file // TODO: Open a universal file for make fake data
var mediaDisk = "public"
var mediaImagePlaceholderDirectoryPath = "./public/placeholders/images"

func FakeMedia() entities.Media {
	return entities.Media{
		Id: uuid.New(),
		// ModelType: , // Will be filled in later
		// ModelId: , // Will be filled in later
		CollectionName: "default",
		// Name: , // Will be filled in later
		// FileName // TODO
		// MimeType // TODO
		// Disk // TODO
		ConversionDisk: sql.NullString{
			Valid: false,
		},
		// Size // TODO
		CreatedAt: time.Now(),
		UpdatedAt: helper.DBRandNullOrTime(time.Now().AddDate(0, 0, 1)),
	}
}

func FakeMediaForComment(comment entities.Comment) entities.Media {
	typ := helper.ReflectGetTypeName(comment)
	f, err := helper.FileRandFromDir(mediaImagePlaceholderDirectoryPath)
	helper.ErrorPanic(err)

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = comment.Id
	media.CollectionName = typ
	media.FileName = f.Name()
	media.Size = helper.FileSize(f)
	mimeType, err := helper.FileContentType(f) // get the mime type
	helper.ErrorPanic(err)
	media.MimeType = mimeType
	media.Disk = mediaDisk

	media.File = f

	return media
}

func FakeMediaForCommentReaction(cr entities.CommentReaction) entities.Media {
	typ := helper.ReflectGetTypeName(cr)
	f, err := helper.FileRandFromDir(mediaImagePlaceholderDirectoryPath)
	helper.ErrorPanic(err)

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = cr.Id
	media.CollectionName = typ
	media.FileName = f.Name()
	media.Size = helper.FileSize(f)
	mimeType, err := helper.FileContentType(f) // get the mime type
	helper.ErrorPanic(err)
	media.MimeType = mimeType
	media.Disk = mediaDisk

	media.File = f

	return media
}

func FakeMediaForNotification(notification entities.Notification) entities.Media {
	typ := helper.ReflectGetTypeName(notification)
	f, err := helper.FileRandFromDir(mediaImagePlaceholderDirectoryPath)
	helper.ErrorPanic(err)

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = notification.Id
	media.CollectionName = typ
	media.FileName = f.Name()
	media.Size = helper.FileSize(f)
	mimeType, err := helper.FileContentType(f) // get the mime type
	helper.ErrorPanic(err)
	media.MimeType = mimeType
	media.Disk = mediaDisk

	media.File = f

	return media
}

func FakeMediaForPage(notification entities.Page) entities.Media {
	typ := helper.ReflectGetTypeName(notification)
	f, err := helper.FileRandFromDir(mediaImagePlaceholderDirectoryPath)
	helper.ErrorPanic(err)

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = notification.Id
	media.CollectionName = typ
	media.FileName = f.Name()
	media.Size = helper.FileSize(f)
	mimeType, err := helper.FileContentType(f) // get the mime type
	helper.ErrorPanic(err)
	media.MimeType = mimeType
	media.Disk = mediaDisk

	media.File = f

	return media
}

func FakeMediaForPageContent(notification entities.PageContent) entities.Media {
	typ := helper.ReflectGetTypeName(notification)
	f, err := helper.FileRandFromDir(mediaImagePlaceholderDirectoryPath)
	helper.ErrorPanic(err)

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = notification.Id
	media.CollectionName = typ
	media.FileName = f.Name()
	media.Size = helper.FileSize(f)
	mimeType, err := helper.FileContentType(f) // get the mime type
	helper.ErrorPanic(err)
	media.MimeType = mimeType
	media.Disk = mediaDisk

	media.File = f

	return media
}

func FakeMediaForSpace(notification entities.Space) entities.Media {
	typ := helper.ReflectGetTypeName(notification)
	f, err := helper.FileRandFromDir(mediaImagePlaceholderDirectoryPath)
	helper.ErrorPanic(err)

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = notification.Id
	media.CollectionName = typ
	media.FileName = f.Name()
	media.Size = helper.FileSize(f)
	mimeType, err := helper.FileContentType(f) // get the mime type
	helper.ErrorPanic(err)
	media.MimeType = mimeType
	media.Disk = mediaDisk

	media.File = f

	return media
}

func FakeMediaForUser(notification entities.User) entities.Media {
	typ := helper.ReflectGetTypeName(notification)
	f, err := helper.FileRandFromDir(mediaImagePlaceholderDirectoryPath)
	helper.ErrorPanic(err)

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = notification.Id
	media.CollectionName = "avatar" // avatar represents user's profile picture
	media.FileName = f.Name()
	media.Size = helper.FileSize(f)
	mimeType, err := helper.FileContentType(f) // get the mime type
	helper.ErrorPanic(err)
	media.MimeType = mimeType
	media.Disk = mediaDisk

	media.File = f

	return media
}
