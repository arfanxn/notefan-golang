package seeders

import (
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"

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
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// Begin
	roleNames := []string{
		"space owner",
		"space member",
	}

	valueArgs := []any{}

	for _, roleName := range roleNames {
		role := entities.Role{
			Id:   uuid.New(),
			Name: roleName,
		}
		valueArgs = append(valueArgs, role.Id.String(), role.Name)
	}

	query := helper.BuildBulkInsertQuery(seeder.tableName, len(roleNames), `id`, `name`)

	stmt, err := seeder.db.Prepare(query)
	helper.LogFatalIfError(err)

	_, err = stmt.Exec(valueArgs...)
	helper.LogFatalIfError(err)

}
