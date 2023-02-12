package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"
	"strings"
	"time"
)

type NotificationSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.NotificationRepo
	userRepo  *repositories.UserRepo
	spaceRepo *repositories.SpaceRepo
}

func NewNotificationSeeder(db *sql.DB) *NotificationSeeder {
	return &NotificationSeeder{
		db:        db,
		tableName: "notifications",
		repo:      repositories.NewNotificationRepo(db),
		userRepo:  repositories.NewUserRepo(db),
		spaceRepo: repositories.NewSpaceRepo(db),
	}
}

func (seeder *NotificationSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunningSeeder(pc)
	defer printFinishRunningSeeder(pc)

	// Begin
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users, err := seeder.userRepo.All(ctx)
	helper.PanicIfError(err)

	spaces, err := seeder.spaceRepo.All(ctx)
	helper.PanicIfError(err)

	notifications := []entities.Notification{}

	for i := 0; i < len(users); i++ {
		for j := 0; j < 5; j++ { // each user has 5 notifications
			space := spaces[rand.Intn(len(spaces))]

			notification := factories.NewNotification()
			notification.ObjectType = strings.ToUpper(helper.GetTypeName(space))
			notification.ObjectId = space.Id

			notifications = append(notifications, notification)
		}
	}

	_, err = seeder.repo.Insert(ctx, notifications...)
	helper.PanicIfError(err)
}
