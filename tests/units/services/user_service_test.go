package services

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestUserService(t *testing.T) {
	// TODO: Complete user service

	/*
		require := require.New(t)

		app, appErr := singletons.GetApp()
		require.Nil(appErr)
		userRepository := repositories.NewUserRepository(app.DB)
		mediaRepository := repositories.NewMediaRepository(app.DB)
		userService := services.NewUserService(userRepository, mediaRepository)
		ctx := context.Background()

		// create testing user
		user := factories.FakeUser()
		_, createErr := userRepository.Create(ctx, &user)
		require.Nil(createErr)

		t.Run("UpdateProfile", func(t *testing.T) {
			userReq := user_reqs.UpdateProfile{
				Id:     user.Id.String(),
				Name:   faker.Name(),
				Avatar: file_reqs.FillFromBytes(factories.FakeImageBuffer().Bytes()),
			}
			userRes, err := userService.UpdateProfile(ctx, userReq)
			require.Nil(err)
			require.Equal(userReq.Id, userRes.Id)
			require.Equal(userReq.Name, userRes.Name)
			require.NotZero(userRes.Email)
			require.NotZero(userRes.CreatedAt)
			require.NotZero(userRes.Avatar.Size)
		})
	*/
}
