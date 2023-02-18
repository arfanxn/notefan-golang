package factories

import (
	"database/sql"
	"time"

	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

var mediaDiskName = "public"
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
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

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
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

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
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

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
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

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
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

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
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

	media.File = f

	return media
}

func FakeMediaForUser(user entities.User) entities.Media {

	typ := helper.ReflectGetTypeName(user)
	f, err := helper.FileRandFromDir(mediaImagePlaceholderDirectoryPath)
	helper.ErrorPanic(err)

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = user.Id
	media.CollectionName = "avatar" // avatar represents user's profile picture
	media.FileName = f.Name()
	media.Size = helper.FileSize(f)
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

	media.File = f

	return media
}
