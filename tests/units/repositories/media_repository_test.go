package repositories

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/notefan-golang/containers/singletons"
	"github.com/notefan-golang/database/factories"
	"github.com/notefan-golang/helpers/fileh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/repositories"
	"github.com/stretchr/testify/require"
)

func TestMediaRepository(t *testing.T) {
	require := require.New(t)

	app, _ := singletons.GetApp()
	mediaRepository := repositories.NewMediaRepository(app.DB)
	ctx := context.Background()

	var media entities.Media

	t.Run("Create", func(t *testing.T) {
		media = factories.FakeMediaForPage(factories.FakePage())
		media.UpdatedAt.Time = time.Time{} // Assign zero time
		media.UpdatedAt.Valid = false      // invalidate updated_at timestamp

		result, err := mediaRepository.Create(ctx, &media)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		require.NotEqual(uuid.Nil, media.Id)
		require.NotZero(media.CreatedAt)
		require.NotEmpty(media.File.Bytes())

		// Assert open saved media file
		_, err = os.Open(media.GetFilePath())
		require.Nil(err)
	})

	t.Run("UpdateFileName", func(t *testing.T) {
		// Prepare value to be updated
		expectedMedia := media
		expectedMedia.FileName = "testfile" + "." + filepath.Ext(expectedMedia.FileName)
		// Do update
		result, err := mediaRepository.UpdateById(ctx, &expectedMedia)
		require.Nil(err)
		require.NotZero(result.RowsAffected())
		require.NotZero(expectedMedia.UpdatedAt.Time)
		require.NotEqual(expectedMedia.FileName, media.FileName)
		require.Equal(expectedMedia.File, media.File)

		// Assert Open media file after updating media.FileName
		_, err = os.Open(expectedMedia.GetFilePath())
		require.Nil(err)

		// Assert only one file at one media's directory
		mediaFileNames, err := fileh.FileNamesFromDir(filepath.Dir(expectedMedia.GetFilePath()))
		require.Equal(1, len(mediaFileNames))
	})

	t.Run("UpdateFile", func(t *testing.T) {
		// Get random file for media update
		file, err := fileh.RandFromDir("./public/placeholders/images")
		require.Nil(err)

		// Prepare value to be updated
		expectedMedia := media
		expectedMedia.FileName = filepath.Base(file.Name())
		// Reset for update
		expectedMedia.File.Reset()
		require.Empty(expectedMedia.File.Bytes())
		// Read new file
		_, err = expectedMedia.File.ReadFrom(file)
		require.Nil(err)
		// Ensure that media file is not empty
		require.NotEmpty(expectedMedia.File.Bytes())

		// Do update
		result, err := mediaRepository.UpdateById(ctx, &expectedMedia)
		require.Nil(err)
		require.NotZero(result.RowsAffected())

		// Assert entity updated
		require.NotZero(expectedMedia.UpdatedAt.Time)
		require.NotEmpty(expectedMedia.File.Bytes())

		// Assert Open media file after updating media.FileName
		file, err = os.Open(expectedMedia.GetFilePath())
		require.Nil(err)

		// Assert only one file at one media's directory
		mediaFileNames, err := fileh.FileNamesFromDir(filepath.Dir(expectedMedia.GetFilePath()))
		require.Equal(1, len(mediaFileNames))
	})

}
