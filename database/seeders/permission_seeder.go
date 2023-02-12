package seeders

import (
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"

	"github.com/google/uuid"
)

type PermissionSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.PermissionRepo
}

func NewPermissionSeeder(db *sql.DB) *PermissionSeeder {
	return &PermissionSeeder{
		db:        db,
		tableName: "permissions",
		repo:      repositories.NewPermissionRepo(db),
	}
}

func (seeder *PermissionSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// Begin
	permissionNames := []string{
		// Notification Module Permissions
		"view notification",
		"delete notification",

		// Space Member (User) Module Permissions
		"view member",
		"update member role",
		"delete member",

		// Space Module Permissions
		"view space",
		"create space",
		"update space",
		"delete space",

		// Page Module Permissions
		"view page",
		"create page",
		"update page",
		"delete page",

		// Page Content  Module Permissions
		"view page content",
		"create page content",
		"update page content",
		"delete page content",

		// Comment Module Permissions
		"view comment",
		"create comment",
		"update comment",
		"delete comment",

		// Comment Reaction Module Permissions
		"view comment reaction",
		"create comment reaction",
		"update comment reaction",
		"delete comment reaction",
	}

	valueArgs := []any{}

	for _, permissionName := range permissionNames {
		permission := entities.Permission{
			Id:   uuid.New(),
			Name: permissionName,
		}
		valueArgs = append(valueArgs, permission.Id.String(), permission.Name)
	}

	query := helper.BuildBulkInsertQuery(seeder.tableName, len(permissionNames), `id`, `name`)

	stmt, err := seeder.db.Prepare(query)
	helper.LogFatalIfError(err)

	_, err = stmt.Exec(valueArgs...)
	helper.LogFatalIfError(err)
}
