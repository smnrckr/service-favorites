package repositories_test

import (
	"testing"
	"user-favorites-service/internals/models"
	"user-favorites-service/internals/repositories"

	"github.com/stretchr/testify/assert"
)

func TestFavoritesListsRepository(t *testing.T) {
	favoritesListRepo := repositories.NewFavoritesListsRepository(testDb)
	t.Run("CreateFavoriteList", func(t *testing.T) {
		favoriteList := models.FavoriteList{
			Name:   "oyunlar",
			UserID: 2,
		}
		err := favoritesListRepo.CreateFavoriteList(&favoriteList)
		assert.NoError(t, err)
		assert.Equal(t, 1, favoriteList.Id)
	})

	t.Run("GetAllFavoriteLists", func(t *testing.T) {
		favoriteLists, err := favoritesListRepo.GetFavoriteListsByUserId(2)
		assert.NoError(t, err)
		assert.NotEmpty(t, favoriteLists)
		assert.Len(t, favoriteLists, 1)

	})

	t.Run("UpdateFavoriteList", func(t *testing.T) {

		updatedData := models.FavoriteList{
			Name: "ayakkabılar",
		}

		err := favoritesListRepo.UpdateFavoriteList(1, updatedData)
		assert.NoError(t, err)

		favoriteLists, err := favoritesListRepo.GetFavoriteListsByUserId(2)

		assert.NoError(t, err)
		assert.NotEmpty(t, favoriteLists)
		assert.Equal(t, "ayakkabılar", favoriteLists[0].Name)
	})

	t.Run("DeleteFavoriteList", func(t *testing.T) {

		err := favoritesListRepo.DeleteFavoriteListById(1, 2)
		assert.NoError(t, err)
		favoriteLists, err := favoritesListRepo.GetFavoriteListsByUserId(2)
		assert.NoError(t, err)
		assert.Empty(t, favoriteLists)

	})
	t.Run("DeleteFavoriteListInvalidId", func(t *testing.T) {

		err := favoritesListRepo.DeleteFavoriteListById(999, 999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrorNoRowsAffected)

	})
}
