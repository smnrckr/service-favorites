package repositories_test

import (
	"testing"
	"user-favorites-service/internals/models"
	"user-favorites-service/internals/repositories"

	"github.com/stretchr/testify/assert"
)

func TestFavoritesRepository(t *testing.T) {
	favoritesRepo := repositories.NewFavoritesRepository(testDb)
	favoritesListRepo := repositories.NewFavoritesListsRepository(testDb)
	t.Run("AddProductToFavoriteList", func(t *testing.T) {

		favoriteList := models.FavoriteList{
			Name:   "oyunlar",
			UserID: 2,
		}
		err := favoritesListRepo.CreateFavoriteList(&favoriteList)
		assert.NoError(t, err)
		assert.Equal(t, 1, favoriteList.Id)

		favorites := models.Favorite{
			UserID:    1,
			ProductID: 1,
			ListID:    1,
		}
		err = favoritesRepo.AddProductToFavoriteList(&favorites)
		assert.NoError(t, err)
		assert.Equal(t, 1, favorites.Id)

	})

	t.Run("GetAllFavoritesFromList", func(t *testing.T) {
		favorites, err := favoritesRepo.GetAllFavoritesFromList(1)
		assert.NoError(t, err)
		assert.NotEmpty(t, favorites)
		assert.Len(t, favorites, 1)

	})

	t.Run("DeleteFavorites", func(t *testing.T) {

		err := favoritesRepo.DeleteFavoriteById(1, 1, 1)
		assert.NoError(t, err)
		favoriteLists, err := favoritesListRepo.GetFavoriteListsByUserId(1)
		assert.NoError(t, err)
		assert.Empty(t, favoriteLists)

	})
	t.Run("DeleteFavoritesInvalidId", func(t *testing.T) {

		err := favoritesRepo.DeleteFavoriteById(999, 999, 999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorNoRowsAffected)

	})

}
