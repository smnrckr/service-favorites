package handlers

import (
	"user-favorites-service/internals/models"
)

type MockProductClient struct {
}

func (productClient *MockProductClient) GetProductById(productID int) (models.ProductResponse, error) {
	productsData := models.ProductResponse{}

	return productsData, nil
}

type MockUserClient struct {
}

func (m *MockUserClient) CheckUserExist(userID int) (bool, error) {
	switch userID {
	case 1:
		return true, nil
	case 2:
		return false, nil
	}
	return false, nil
}

type MockFavoriteListsRepo struct{}

func (m *MockFavoriteListsRepo) CheckFavoriteListExist(listId int, userId int) (bool, error) {
	switch {
	case listId == 1 && userId == 1:
		return true, nil
	case listId == 2 && userId == 1:
		return false, nil
	default:
		return false, nil
	}
}
