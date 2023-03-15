package repositories

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/notefan-golang/containers/singletons"
	"github.com/notefan-golang/database/factories"
	token_types "github.com/notefan-golang/enums/token/types"
	"github.com/notefan-golang/helpers/numberh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
	"github.com/stretchr/testify/require"
)

func TestTokenRepository(t *testing.T) {
	require := require.New(t)

	app, appErr := singletons.GetApp()
	require.Nil(appErr)
	tokenRepository := repositories.NewTokenRepository(app.DB)
	userRepository := repositories.NewUserRepository(app.DB)
	ctx := context.Background()

	var (
		token entities.Token
		user  entities.User
	)

	t.Run("Create", func(t *testing.T) {
		user = factories.FakeUser()
		result, err := userRepository.Create(ctx, &user)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		token = factories.FakeToken()
		token.TokenableType = reflecth.GetTypeName(user)
		token.TokenableId = user.Id
		token.Type = token_types.ResetPassword
		token.UpdatedAt.Time = time.Time{} // Assign zero time
		token.UpdatedAt.Valid = false      // invalidate updated_at timestamp
		result, err = tokenRepository.Create(ctx, &token)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		require.NotEqual(uuid.Nil, token.Id)
		require.NotZero(token.CreatedAt)
		require.Zero(token.UpdatedAt.Time)
	})

	t.Run("Find", func(t *testing.T) {
		expectedToken := token

		actualToken, err := tokenRepository.Find(ctx, token.Id.String())
		require.Nil(err)
		require.Equal(expectedToken.Id.String(), actualToken.Id.String())
		require.NotZero(expectedToken.CreatedAt)
	})

	t.Run("FindByTokenableAndType", func(t *testing.T) {
		expectedToken := token

		actualToken, err := tokenRepository.FindByTokenableAndType(
			ctx,
			reflecth.GetTypeName(entities.User{}),
			user.Id.String(),
			token_types.ResetPassword,
		)
		require.Nil(err)
		require.Equal(expectedToken.Id.String(), actualToken.Id.String())

		require.NotEqual(uuid.Nil, token.Id)
		require.NotZero(expectedToken.CreatedAt)
	})

	t.Run("UpdateById", func(t *testing.T) {
		// Prepare value to be updated
		expectedToken := token
		expectedToken.Body = strconv.Itoa(numberh.Random(100000, 999999))
		// Do update
		result, err := tokenRepository.UpdateById(ctx, &expectedToken)
		require.Nil(err)
		require.NotZero(result.RowsAffected())
		require.NotZero(expectedToken.UpdatedAt.Time)
		token = expectedToken
	})

	t.Run("DeleteByIds", func(t *testing.T) {
		tokenOne := token
		ids := []string{tokenOne.Id.String()}
		result, err := tokenRepository.DeleteByIds(ctx, ids...)
		require.Nil(err)
		rowsAffected, err := result.RowsAffected()
		require.Nil(err)
		require.NotZero(rowsAffected)
	})
}
