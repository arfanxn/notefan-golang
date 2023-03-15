package user_ress

import (
	"time"

	"github.com/notefan-golang/models/entities"
	mediaRess "github.com/notefan-golang/models/responses/media_ress"
	"github.com/notefan-golang/models/responses/role_ress"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`

	// Role represents user role
	Role   role_ress.Role  `json:"role"`
	Avatar mediaRess.Media `json:"avatar,omitempty"`
}

func FillFromEntity(entity entities.User) User {
	res := User{
		Id:        entity.Id.String(),
		Name:      entity.Name,
		Email:     entity.Email,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: null.NewTime(entity.UpdatedAt.Time, entity.UpdatedAt.Valid),
		Role:      role_ress.FillFromEntity(entity.Role),
	}

	return res
}
