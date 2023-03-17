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

func TestPageContentRepository(t *testing.T) {
	require := require.New(t)

	app, appErr := singletons.GetApp()
	require.Nil(appErr)
	spaceRepository := repositories.NewSpaceRepository(app.DB)
	pageRepository := repositories.NewPageRepository(app.DB)
	pageContentRepository := repositories.NewPageContentRepository(app.DB)
	ctx := context.Background()

	var (
		space       entities.Space
		page        entities.Page
		pageContent entities.PageContent
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

		pageContent = factories.FakePageContent()
		pageContent.PageId = page.Id             // Associate with page id
		pageContent.UpdatedAt.Time = time.Time{} // Assign zero time
		pageContent.UpdatedAt.Valid = false      // invalidate updated_at timestamp
		result, err = pageContentRepository.Create(ctx, &pageContent)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		require.NotEqual(uuid.Nil, pageContent.Id)
		require.NotZero(pageContent.CreatedAt)
	})

	t.Run("UpdateById", func(t *testing.T) {
		actualPageContent := pageContent
		expectedPageContent := actualPageContent
		expectedPageContent.Body = faker.Sentence()
		expectedPageContent.Order = rand.Intn(10)

		require.Equal(expectedPageContent.Id, actualPageContent.Id)

		result, err := pageContentRepository.UpdateById(ctx, &expectedPageContent)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		require.NotZero(expectedPageContent.UpdatedAt.Time)

		pageContent = expectedPageContent
	})

	t.Run("All", func(t *testing.T) { // test get all entities/rows from database table
		pageContents, err := pageContentRepository.All(ctx)
		require.Nil(err)
		require.NotEmpty(pageContents)
	})

	t.Run("GetBySpaceId", func(t *testing.T) { // test find entity/row from database table by id
		expectedPageContent := pageContent
		actualPageContents, err := pageContentRepository.GetByPageId(ctx, expectedPageContent.PageId.String())
		require.Nil(err)
		require.NotEmpty(actualPageContents)
		require.NotNil(actualPageContents[0])
		require.Equal(actualPageContents[0].Id.String(), expectedPageContent.Id.String())
		require.NotZero(actualPageContents[0].CreatedAt)
		pageContent = expectedPageContent
	})

	t.Run("FindById", func(t *testing.T) { // test find entity/row from database table by id
		actualPageContent := pageContent
		expectedPageContent, err := pageContentRepository.Find(ctx, actualPageContent.Id.String())
		require.Nil(err)
		require.Equal(expectedPageContent.Id.String(), actualPageContent.Id.String())
		require.Equal(expectedPageContent.Body, actualPageContent.Body)
		require.Equal(expectedPageContent.Order, actualPageContent.Order)
		require.NotZero(actualPageContent.CreatedAt)
		require.NotZero(actualPageContent.UpdatedAt.Time)

		pageContent = expectedPageContent
	})

	t.Run("DeleteByIds", func(t *testing.T) {
		pageContents := []entities.PageContent{pageContent}
		ids := sliceh.Map(pageContents, func(pageContent entities.PageContent) string { return pageContent.Id.String() })
		result, err := pageContentRepository.DeleteByIds(ctx, ids...)
		require.Nil(err)
		require.NotZero(result.RowsAffected())
	})
}
