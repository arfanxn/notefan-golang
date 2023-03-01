package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
)

type NotificationUserSeeder struct {
	db                     *sql.DB
	repository             *repositories.NotificationUserRepository
	notificationRepository *repositories.NotificationRepository
	userRepository         *repositories.UserRepository
}

func NewNotificationUserSeeder(db *sql.DB) *NotificationUserSeeder {
	return &NotificationUserSeeder{
		db:                     db,
		repository:             repositories.NewNotificationUserRepository(db),
		notificationRepository: repositories.NewNotificationRepository(db),
		userRepository:         repositories.NewUserRepository(db),
	}
}

func (seeder *NotificationUserSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	notifications, err := seeder.notificationRepository.All(ctx)
	errorh.Panic(err)

	users, err := seeder.userRepository.All(ctx)
	errorh.Panic(err)

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

	_, err = seeder.repository.Insert(ctx, notificationUsers...)
	errorh.Panic(err)
}
