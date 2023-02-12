package seeders

import (
	"context"
	"database/sql"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"
	"time"
)

type PermissionRoleSeeder struct {
	db             *sql.DB
	tableName      string
	repo           *repositories.PermissionRoleRepo
	permissionRepo *repositories.PermissionRepo
	roleRepo       *repositories.RoleRepo
}

func NewPermissionRoleSeeder(db *sql.DB) *PermissionRoleSeeder {
	return &PermissionRoleSeeder{
		db:             db,
		tableName:      "permission_role",
		repo:           repositories.NewPermissionRoleRepo(db),
		permissionRepo: repositories.NewPermissionRepo(db),
		roleRepo:       repositories.NewRoleRepo(db),
	}
}

func (seeder *PermissionRoleSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunningSeeder(pc)
	defer printFinishRunningSeeder(pc)

	/* ---- Begin ---- */
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	// Roles
	roleSpaceOwner, err := seeder.roleRepo.FindByName(ctx, "space owner")
	helper.PanicIfError(err)
	roleSpaceMember, err := seeder.roleRepo.FindByName(ctx, "space member")
	helper.PanicIfError(err)

	// Permissions
	roleSpaceOwnerPermissions, err := seeder.permissionRepo.All(ctx)
	helper.PanicIfError(err)
	roleSpaceMemberPermissions, err := seeder.getRoleSpaceMemberPermissions(ctx)
	helper.PanicIfError(err)

	// PermissionRoles
	permissionRoles := []entities.PermissionRole{}

	for _, permission := range roleSpaceOwnerPermissions {
		permissionRole := entities.PermissionRole{
			PermissionId: permission.Id,
			RoleId:       roleSpaceOwner.Id,
			CreatedAt:    time.Now(),
		}
		permissionRoles = append(permissionRoles, permissionRole)
	}
	for _, permission := range roleSpaceMemberPermissions {
		permissionRole := entities.PermissionRole{
			PermissionId: permission.Id,
			RoleId:       roleSpaceMember.Id,
			CreatedAt:    time.Now(),
		}
		permissionRoles = append(permissionRoles, permissionRole)
	}

	_, err = seeder.repo.Insert(ctx, permissionRoles...)
	helper.PanicIfError(err)
}

func (seeder PermissionRoleSeeder) getRoleSpaceMemberPermissions(ctx context.Context) (
	[]entities.Permission, error) {
	return seeder.permissionRepo.FindByNames(ctx,
		// Notification Module Permissions
		"view notification",

		// Space Member (User) Module Permissions
		"view member",

		// Space Module Permissions
		"view space",

		// Page Module Permissions
		"view page",

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
	)
}
