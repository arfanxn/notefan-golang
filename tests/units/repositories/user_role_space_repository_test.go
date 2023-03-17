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
	userRepository := repositories.NewUserRepository(app.DB)
	roleRepository := repositories.NewRoleRepository(app.DB)
	spaceRepository := repositories.NewSpaceRepository(app.DB)
	ursRepository := repositories.NewUserRoleSpaceRepository(app.DB)
	pageRepository := repositories.NewPageRepository(app.DB)
	pageContentRepository := repositories.NewPageContentRepository(app.DB)
	ctx := context.Background()

	var (
		user        entities.User
		role        entities.Role
		space       entities.Space
		urs         entities.UserRoleSpace
		page        entities.Page
		pageContent entities.PageContent
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

		// Create Page
		page = factories.FakePage()
		page.SpaceId = space.Id
		result, err = pageRepository.Create(ctx, &page)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		// Create Page Content
		pageContent = factories.FakePageContent()
		pageContent.PageId = page.Id
		result, err = pageContentRepository.Create(ctx, &pageContent)
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

	t.Run("FindByUserIdAndPageId", func(t *testing.T) {
		expectedUrs := urs
		actualUrs, err := ursRepository.FindByUserIdAndPageId(ctx,
			expectedUrs.UserId.String(), page.Id.String())
		require.Nil(err)
		require.Equal(expectedUrs.UserId.String(), actualUrs.UserId.String())
		require.Equal(expectedUrs.RoleId.String(), actualUrs.RoleId.String())
		require.Equal(expectedUrs.SpaceId.String(), actualUrs.SpaceId.String())
	})

	t.Run("FindByUserIdAndPageContentId", func(t *testing.T) {
		expectedUrs := urs
		actualUrs, err := ursRepository.FindByUserIdAndPageContentId(ctx,
			expectedUrs.UserId.String(), pageContent.Id.String())
		require.Nil(err)
		require.Equal(expectedUrs.UserId.String(), actualUrs.UserId.String())
		require.Equal(expectedUrs.RoleId.String(), actualUrs.RoleId.String())
		require.Equal(expectedUrs.SpaceId.String(), actualUrs.SpaceId.String())
	})

	t.Run("UpdateById", func(t *testing.T) {
		actualUrs := urs
		expectedUrs := actualUrs
		require.Equal(expectedUrs.Id, actualUrs.Id)
		roles, err := roleRepository.All(ctx)
		require.Nil(err)
		require.NotEmpty(roles)
		expectedUrs.RoleId = roles[0].Id

		result, err := ursRepository.UpdateById(ctx, &expectedUrs)
		require.Nil(err)
		require.NotZero(result.RowsAffected())
		require.NotZero(expectedUrs.UpdatedAt.Time)
		urs = expectedUrs
	})

	t.Run("DeleteByIds", func(t *testing.T) {
		result, err := ursRepository.DeleteByIds(ctx, urs.Id.String())
		require.Nil(err)
		require.NotZero(result.RowsAffected())
	})
}
