package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"
	"time"
)

type NotificationUserSeeder struct {
	db               *sql.DB
	tableName        string
	repo             *repositories.NotificationUserRepo
	notificationRepo *repositories.NotificationRepo
	userRepo         *repositories.UserRepo
}

func NewNotificationUserSeeder(db *sql.DB) *NotificationUserSeeder {
	return &NotificationUserSeeder{
		db:               db,
		tableName:        "notification_user",
		repo:             repositories.NewNotificationUserRepo(db),
		notificationRepo: repositories.NewNotificationRepo(db),
		userRepo:         repositories.NewUserRepo(db),
	}
}

func (seeder *NotificationUserSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	notifications, err := seeder.notificationRepo.All(ctx)
	helper.PanicIfError(err)
	users, err := seeder.userRepo.All(ctx)
	helper.PanicIfError(err)

	notificationUsers := []entities.NotificationUser{}

	for _, notification := range notifications {
		notifier := users[rand.Intn(len(users))]
		notified := users[rand.Intn(len(users))]

		notificationUser := entities.NotificationUser{
			NotificationId: notification.Id,
			NotifierId:     notifier.Id,
			NotifiedId:     notified.Id,
		}
		notificationUsers = append(notificationUsers, notificationUser)

	}

	_, err = seeder.repo.Insert(ctx, notificationUsers...)
	helper.PanicIfError(err)
}
