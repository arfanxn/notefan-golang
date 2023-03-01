package factories

import (
	"bytes"
	"database/sql"
	"path/filepath"
	"time"

	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/fileh"
	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/helpers/reflecth"
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
		// FileName // Will be filled in later
		// MimeType // Will be filled in later
		// Disk // Will be filled in later
		ConversionDisk: sql.NullString{
			Valid: false,
		},
		// Size // Will be filled in later
		CreatedAt: time.Now(),
		UpdatedAt: nullh.RandSqlNullTime(time.Now().AddDate(0, 0, 1)),

		File: bytes.NewBuffer([]byte{}),
	}
}

func FakeMediaForComment(comment entities.Comment) entities.Media {
	typ := reflecth.GetTypeName(comment)
	f, err := fileh.RandFromDir(mediaImagePlaceholderDirectoryPath)
	errorh.Panic(err)
	defer f.Close()

	filename := filepath.Base(f.Name())

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = comment.Id
	media.CollectionName = typ
	media.FileName = filename
	media.Size = fileh.GetSize(f)
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

	media.File.ReadFrom(f)

	return media
}

func FakeMediaForCommentReaction(cr entities.CommentReaction) entities.Media {
	typ := reflecth.GetTypeName(cr)
	f, err := fileh.RandFromDir(mediaImagePlaceholderDirectoryPath)
	errorh.Panic(err)
	defer f.Close()

	filename := filepath.Base(f.Name())

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = cr.Id
	media.CollectionName = typ
	media.FileName = filename
	media.Size = fileh.GetSize(f)
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

	media.File.ReadFrom(f)

	return media
}

func FakeMediaForNotification(notification entities.Notification) entities.Media {
	typ := reflecth.GetTypeName(notification)
	f, err := fileh.RandFromDir(mediaImagePlaceholderDirectoryPath)
	errorh.Panic(err)
	defer f.Close()

	filename := filepath.Base(f.Name())

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = notification.Id
	media.CollectionName = typ
	media.FileName = filename
	media.Size = fileh.GetSize(f)
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

	media.File.ReadFrom(f)

	return media
}

func FakeMediaForPage(notification entities.Page) entities.Media {
	typ := reflecth.GetTypeName(notification)
	f, err := fileh.RandFromDir(mediaImagePlaceholderDirectoryPath)
	errorh.Panic(err)
	defer f.Close()

	filename := filepath.Base(f.Name())

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = notification.Id
	media.CollectionName = typ
	media.FileName = filename
	media.Size = fileh.GetSize(f)
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

	media.File.ReadFrom(f)

	return media
}

func FakeMediaForPageContent(notification entities.PageContent) entities.Media {
	typ := reflecth.GetTypeName(notification)
	f, err := fileh.RandFromDir(mediaImagePlaceholderDirectoryPath)
	errorh.Panic(err)
	defer f.Close()

	filename := filepath.Base(f.Name())

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = notification.Id
	media.CollectionName = typ
	media.FileName = filename
	media.Size = fileh.GetSize(f)
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

	media.File.ReadFrom(f)

	return media
}

func FakeMediaForSpace(space entities.Space) entities.Media {
	typ := reflecth.GetTypeName(space)
	f, err := fileh.RandFromDir(mediaImagePlaceholderDirectoryPath)
	errorh.Panic(err)
	defer f.Close()

	filename := filepath.Base(f.Name())

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = space.Id
	media.CollectionName = typ
	media.FileName = filename
	media.Size = fileh.GetSize(f)
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

	media.File.ReadFrom(f)

	return media
}

func FakeMediaForUser(user entities.User) entities.Media {
	typ := reflecth.GetTypeName(user)
	f, err := fileh.RandFromDir(mediaImagePlaceholderDirectoryPath)
	errorh.Panic(err)
	defer f.Close()

	filename := filepath.Base(f.Name())

	media := FakeMedia()
	media.ModelType = typ
	media.ModelId = user.Id
	media.CollectionName = "avatar" // avatar represents user's profile picture
	media.FileName = filename
	media.Size = fileh.GetSize(f)
	media.MimeType = "image/jpeg"
	media.Disk = mediaDiskName

	media.File.ReadFrom(f)

	return media
}
