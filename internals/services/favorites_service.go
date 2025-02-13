package services

import (
	"user-favorites-service/internals/models"
)

type favoriteListsRepoFavoriteService interface {
	ChechFavoriteListExist(listId int, userId int) (bool, error)
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
	userFavoritesRepository favoritesRepository
	favoriteListRepository  favoriteListsRepoFavoriteService
	userClient              favoritedUserClient
}

func NewFavoritesService(repository favoritesRepository, favoriteListRepo favoriteListsRepoFavoriteService, userClient favoritedUserClient) *FavoritesService {
	return &FavoritesService{
		userFavoritesRepository: repository,
		userClient:              userClient,
		favoriteListRepository:  favoriteListRepo,
	}
}

func (s *FavoritesService) AddProductToFavoriteList(favorite *models.Favorite) error {
	isFavoriteListExist, err := s.favoriteListRepository.ChechFavoriteListExist(favorite.ListID, favorite.UserID)
	if err != nil {
		return err
	}

	if !isFavoriteListExist {
		return models.ErrorListNotFound
	}
	return s.userFavoritesRepository.AddProductToFavoriteList(favorite)
}

func (s *FavoritesService) DeleteFavoritetById(userId int, listId int, favoriteProductId int) error {
	return s.userFavoritesRepository.DeleteFavoriteById(userId, listId, favoriteProductId)
}

func (s *FavoritesService) GetAllFavoritesFromList(listId int, userId int) ([]models.Favorite, error) {
	isUserExist, err := s.userClient.CheckUserExist(userId)
	if err != nil || !isUserExist {
		return nil, err
	}
	return s.userFavoritesRepository.GetAllFavoritesFromList(listId)
}
