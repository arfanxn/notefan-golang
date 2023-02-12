package seeders

import (
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"runtime"

	"github.com/google/uuid"
)

func RoleSeeder(seeder DatabaseSeeder) {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// Begin
	roleNames := []string{
		"space owner",
		"space member",
	}

	tableName := "roles"
	valueArgs := []any{}

	for _, roleName := range roleNames {
		role := entities.Role{
			Id:   uuid.New(),
			Name: roleName,
		}
		valueArgs = append(valueArgs, role.Id.String(), role.Name)
	}

	query := helper.BuildBulkInsertQuery(tableName, len(roleNames), `id`, `name`)

	stmt, err := seeder.db.Prepare(query)
	helper.LogFatalIfError(err)

	_, err = stmt.Exec(valueArgs...)
	helper.LogFatalIfError(err)

}
