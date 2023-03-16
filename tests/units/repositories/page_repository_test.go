package repositories

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/notefan-golang/containers/singletons"
	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/helpers/sliceh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
	"github.com/stretchr/testify/require"
)

func TestPageRepository(t *testing.T) {
	require := require.New(t)

	app, appErr := singletons.GetApp()
	require.Nil(appErr)
	pageRepository := repositories.NewPageRepository(app.DB)
	spaceRepository := repositories.NewSpaceRepository(app.DB)
	ctx := context.Background()

	var (
		space entities.Space
		page  entities.Page
	)

	t.Run("Create", func(t *testing.T) {
		space = factories.FakeSpace()
		space.UpdatedAt.Time = time.Time{} // Assign zero time
		space.UpdatedAt.Valid = false      // invalidate updated_at timestamp
		result, err := spaceRepository.Create(ctx, &space)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		page = factories.FakePage()
		page.SpaceId = space.Id           // Associate with space id
		page.UpdatedAt.Time = time.Time{} // Assign zero time
		page.UpdatedAt.Valid = false      // invalidate updated_at timestamp
		result, err = pageRepository.Create(ctx, &page)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		require.NotEqual(uuid.Nil, page.Id)
		require.NotZero(page.CreatedAt)
	})

	t.Run("UpdateById", func(t *testing.T) {
		actualPage := page
		expectedPage := actualPage
		expectedPage.Title = faker.Word()
		expectedPage.Order = rand.Intn(10)

		require.Equal(expectedPage.Id, actualPage.Id)

		result, err := pageRepository.UpdateById(ctx, &expectedPage)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		require.NotZero(expectedPage.UpdatedAt.Time)

		page = expectedPage
	})

	t.Run("All", func(t *testing.T) { // test get all entities/rows from database table
		pages, err := pageRepository.All(ctx)
		require.Nil(err)
		require.NotEmpty(pages)
	})

	t.Run("GetBySpaceId", func(t *testing.T) { // test find entity/row from database table by id
		expectedPage := page
		actualPages, err := pageRepository.GetBySpaceId(ctx, expectedPage.SpaceId.String())
		require.Nil(err)
		require.NotEmpty(actualPages)
		require.NotNil(actualPages[0])
		require.Equal(actualPages[0].Id.String(), expectedPage.Id.String())
		require.NotZero(actualPages[0].CreatedAt)
		page = expectedPage
	})

	t.Run("FindById", func(t *testing.T) { // test find entity/row from database table by id
		actualPage := page
		expectedPage, err := pageRepository.Find(ctx, actualPage.Id.String())
		require.Nil(err)
		require.Equal(expectedPage.Id.String(), actualPage.Id.String())
		require.Equal(expectedPage.Title, actualPage.Title)
		require.Equal(expectedPage.Order, actualPage.Order)
		require.NotZero(actualPage.CreatedAt)
		require.NotZero(actualPage.UpdatedAt.Time)

		page = expectedPage
	})

	t.Run("DeleteByIds", func(t *testing.T) {
		pages := []entities.Page{page}
		ids := sliceh.Map(pages, func(page entities.Page) string { return page.Id.String() })
		result, err := pageRepository.DeleteByIds(ctx, ids...)
		require.Nil(err)
		require.NotZero(result.RowsAffected())
	})
}
