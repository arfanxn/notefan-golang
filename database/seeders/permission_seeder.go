package seeders

import (
	"context"
	"database/sql"
	"time"

	perm_names "github.com/notefan-golang/enums/permission/names"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"

	"github.com/google/uuid"
)

type PermissionSeeder struct {
	db         *sql.DB
	repository *repositories.PermissionRepository
}

func NewPermissionSeeder(db *sql.DB) *PermissionSeeder {
	return &PermissionSeeder{
		db:         db,
		repository: repositories.NewPermissionRepository(db),
	}
}

func (seeder *PermissionSeeder) Run() {

	// Begin
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	permissionNames := []string{
		// Notification Module Permissions
		perm_names.NotificationView,
		perm_names.NotificationDelete,

		// Space Module Permissions
		perm_names.PageView,
		perm_names.PageCreate,
		perm_names.PageUpdate,
		perm_names.PageDelete,

		// Space Member (User) Module Permissions
		perm_names.SpaceMemberView,
		perm_names.SpaceMemberUpdateRole,
		perm_names.SpaceMemberInvite,
		perm_names.SpaceMemberRemove,

		// Page Module Permissions
		perm_names.PageView,
		perm_names.PageCreate,
		perm_names.PageUpdate,
		perm_names.PageDelete,

		// Page Content  Module Permissions
		perm_names.PageContentView,
		perm_names.PageContentCreate,
		perm_names.PageContentUpdate,
		perm_names.PageContentDelete,

		// Comment Module Permissions
		perm_names.CommentView,
		perm_names.CommentCreate,
		perm_names.CommentUpdate,
		perm_names.CommentDelete,

		// Comment Reaction Module Permissions
		perm_names.CommentReactionView,
		perm_names.CommentReactionCreate,
		perm_names.CommentReactionUpdate,
		perm_names.CommentReactionDelete,
	}
	var permissions []*entities.Permission

	for _, permissionName := range permissionNames {
		permission := entities.Permission{
			Id:   uuid.New(),
			Name: permissionName,
		}
		permissions = append(permissions, &permission)
	}

	_, err := seeder.repository.Insert(ctx, permissions...)
	errorh.LogPanic(err)
}
