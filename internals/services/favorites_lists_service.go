package services

import (
	"user-favorites-service/internals/models"
)

type favoritesListsRepository interface {
	CreateFavoriteList(list *models.FavoriteList) error
	DeleteFavoriteListById(listId int, userId int) error
	GetFavoriteListsByUserId(userId int) ([]models.FavoriteList, error)
	UpdateFavoriteList(id int, updatedData models.FavoriteList) error
}

type favoriteListsUserClient interface {
	CheckUserExist(userID int) (bool, error)
}

type FavoritesListsService struct {
	favoritesListsRepository favoritesListsRepository
	userClient               favoriteListsUserClient
}

func NewFavoritesListsService(repository favoritesListsRepository, userClient favoriteListsUserClient) *FavoritesListsService {
	return &FavoritesListsService{
		favoritesListsRepository: repository,
		userClient:               userClient,
	}
}

func (s *FavoritesListsService) CreateFavoriteList(list *models.FavoriteList) error {
	err := s.checkUserExist(list.UserID)
	if err != nil {
		return err
	}

	return s.favoritesListsRepository.CreateFavoriteList(list)
}

func (s *FavoritesListsService) DeleteFavoriteListById(listId int, userId int) error {
	err := s.checkUserExist(userId)
	if err != nil {
		return err
	}
	return s.favoritesListsRepository.DeleteFavoriteListById(listId, userId)
}

func (s *FavoritesListsService) GetFavoriteListsByUserId(userId int) ([]models.FavoriteList, error) {
	err := s.checkUserExist(userId)
	if err != nil {
		return nil, err
	}
	return s.favoritesListsRepository.GetFavoriteListsByUserId(userId)
}

func (s *FavoritesListsService) UpdateFavoriteList(listId int, updatedData models.FavoriteList) error {
	return s.favoritesListsRepository.UpdateFavoriteList(listId, updatedData)
}

func (s *FavoritesListsService) checkUserExist(userId int) error {
	isUserExist, err := s.userClient.CheckUserExist(userId)
	if err != nil {
		return err
	}
	if !isUserExist {
		return models.ErrorUserNotFound
	}
	return nil
}
