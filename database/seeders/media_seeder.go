package seeders

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/notefan-golang/config"
	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
	"github.com/notefan-golang/services"
)

type MediaSeeder struct {
	db                        *sql.DB
	service                   *services.MediaService
	repository                *repositories.MediaRepository
	commentRepository         *repositories.CommentRepository
	commentReactionRepository *repositories.CommentReactionRepository
	notificationRepository    *repositories.NotificationRepository
	pageRepository            *repositories.PageRepository
	pageContentRepository     *repositories.PageContentRepository
	spaceRepository           *repositories.SpaceRepository
	userRepository            *repositories.UserRepository

	waitGroup *sync.WaitGroup
}

func NewMediaSeeder(db *sql.DB) *MediaSeeder {
	return &MediaSeeder{
		db:                        db,
		service:                   services.NewMediaService(repositories.NewMediaRepository(db)),
		commentRepository:         repositories.NewCommentRepository(db),
		commentReactionRepository: repositories.NewCommentReactionRepository(db),
		notificationRepository:    repositories.NewNotificationRepository(db),
		pageRepository:            repositories.NewPageRepository(db),
		pageContentRepository:     repositories.NewPageContentRepository(db),
		spaceRepository:           repositories.NewSpaceRepository(db),
		userRepository:            repositories.NewUserRepository(db),

		waitGroup: new(sync.WaitGroup),
	}
}

func (seeder *MediaSeeder) Run() {
	seeder.cleanMediaStorages()

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute) // Give a minute timeout
	defer cancel()

	// seeder.waitGroup.Add(7)
	seeder.seedMediaForComments(ctx)
	seeder.seedMediaForCommentReactions(ctx)
	seeder.seedMediaForNotifications(ctx)
	seeder.seedMediaForPages(ctx)
	seeder.seedMediaForPageContents(ctx)
	seeder.seedMediaForSpaces(ctx)
	seeder.seedMediaForUsers(ctx)

	// seeder.waitGroup.Wait()
}

// cleanMediaStorage will clean the media storages
func (*MediaSeeder) cleanMediaStorages() {
	for _, disk := range config.FSDisks {
		os.RemoveAll(filepath.Join(disk.Root, "medias"))
	}
}

// seedMediaForComments seeds media for each comment
func (seeder *MediaSeeder) seedMediaForComments(ctx context.Context) {
	// defer seeder.waitGroup.Done()

	comments, err := seeder.commentRepository.All(ctx)
	errorh.LogPanic(err)

	medias := []entities.Media{}

	for _, comment := range comments {
		medias = append(medias, factories.FakeMediaForComment(comment))
	}

	_, err = seeder.repository.Insert(ctx, medias...)
	errorh.LogPanic(err)
}

// seedMediaForCommentReactions seeds media for each comment reaction
func (seeder *MediaSeeder) seedMediaForCommentReactions(ctx context.Context) {
	// defer seeder.waitGroup.Done()

	commentReactions, err := seeder.commentReactionRepository.All(ctx)
	errorh.LogPanic(err)

	medias := []entities.Media{}

	for _, cr := range commentReactions {
		medias = append(medias, factories.FakeMediaForCommentReaction(cr))
	}

	_, err = seeder.repository.Insert(ctx, medias...)
	errorh.LogPanic(err)
}

// seedMediaForNotifications seeds media for each notification
func (seeder *MediaSeeder) seedMediaForNotifications(ctx context.Context) {
	// defer seeder.waitGroup.Done()

	notifications, err := seeder.notificationRepository.All(ctx)
	errorh.LogPanic(err)

	medias := []entities.Media{}

	for _, notification := range notifications {
		medias = append(medias, factories.FakeMediaForNotification(notification))
	}

	_, err = seeder.repository.Insert(ctx, medias...)
	errorh.LogPanic(err)
}

// seedMediaForPages seeds media for each page
func (seeder *MediaSeeder) seedMediaForPages(ctx context.Context) {
	// defer seeder.waitGroup.Done()

	pages, err := seeder.pageRepository.All(ctx)
	errorh.LogPanic(err)

	medias := []entities.Media{}

	for _, page := range pages {
		medias = append(medias, factories.FakeMediaForPage(page))
	}

	_, err = seeder.repository.Insert(ctx, medias...)
	errorh.LogPanic(err)
}

// seedMediaForPageContents seeds media for each page content
func (seeder *MediaSeeder) seedMediaForPageContents(ctx context.Context) {
	// defer seeder.waitGroup.Done()

	pageContents, err := seeder.pageContentRepository.All(ctx)
	errorh.LogPanic(err)

	medias := []entities.Media{}

	for _, pageContent := range pageContents {
		medias = append(medias, factories.FakeMediaForPageContent(pageContent))
	}

	_, err = seeder.repository.Insert(ctx, medias...)
	errorh.LogPanic(err)
}

// seedMediaForSpaces seeds media for each space
func (seeder *MediaSeeder) seedMediaForSpaces(ctx context.Context) {
	// defer seeder.waitGroup.Done()

	spaces, err := seeder.spaceRepository.All(ctx)
	errorh.LogPanic(err)

	medias := []entities.Media{}

	for _, space := range spaces {
		medias = append(medias, factories.FakeMediaForSpace(space))
	}

	_, err = seeder.repository.Insert(ctx, medias...)
	errorh.LogPanic(err)
}

// seedMediaForUsers seeds media for each user
func (seeder *MediaSeeder) seedMediaForUsers(ctx context.Context) {
	// defer seeder.waitGroup.Done()

	users, err := seeder.userRepository.All(ctx)
	errorh.LogPanic(err)

	medias := []entities.Media{}

	for _, user := range users {
		medias = append(medias, factories.FakeMediaForUser(user))
	}

	_, err = seeder.repository.Insert(ctx, medias...)
	errorh.LogPanic(err)
}
