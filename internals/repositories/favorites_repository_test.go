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

	t.Run("GetAllFavoritesFromList", func(t *testing.T) {
		favorites, err := favoritesRepo.GetAllFavoritesFromList(1)
		assert.NoError(t, err)
		assert.NotEmpty(t, favorites)
		assert.Len(t, favorites, 2)

	})

	t.Run("DeletAllFavorites", func(t *testing.T) {

		err := favoritesRepo.DeleteAllFavoritesByListId(1, 1)
		assert.NoError(t, err)
		favorites, err := favoritesRepo.GetAllFavoritesFromList(1)
		assert.NoError(t, err)
		assert.Empty(t, favorites)

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
		favoritesFromList, err := favoritesRepo.GetAllFavoritesFromList(1)
		assert.NoError(t, err)
		assert.Len(t, favoritesFromList, 0)

	})

	t.Run("DeleteFavoritesInvalidId", func(t *testing.T) {
		err := favoritesRepo.DeleteFavoriteById(999, 999, 999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorNoRowsAffected)

	})

}
