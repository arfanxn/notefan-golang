package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/notefan-golang/containers/singletons"
	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
	"github.com/stretchr/testify/require"
)

func TestUserRepository(t *testing.T) {
	require := require.New(t)

	app, _ := singletons.GetApp()
	userRepository := repositories.NewUserRepository(app.DB)
	ctx := context.Background()

	var user entities.User

	t.Run("Create", func(t *testing.T) {
		user = factories.FakeUser()
		user.UpdatedAt.Time = time.Time{} // Assign zero time
		user.UpdatedAt.Valid = false      // invalidate updated_at timestamp

		result, err := userRepository.Create(ctx, &user)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		require.NotEqual(uuid.Nil, user.Id)
		require.NotZero(user.CreatedAt)
	})

	t.Run("UpdateById", func(t *testing.T) {
		actualUser := user
		expectedUser := actualUser
		expectedUser.Name = faker.Name()
		expectedUser.Email = faker.Email()

		require.Equal(expectedUser.Id, actualUser.Id)

		result, err := userRepository.UpdateById(ctx, &expectedUser)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		require.NotZero(expectedUser.UpdatedAt.Time)

		user = expectedUser
	})

	t.Run("All", func(t *testing.T) {
		users, err := userRepository.All(ctx)
		require.Nil(err)
		require.NotEmpty(users)
	})

	t.Run("FindById", func(t *testing.T) {
		actualUser := user
		expectedUser, err := userRepository.Find(ctx, actualUser.Id.String())
		require.Nil(err)
		require.Equal(expectedUser.Id.String(), actualUser.Id.String())
		require.Equal(expectedUser.Name, actualUser.Name)
		require.Equal(expectedUser.Email, actualUser.Email)
		require.NotEmpty(actualUser.Password)
		require.NotZero(actualUser.CreatedAt)

		user = expectedUser
	})

	t.Run("FindByEmail", func(t *testing.T) {
		actualUser := user
		expectedUser, err := userRepository.FindByEmail(ctx, actualUser.Email)
		require.Nil(err)
		require.Equal(expectedUser.Id.String(), actualUser.Id.String())
		require.Equal(expectedUser.Name, actualUser.Name)
		require.Equal(expectedUser.Email, actualUser.Email)
		require.NotEmpty(actualUser.Password)
		require.NotZero(actualUser.CreatedAt)

		user = expectedUser
	})
}
