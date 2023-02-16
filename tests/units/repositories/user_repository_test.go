package repositories

import (
	"context"
	"testing"

	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/repositories"
	"github.com/notefan-golang/tests"
	"github.com/stretchr/testify/require"
)

func TestUserRepository(t *testing.T) {
	require := require.New(t)

	app := tests.GetApp()
	userRepository := repositories.NewUserRepository(app.DB)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			expectedUser := factories.FakeUser()
			actualUser, err := userRepository.Create(ctx, expectedUser)
			require.Nil(err)
			require.Equal(expectedUser.Name, actualUser.Name)
			require.Equal(expectedUser.Email, actualUser.Email)
		})
	})
}
