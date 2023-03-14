package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/notefan-golang/containers/singletons"
	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/helpers/sliceh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
	"github.com/stretchr/testify/require"
)

func TestPermissionRepository(t *testing.T) {
	require := require.New(t)

	app, appErr := singletons.GetApp()
	require.Nil(appErr)
	permissionRepository := repositories.NewPermissionRepository(app.DB)
	roleRepository := repositories.NewRoleRepository(app.DB)
	permissionRoleRepository := repositories.NewPermissionRoleRepository(app.DB)
	ctx := context.Background()

	var permission entities.Permission

	t.Run("Create", func(t *testing.T) {
		permission = factories.FakePermission()

		result, err := permissionRepository.Create(ctx, &permission)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		require.NotEqual(uuid.Nil, permission.Id)
		require.NotZero(permission.Name)
	})

	t.Run("All", func(t *testing.T) {
		permissions, err := permissionRepository.All(ctx)
		require.Nil(err)
		require.NotEmpty(permissions)
	})

	t.Run("GetByRoleId", func(t *testing.T) {
		// Create the Role that own the Permission
		role := factories.FakeRole()
		_, err := roleRepository.Create(ctx, &role)
		require.Nil(err)

		// Asscociate the relation between Role and Permission
		permissionRole := entities.PermissionRole{
			PermissionId: permission.Id,
			RoleId:       role.Id,
			CreatedAt:    time.Now(),
		}
		_, err = permissionRoleRepository.Create(ctx, &permissionRole)
		require.Nil(err)

		actualPermission := permission
		expectedPermissions, err := permissionRepository.GetByRoleId(ctx, role.Id.String())
		require.Nil(err)

		expectedPermissions = sliceh.Filter(expectedPermissions,
			func(expectedPermission entities.Permission) bool {
				return expectedPermission.Id.String() == actualPermission.Id.String()
			},
		)
		require.NotEmpty(expectedPermissions)
		expectedPermission := expectedPermissions[0]
		permission = expectedPermission
	})

	t.Run("GetByNames", func(t *testing.T) {
		actualPermission := permission
		expectedPermissions, err := permissionRepository.GetByNames(ctx, actualPermission.Name)
		require.Nil(err)
		expectedPermission := expectedPermissions[0]
		require.Equal(expectedPermission.Id.String(), actualPermission.Id.String())
		require.Equal(expectedPermission.Name, actualPermission.Name)

		permission = expectedPermission
	})

	t.Run("FindByNameAndRoleId", func(t *testing.T) {
		// Create the Role that own the Permission
		role := factories.FakeRole()
		_, err := roleRepository.Create(ctx, &role)
		require.Nil(err)

		// Asscociate the relation between Role and Permission
		permissionRole := entities.PermissionRole{
			PermissionId: permission.Id,
			RoleId:       role.Id,
			CreatedAt:    time.Now(),
		}
		_, err = permissionRoleRepository.Create(ctx, &permissionRole)
		require.Nil(err)

		excpectedPermission := permission
		actualPermission, err := permissionRepository.
			FindByNameAndRoleId(ctx, excpectedPermission.Name, role.Id.String())
		require.Nil(err)

		require.Equal(excpectedPermission.Id.String(), actualPermission.Id.String())
	})
}
