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
	GetAllFavoritesFromList(listId int) ([]models.Favorite, error)
}

type favoritedUserClient interface {
	CheckUserExist(userID int) (bool, error)
}

type FavoritesService struct {
	favoritesRepository    favoritesRepository
	favoriteListRepository favoriteListsRepo
	userClient             favoritedUserClient
}

func NewFavoritesService(repository favoritesRepository, favoriteListRepo favoriteListsRepo, userClient favoritedUserClient) *FavoritesService {
	return &FavoritesService{
		favoritesRepository:    repository,
		userClient:             userClient,
		favoriteListRepository: favoriteListRepo,
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

func (s *FavoritesService) GetAllFavoritesFromList(listId int, userId int) ([]models.Favorite, error) {
	isUserExist, err := s.userClient.CheckUserExist(userId)
	if err != nil || !isUserExist {
		return nil, err
	}
	return s.favoritesRepository.GetAllFavoritesFromList(listId)
}
