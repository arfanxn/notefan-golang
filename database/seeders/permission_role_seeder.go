package seeders

import (
	"context"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"
	"time"
)

func PermissionRoleSeeder(seeder DatabaseSeeder) {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	/* ---- Begin ---- */
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	roleRepo := repositories.NewRoleRepo(seeder.db)
	permissionRepo := repositories.NewPermissionRepo(seeder.db)

	roleSpaceOwner, err := roleRepo.FindByName(ctx, "space owner")
	helper.PanicIfError(err)
	roleSpaceMember, err := roleRepo.FindByName(ctx, "space member")
	helper.PanicIfError(err)

	roleSpaceOwnerPermissions, err := permissionRepo.All(ctx)
	helper.PanicIfError(err)

	roleSpaceMemberPermissions, err := permissionRepo.GetByNames(ctx,
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
	helper.PanicIfError(err)

	tableName := "permission_role"
	valueArgs := []any{}
	for _, permission := range roleSpaceOwnerPermissions {
		permissionRole := entities.PermissionRole{
			PermissionId: permission.Id,
			RoleId:       roleSpaceOwner.Id,
			CreatedAt:    time.Now(),
		}
		valueArgs = append(
			valueArgs,
			permissionRole.PermissionId.String(), permissionRole.RoleId.String(), permissionRole.CreatedAt)
	}
	for _, permission := range roleSpaceMemberPermissions {
		permissionRole := entities.PermissionRole{
			PermissionId: permission.Id,
			RoleId:       roleSpaceMember.Id,
			CreatedAt:    time.Now(),
		}
		valueArgs = append(
			valueArgs,
			permissionRole.PermissionId.String(), permissionRole.RoleId.String(), permissionRole.CreatedAt)
	}

	query := helper.BuildBulkInsertQuery(
		tableName,
		len(roleSpaceOwnerPermissions)+len(roleSpaceMemberPermissions),
		`permission_id`, `role_id`, `created_at`)

	stmt, err := seeder.db.Prepare(query)
	helper.PanicIfError(err)

	_, err = stmt.Exec(valueArgs...)
	helper.PanicIfError(err)

	cancel()
}
