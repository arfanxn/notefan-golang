package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"strings"
	"time"

	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
)

type NotificationSeeder struct {
	db              *sql.DB
	repository      *repositories.NotificationRepository
	userRepository  *repositories.UserRepository
	spaceRepository *repositories.SpaceRepository
}

func NewNotificationSeeder(db *sql.DB) *NotificationSeeder {
	return &NotificationSeeder{
		db:              db,
		repository:      repositories.NewNotificationRepository(db),
		userRepository:  repositories.NewUserRepository(db),
		spaceRepository: repositories.NewSpaceRepository(db),
	}
}

func (seeder *NotificationSeeder) Run() {

	// Begin
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	users, err := seeder.userRepository.All(ctx)
	errorh.Panic(err)

	spaces, err := seeder.spaceRepository.All(ctx)
	errorh.Panic(err)

	notifications := []entities.Notification{}

	for i := 0; i < len(users); i++ {
		for j := 0; j < 5; j++ { // each user has 5 notifications
			space := spaces[rand.Intn(len(spaces))]

			notification := factories.FakeNotification()
			notification.ObjectType = strings.ToUpper(reflecth.GetTypeName(space))
			notification.ObjectId = space.Id

			notifications = append(notifications, notification)
		}
	}

	_, err = seeder.repository.Insert(ctx, notifications...)
	errorh.Panic(err)
}
