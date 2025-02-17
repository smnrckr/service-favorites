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
		assert.Equal(t, 2, favoriteList.Id)
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

		updatedList, err := favoritesListRepo.UpdateFavoriteList(2, updatedData)
		assert.NoError(t, err)

		favoriteLists, err := favoritesListRepo.GetFavoriteListsByUserId(2)

		assert.NoError(t, err)
		assert.NotEmpty(t, favoriteLists)
		assert.Equal(t, updatedList.Name, favoriteLists[0].Name)
	})

	t.Run("DeleteFavoriteList", func(t *testing.T) {

		err := favoritesListRepo.DeleteFavoriteListById(2, 2)
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
