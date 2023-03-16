package repositories

import (
	"context"
	"testing"

	"github.com/notefan-golang/containers/singletons"
	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
	"github.com/stretchr/testify/require"
)

func TestUserRoleSpaceRepository(t *testing.T) {
	require := require.New(t)

	app, appErr := singletons.GetApp()
	require.Nil(appErr)
	spaceRepository := repositories.NewSpaceRepository(app.DB)
	userRepository := repositories.NewUserRepository(app.DB)
	ursRepository := repositories.NewUserRoleSpaceRepository(app.DB)
	roleRepository := repositories.NewRoleRepository(app.DB)
	ctx := context.Background()

	var (
		urs   entities.UserRoleSpace
		user  entities.User
		role  entities.Role
		space entities.Space
	)

	t.Run("Create", func(t *testing.T) {
		// Create User
		user = factories.FakeUser()
		result, err := userRepository.Create(ctx, &user)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		// Create Role
		role = factories.FakeRole()
		result, err = roleRepository.Create(ctx, &role)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		// Create Space
		space = factories.FakeSpace()
		result, err = spaceRepository.Create(ctx, &space)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		// Create UserRoleSpace
		urs = entities.UserRoleSpace{
			UserId:  user.Id,
			RoleId:  role.Id,
			SpaceId: space.Id,
		}
		result, err = ursRepository.Create(ctx, &urs)
		require.Nil(err)
		require.NotZero(result.RowsAffected())
	})

	t.Run("FindByUserIdAndSpaceId", func(t *testing.T) {
		expectedUrs := urs
		actualUrs, err := ursRepository.FindByUserIdAndSpaceId(ctx,
			expectedUrs.UserId.String(), expectedUrs.SpaceId.String())
		require.Nil(err)
		require.Equal(expectedUrs.UserId.String(), actualUrs.UserId.String())
		require.Equal(expectedUrs.RoleId.String(), actualUrs.RoleId.String())
		require.Equal(expectedUrs.SpaceId.String(), actualUrs.SpaceId.String())
	})

	t.Run("DeleteByIds", func(t *testing.T) {
		result, err := ursRepository.DeleteByIds(ctx, urs.Id.String())
		require.Nil(err)
		require.NotZero(result.RowsAffected())
	})
}
