package responses

import (
	"time"

	"github.com/notefan-golang/models/entities"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`

	// TODO: User response with Avatar media
	// Avatar Media `json:"avatar,omitempty"`
}

func NewUserFromEntity(entity entities.User) User {
	return User{
		Id:        entity.Id.String(),
		Name:      entity.Name,
		Email:     entity.Email,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: null.NewTime(entity.UpdatedAt.Time, entity.UpdatedAt.Valid),
	}
}
