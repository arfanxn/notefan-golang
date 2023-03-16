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

func TestRoleRepository(t *testing.T) {
	require := require.New(t)

	app, appErr := singletons.GetApp()
	require.Nil(appErr)
	roleRepository := repositories.NewRoleRepository(app.DB)
	ctx := context.Background()

	var (
		role entities.Role
	)

	t.Run("Create", func(t *testing.T) {
		// Create Role
		role = factories.FakeRole()
		result, err := roleRepository.Create(ctx, &role)
		require.Nil(err)
		require.NotZero(result.RowsAffected())
	})

	t.Run("All", func(t *testing.T) {
		expectedRoles := []entities.Role{role}
		actualRoles, err := roleRepository.All(ctx)
		require.Nil(err)
		require.Equal(
			expectedRoles[0].Id.String(),
			actualRoles[0].Id.String(),
		)
	})

	t.Run("FindByName", func(t *testing.T) {
		expectedRole := role

		actualRole, err := roleRepository.FindByName(ctx, expectedRole.Name)
		require.Nil(err)
		require.Equal(expectedRole.Id.String(), actualRole.Id.String())
		role = expectedRole
	})
}
