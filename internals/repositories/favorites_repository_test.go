package repositories_test

import (
	"testing"
	"user-favorites-service/internals/models"
	"user-favorites-service/internals/repositories"

	"github.com/stretchr/testify/assert"
)

func TestFavoritesRepository(t *testing.T) {
	favoritesRepo := repositories.NewFavoritesRepository(testDb)

	t.Run("AddProductToFavoriteList", func(t *testing.T) {

		favorites := models.Favorite{
			UserID:    1,
			ProductID: 2,
			ListID:    1,
		}
		err := favoritesRepo.AddProductToFavoriteList(&favorites)
		assert.NoError(t, err)
		assert.Equal(t, 2, favorites.Id)

	})

	t.Run("DeletAllFavorites", func(t *testing.T) {

		err := favoritesRepo.DeleteAllFavoritesByListId(1, 1)
		assert.NoError(t, err)

	})

	t.Run("DeleteFavorites", func(t *testing.T) {
		favorites := models.Favorite{
			UserID:    1,
			ProductID: 2,
			ListID:    1,
		}
		err := favoritesRepo.AddProductToFavoriteList(&favorites)
		assert.NoError(t, err)
		err = favoritesRepo.DeleteFavoriteById(1, 1, 2)
		assert.NoError(t, err)

	})

	t.Run("DeleteFavoritesInvalidId", func(t *testing.T) {
		err := favoritesRepo.DeleteFavoriteById(999, 999, 999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorNoRowsAffected)

	})

}
