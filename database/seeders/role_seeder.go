package seeders

import (
	"context"
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"
	"time"

	"github.com/google/uuid"
)

type RoleSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.RoleRepo
}

func NewRoleSeeder(db *sql.DB) *RoleSeeder {
	return &RoleSeeder{
		db:        db,
		tableName: "roles",
		repo:      repositories.NewRoleRepo(db),
	}
}

func (seeder *RoleSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunningSeeder(pc)
	defer printFinishRunningSeeder(pc)

	// Begin
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	roleNames := []string{
		"space owner",
		"space member",
	}
	roles := []entities.Role{}

	for _, roleName := range roleNames {
		role := entities.Role{
			Id:   uuid.New(),
			Name: roleName,
		}
		roles = append(roles, role)
	}

	_, err := seeder.repo.Insert(ctx, roles...)
	helper.PanicIfError(err)

}
