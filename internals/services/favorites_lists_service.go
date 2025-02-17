package services

import (
	"user-favorites-service/internals/models"
)

type favoritesListsRepository interface {
	CreateFavoriteList(list *models.FavoriteList) error
	DeleteFavoriteListById(listId int, userId int) error
	GetFavoriteListsByUserId(userId int) ([]models.FavoriteList, error)
	UpdateFavoriteList(id int, updatedData models.FavoriteList) (models.FavoriteList, error)
}

type favoriteListsUserClient interface {
	CheckUserExist(userID int) (bool, error)
}

type favoritesDeleteRepository interface {
	DeleteAllFavoritesByListId(userId int, list_id int) error
}

type FavoritesListsService struct {
	favoritesListsRepository favoritesListsRepository
	favoritesRespository     favoritesDeleteRepository
	userClient               favoriteListsUserClient
}

func NewFavoritesListsService(favoritesListRepository favoritesListsRepository, favoritesRepository favoritesDeleteRepository, userClient favoriteListsUserClient) *FavoritesListsService {
	return &FavoritesListsService{
		favoritesListsRepository: favoritesListRepository,
		favoritesRespository:     favoritesRepository,
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
	err = s.favoritesRespository.DeleteAllFavoritesByListId(userId, listId)
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

func (s *FavoritesListsService) UpdateFavoriteList(listId int, updatedData models.FavoriteList) (models.FavoriteList, error) {
	updatedList, err := s.favoritesListsRepository.UpdateFavoriteList(listId, updatedData)
	if err != nil {
		return models.FavoriteList{}, err
	}

	return updatedList, nil
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
