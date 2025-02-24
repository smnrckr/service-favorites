package repositories

import (
	"user-favorites-service/internals/models"
)

type FavoritesRepository struct {
	storage PostgreDB
}

func NewFavoritesRepository(s PostgreDB) *FavoritesRepository {
	return &FavoritesRepository{
		storage: s,
	}
}

func (r *FavoritesRepository) AddProductToFavoriteList(favorite *models.Favorite) error {
	return r.storage.GetConnection().Create(favorite).Error
}

func (r *FavoritesRepository) DeleteFavoriteById(userId int, list_id int, productId int) error {

	result := r.storage.GetConnection().Where("user_id = ? AND list_id = ? AND  product_id = ?", userId, list_id, productId).Delete(&models.Favorite{})
	if result.RowsAffected == 0 {
		return models.ErrorNoRowsAffected
	}
	return result.Error
}

func (r *FavoritesRepository) DeleteAllFavoritesByListId(userId int, list_id int) error {
	result := r.storage.GetConnection().Where("user_id = ? AND list_id = ? ", userId, list_id).Delete(&models.Favorite{})
	if result.RowsAffected == 0 {
		return models.ErrorNoRowsAffected
	}
	return result.Error
}

func (r *FavoritesRepository) GetAllFavoritesFromList(listId int) ([]models.Favorite, error) {
	var favorites []models.Favorite
	err := r.storage.GetConnection().Where("list_id = ?", listId).Find(&favorites).Error
	return favorites, err

}
