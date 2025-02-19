package services

import (
	"user-favorites-service/internals/models"
)

type favoriteListsRepo interface {
	CheckFavoriteListExist(listId int, userId int) (bool, error)
}

type favoritesRepository interface {
	AddProductToFavoriteList(favorite *models.Favorite) error
	DeleteFavoriteById(userId int, list_id int, productId int) error
}

type favoritedUserClient interface {
	CheckUserExist(userID int) (bool, error)
}

type favoritesProductClient interface {
	GetProductById(productID int) (models.ProductResponse, error)
}

type FavoritesService struct {
	favoritesRepository    favoritesRepository
	favoriteListRepository favoriteListsRepo
	userClient             favoritedUserClient
	productClient          favoritesProductClient
}

func NewFavoritesService(repository favoritesRepository, favoriteListRepo favoriteListsRepo, userClient favoritedUserClient, productClient favoritesProductClient) *FavoritesService {
	return &FavoritesService{
		favoritesRepository:    repository,
		userClient:             userClient,
		favoriteListRepository: favoriteListRepo,
		productClient:          productClient,
	}
}

func (s *FavoritesService) AddProductToFavoriteList(favorite *models.Favorite) error {
	isFavoriteListExist, err := s.favoriteListRepository.CheckFavoriteListExist(favorite.ListID, favorite.UserID)
	if err != nil {
		return err
	}

	if !isFavoriteListExist {
		return models.ErrorListNotFound
	}
	return s.favoritesRepository.AddProductToFavoriteList(favorite)
}

func (s *FavoritesService) DeleteFavoritetById(userId int, listId int, favoriteProductId int) error {
	return s.favoritesRepository.DeleteFavoriteById(userId, listId, favoriteProductId)
}
