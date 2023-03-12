package seeders

import (
	"context"
	"database/sql"
	"time"

	roleNames "github.com/notefan-golang/enums/role/names"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"

	"github.com/google/uuid"
)

type RoleSeeder struct {
	db         *sql.DB
	repository *repositories.RoleRepository
}

func NewRoleSeeder(db *sql.DB) *RoleSeeder {
	return &RoleSeeder{
		db:         db,
		repository: repositories.NewRoleRepository(db),
	}
}

func (seeder *RoleSeeder) Run() {

	// Begin
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	roleNames := []string{
		roleNames.SpaceOwner,
		roleNames.SpaceMember,
	}
	var roles []*entities.Role

	for _, roleName := range roleNames {
		role := entities.Role{
			Id:   uuid.New(),
			Name: roleName,
		}
		roles = append(roles, &role)
	}

	_, err := seeder.repository.Insert(ctx, roles...)
	errorh.LogPanic(err)

}
