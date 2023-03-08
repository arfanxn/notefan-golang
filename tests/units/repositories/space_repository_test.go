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

func TestSpaceRepository(t *testing.T) {
	require := require.New(t)

	app, appErr := singletons.GetApp()
	require.Nil(appErr)
	spaceRepository := repositories.NewSpaceRepository(app.DB)
	ctx := context.Background()

	var space entities.Space

	t.Run("Create", func(t *testing.T) { // test create an entity/row
		space = factories.FakeSpace()
		space.UpdatedAt.Time = time.Time{} // Assign zero time
		space.UpdatedAt.Valid = false      // invalidate updated_at timestamp

		result, err := spaceRepository.Create(ctx, &space)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		require.NotEqual(uuid.Nil, space.Id)
		require.NotZero(space.CreatedAt)
	})

	t.Run("All", func(t *testing.T) { // test get all entities/rows from database table
		spaces, err := spaceRepository.All(ctx)
		require.Nil(err)
		require.NotEmpty(spaces)
	})

	t.Run("FindById", func(t *testing.T) { // test find entity/row from database table by id
		actualSpace := space
		expectedSpace, err := spaceRepository.Find(ctx, actualSpace.Id.String())
		require.Nil(err)
		require.Equal(expectedSpace.Id.String(), actualSpace.Id.String())
		require.Equal(expectedSpace.Name, actualSpace.Name)
		require.Equal(expectedSpace.Description, actualSpace.Description)
		require.Equal(expectedSpace.Domain, actualSpace.Domain)
		require.NotZero(actualSpace.CreatedAt)

		space = expectedSpace
	})

	t.Run("UpdateById", func(t *testing.T) { // test update entity/row from database table by id
		actualSpace := space
		expectedSpace := actualSpace
		expectedSpace.Name = faker.Word()
		expectedSpace.Description = faker.Sentence()

		require.Equal(expectedSpace.Id, actualSpace.Id)

		result, err := spaceRepository.UpdateById(ctx, &expectedSpace)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		require.NotZero(expectedSpace.UpdatedAt.Time)

		space = expectedSpace
	})

	t.Run("DeleteByIds", func(t *testing.T) {
		spaceOne := space

		spaceTwo := factories.FakeSpace()
		result, err := spaceRepository.Create(ctx, &spaceTwo)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		ids := []string{spaceOne.Id.String(), spaceTwo.Id.String()}
		result, err = spaceRepository.DeleteByIds(ctx, ids...)
		require.Nil(err)
		rowsAffected, err := result.RowsAffected()
		require.Nil(err)
		require.Equal(int64(2), rowsAffected)
	})
}
