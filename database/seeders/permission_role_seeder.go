package seeders

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
)

type PermissionRoleSeeder struct {
	db                   *sql.DB
	repository           *repositories.PermissionRoleRepository
	permissionRepository *repositories.PermissionRepository
	roleRepository       *repositories.RoleRepository
}

func NewPermissionRoleSeeder(db *sql.DB) *PermissionRoleSeeder {
	return &PermissionRoleSeeder{
		db:                   db,
		repository:           repositories.NewPermissionRoleRepository(db),
		permissionRepository: repositories.NewPermissionRepository(db),
		roleRepository:       repositories.NewRoleRepository(db),
	}
}

func (seeder *PermissionRoleSeeder) Run() {

	/* ---- Begin ---- */
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	// Roles
	roleSpaceOwner, err := seeder.roleRepository.FindByName(ctx, "space owner")
	errorh.Panic(err)
	roleSpaceMember, err := seeder.roleRepository.FindByName(ctx, "space member")
	errorh.Panic(err)

	// Permissions
	roleSpaceOwnerPermissions, err := seeder.permissionRepository.All(ctx)
	errorh.Panic(err)
	roleSpaceMemberPermissions, err := seeder.getRoleSpaceMemberPermissions(ctx)
	errorh.Panic(err)

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

	_, err = seeder.repository.Insert(ctx, permissionRoles...)
	errorh.Panic(err)
}

func (seeder PermissionRoleSeeder) getRoleSpaceMemberPermissions(ctx context.Context) (
	[]entities.Permission, error) {
	return seeder.permissionRepository.FindByNames(ctx,
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
