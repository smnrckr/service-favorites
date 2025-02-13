package repositories

import (
	"user-favorites-service/internals/models"

	"gorm.io/gorm"
)

type FavoritesListsStorage interface {
	GetConnection() *gorm.DB
	Close()
}

type FavoritesListsRepository struct {
	storage FavoritesListsStorage
}

func NewFavoritesListsRepository(s FavoritesStorage) *FavoritesListsRepository {
	return &FavoritesListsRepository{
		storage: s,
	}
}

func (r *FavoritesListsRepository) CreateFavoriteList(list *models.FavoriteList) error {
	return r.storage.GetConnection().Create(&list).Error
}

func (r *FavoritesListsRepository) DeleteFavoriteListById(id int, userId int) error {
	result := r.storage.GetConnection().Where("id = ? AND user_id = ?", id, userId).Delete(&models.FavoriteList{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrorNoRowsAffected
	}
	return nil
}

func (r *FavoritesListsRepository) GetFavoriteListsByUserId(userId int) ([]models.FavoriteList, error) {
	favoriteLists := []models.FavoriteList{}
	result := r.storage.GetConnection().Where("user_id = ?", userId).Find(&favoriteLists)
	return favoriteLists, result.Error
}

func (r *FavoritesListsRepository) ChechFavoriteListExist(listId int, userId int) (bool, error) {
	favoriteList := models.FavoriteList{}
	result := r.storage.GetConnection().Where("user_id = ? AND id = ?  ", userId, listId).Find(&favoriteList)
	if result.Error != nil {
		return false, result.Error
	}

	if favoriteList == (models.FavoriteList{}) {
		return false, nil
	}
	return true, nil
}

func (r *FavoritesListsRepository) UpdateFavoriteList(listId int, updatedData models.FavoriteListUpdateRequest) error {
	newData := models.FavoriteList{
		Name: updatedData.Name,
	}
	if updatedData.Name == "" {
		return nil
	}

	result := r.storage.GetConnection().Where("id = ?", listId).Updates(&newData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrorListNotFound
	}

	return nil

}
