package repositories

import (
	"context"
	"testing"
	"time"

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

	t.Run("Insert", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			expectedUsers := []entities.User{
				factories.FakeUser(),
				factories.FakeUser(),
			}
			for i := range expectedUsers {
				expectedUser := &expectedUsers[i]
				expectedUser.UpdatedAt.Time = time.Time{} // Assign zero time
				expectedUser.UpdatedAt.Valid = false      // invalidate updated_at timestamp
			}
			actualUsers, err := userRepository.Insert(ctx, expectedUsers[0], expectedUsers[1])
			require.Nil(err)
			for index, actualUser := range actualUsers {
				expectedUser := expectedUsers[index]
				require.NotEmpty(actualUser.Id.String())
				require.Equal(expectedUser.Name, actualUser.Name)
				require.Equal(expectedUser.Email, actualUser.Email)
				require.NotEmpty(actualUser.Password)
				require.NotZero(actualUser.CreatedAt)
				require.Zero(actualUser.UpdatedAt)
			}
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			expectedUser := factories.FakeUser()
			expectedUser.UpdatedAt.Time = time.Time{} // Assign zero time
			expectedUser.UpdatedAt.Valid = false      // invalidate updated_at timestamp
			actualUser, err := userRepository.Create(ctx, expectedUser)
			require.Nil(err)
			require.NotEmpty(actualUser.Id.String())
			require.Equal(expectedUser.Name, actualUser.Name)
			require.Equal(expectedUser.Email, actualUser.Email)
			require.NotEmpty(actualUser.Password)
			require.NotZero(actualUser.CreatedAt)
			require.Zero(actualUser.UpdatedAt)
		})
	})

	t.Run("All", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			users, err := userRepository.All(ctx)
			require.Nil(err)
			require.NotEmpty(users)
		})
	})

	t.Run("FindById", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			users, err := userRepository.All(ctx)
			require.Nil(err)
			actualUser := users[len(users)-1]
			expectedUser, err := userRepository.Find(ctx, actualUser.Id.String())
			require.Nil(err)
			require.Equal(expectedUser.Id.String(), actualUser.Id.String())
			require.Equal(expectedUser.Name, actualUser.Name)
			require.Equal(expectedUser.Email, actualUser.Email)
			require.NotEmpty(actualUser.Password)
			require.NotZero(actualUser.CreatedAt)
			require.Zero(actualUser.UpdatedAt)
		})
	})

	t.Run("FindByEmail", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			users, err := userRepository.All(ctx)
			require.Nil(err)
			actualUser := users[len(users)-1]
			expectedUser, err := userRepository.FindByEmail(ctx, actualUser.Email)
			require.Nil(err)
			require.Equal(expectedUser.Id.String(), actualUser.Id.String())
			require.Equal(expectedUser.Name, actualUser.Name)
			require.Equal(expectedUser.Email, actualUser.Email)
			require.NotEmpty(actualUser.Password)
			require.NotZero(actualUser.CreatedAt)
			require.Zero(actualUser.UpdatedAt)
		})
	})
}
