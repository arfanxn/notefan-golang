package seeders

import (
	"context"
	"database/sql"
	"notefan-golang/config"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type MediaSeeder struct {
	db                  *sql.DB
	repo                *repositories.MediaRepo
	commentRepo         *repositories.CommentRepo
	commentReactionRepo *repositories.CommentReactionRepo
	notificationRepo    *repositories.NotificationRepo
	pageRepo            *repositories.PageRepo
	pageContentRepo     *repositories.PageContentRepo
	spaceRepo           *repositories.SpaceRepo
	userRepo            *repositories.UserRepo
}

func NewMediaSeeder(db *sql.DB) *MediaSeeder {
	return &MediaSeeder{
		db:                  db,
		repo:                repositories.NewMediaRepo(db),
		commentRepo:         repositories.NewCommentRepo(db),
		commentReactionRepo: repositories.NewCommentReactionRepo(db),
		notificationRepo:    repositories.NewNotificationRepo(db),
		pageRepo:            repositories.NewPageRepo(db),
		pageContentRepo:     repositories.NewPageContentRepo(db),
		spaceRepo:           repositories.NewSpaceRepo(db),
		userRepo:            repositories.NewUserRepo(db),
	}
}

func (seeder *MediaSeeder) Run() {
	seeder.cleanMediaStorages()

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute) // Give a minute timeout
	defer cancel()

	wg := new(sync.WaitGroup)

	wg.Add(7)
	go seeder.seedMediaForComments(ctx, wg)
	go seeder.seedMediaForCommentReactions(ctx, wg)
	go seeder.seedMediaForNotifications(ctx, wg)
	go seeder.seedMediaForPages(ctx, wg)
	go seeder.seedMediaForPageContents(ctx, wg)
	go seeder.seedMediaForSpaces(ctx, wg)
	go seeder.seedMediaForUsers(ctx, wg)

	wg.Wait()
}

// cleanMediaStorage will clean the media storages
func (*MediaSeeder) cleanMediaStorages() {
	for _, disk := range config.FSDisks {
		os.RemoveAll(filepath.Join(disk.Root, "medias"))
	}
}

// seedMediaForComments seeds media for each comment
func (seeder *MediaSeeder) seedMediaForComments(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	comments, err := seeder.commentRepo.All(ctx)
	helper.ErrorPanic(err)

	medias := []entities.Media{}

	for _, comment := range comments {
		medias = append(medias, factories.FakeMediaForComment(comment))
	}

	_, err = seeder.repo.Insert(ctx, medias...)
	helper.ErrorPanic(err)
}

// seedMediaForCommentReactions seeds media for each comment reaction
func (seeder *MediaSeeder) seedMediaForCommentReactions(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	commentReactions, err := seeder.commentReactionRepo.All(ctx)
	helper.ErrorPanic(err)

	medias := []entities.Media{}

	for _, cr := range commentReactions {
		medias = append(medias, factories.FakeMediaForCommentReaction(cr))
	}

	_, err = seeder.repo.Insert(ctx, medias...)
	helper.ErrorPanic(err)
}

// seedMediaForNotifications seeds media for each notification
func (seeder *MediaSeeder) seedMediaForNotifications(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	notifications, err := seeder.notificationRepo.All(ctx)
	helper.ErrorPanic(err)

	medias := []entities.Media{}

	for _, notification := range notifications {
		medias = append(medias, factories.FakeMediaForNotification(notification))
	}

	_, err = seeder.repo.Insert(ctx, medias...)
	helper.ErrorPanic(err)
}

// seedMediaForPages seeds media for each page
func (seeder *MediaSeeder) seedMediaForPages(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	pages, err := seeder.pageRepo.All(ctx)
	helper.ErrorPanic(err)

	medias := []entities.Media{}

	for _, page := range pages {
		medias = append(medias, factories.FakeMediaForPage(page))
	}

	_, err = seeder.repo.Insert(ctx, medias...)
	helper.ErrorPanic(err)
}

// seedMediaForPageContents seeds media for each page content
func (seeder *MediaSeeder) seedMediaForPageContents(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	pageContents, err := seeder.pageContentRepo.All(ctx)
	helper.ErrorPanic(err)

	medias := []entities.Media{}

	for _, pageContent := range pageContents {
		medias = append(medias, factories.FakeMediaForPageContent(pageContent))
	}

	_, err = seeder.repo.Insert(ctx, medias...)
	helper.ErrorPanic(err)
}

// seedMediaForSpaces seeds media for each space
func (seeder *MediaSeeder) seedMediaForSpaces(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	spaces, err := seeder.spaceRepo.All(ctx)
	helper.ErrorPanic(err)

	medias := []entities.Media{}

	for _, space := range spaces {
		medias = append(medias, factories.FakeMediaForSpace(space))
	}

	_, err = seeder.repo.Insert(ctx, medias...)
	helper.ErrorPanic(err)
}

// seedMediaForUsers seeds media for each user
func (seeder *MediaSeeder) seedMediaForUsers(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	users, err := seeder.userRepo.All(ctx)
	helper.ErrorPanic(err)

	medias := []entities.Media{}

	for _, user := range users {
		medias = append(medias, factories.FakeMediaForUser(user))
	}

	_, err = seeder.repo.Insert(ctx, medias...)
	helper.ErrorPanic(err)
}
