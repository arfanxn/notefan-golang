package repositories

import (
	"context"
	"testing"
	"time"

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
			result, err := userRepository.Insert(ctx, &expectedUsers[0], &expectedUsers[1])
			require.Nil(err)
			require.NotZero(result.RowsAffected())

			for _, expectedUser := range expectedUsers {
				require.NotEqual(uuid.Nil, expectedUser.Id)
				require.NotZero(expectedUser.CreatedAt)
			}
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			expectedUser := factories.FakeUser()
			expectedUser.UpdatedAt.Time = time.Time{} // Assign zero time
			expectedUser.UpdatedAt.Valid = false      // invalidate updated_at timestamp

			result, err := userRepository.Create(ctx, &expectedUser)
			require.Nil(err)
			require.NotZero(result.RowsAffected())

			require.NotEqual(uuid.Nil, expectedUser.Id)
			require.NotZero(expectedUser.CreatedAt)
		})
	})

	t.Run("UpdateById", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			actualUser := factories.FakeUser()
			actualUser.UpdatedAt.Time = time.Time{} // Assign zero time
			actualUser.UpdatedAt.Valid = false      // invalidate updated_at timestamp
			_, err := userRepository.Create(ctx, &actualUser)
			require.Nil(err)
			require.NotEqual(actualUser.Id, uuid.Nil)

			expectedUser := actualUser
			expectedUser.Name = "Arfan"
			expectedUser.Email = "arf@gm.com"

			require.Equal(expectedUser.Id, actualUser.Id)

			result, err := userRepository.UpdateById(ctx, &expectedUser)
			require.Nil(err)
			require.NotZero(result.RowsAffected())

			require.NotZero(expectedUser.UpdatedAt.Time)
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
