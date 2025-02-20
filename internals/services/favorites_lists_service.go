package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
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

type favoriteListProductClient interface {
	GetProductById(productID int) (models.ProductResponse, error)
}

type favoritesRepoForLists interface {
	DeleteAllFavoritesByListId(userId int, list_id int) error
	GetAllFavoritesFromList(listId int) ([]models.Favorite, error)
}

type FavoritesListsService struct {
	favoritesListsRepository favoritesListsRepository
	favoritesRepository      favoritesRepoForLists
	userClient               favoriteListsUserClient
	productClient            favoriteListProductClient
}

func NewFavoritesListsService(favoritesListRepository favoritesListsRepository, favoritesRepository favoritesRepoForLists, userClient favoriteListsUserClient, productClient favoriteListProductClient) *FavoritesListsService {
	return &FavoritesListsService{
		favoritesListsRepository: favoritesListRepository,
		favoritesRepository:      favoritesRepository,
		userClient:               userClient,
		productClient:            productClient,
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
	err = s.favoritesRepository.DeleteAllFavoritesByListId(userId, listId)
	if err != nil {
		return err
	}
	return s.favoritesListsRepository.DeleteFavoriteListById(listId, userId)
}

func (s *FavoritesListsService) GetFavoriteListsByUserId(userId int) ([]models.UserFavoritesList, error) {
	err := s.checkUserExist(userId)
	if err != nil {
		return nil, err
	}

	favoriteLists, err := s.favoritesListsRepository.GetFavoriteListsByUserId(userId)
	if err != nil {
		return nil, err
	}

	var response []models.UserFavoritesList

	for _, list := range favoriteLists {
		products, err := s.GetAllFavoritesFromList(list.Id, userId)
		if err != nil {
			return nil, err
		}

		response = append(response, models.UserFavoritesList{
			ListID:  list.Id,
			Name:    list.Name,
			Product: products,
		})
	}

	return response, nil
}

func (s *FavoritesListsService) UpdateFavoriteList(listId int, updatedData models.FavoriteList) (models.FavoriteList, error) {
	updatedList, err := s.favoritesListsRepository.UpdateFavoriteList(listId, updatedData)
	if err != nil {
		return models.FavoriteList{}, err
	}

	return updatedList, nil
}

func (s *FavoritesListsService) GetAllFavoritesFromList(listId int, userId int) ([]models.ProductResponse, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxGoroutines := 5
	productInfo := []models.ProductResponse{}

	isUserExist, err := s.userClient.CheckUserExist(userId)
	if err != nil || !isUserExist {
		return nil, err
	}

	favoritesFromList, err := s.favoritesRepository.GetAllFavoritesFromList(listId)
	if err != nil || !isUserExist {
		return nil, err
	}

	ch := make(chan models.ProductResponse, len(favoritesFromList))
	sem := make(chan struct{}, maxGoroutines)
	errs := make(chan error, 1)

	for i := 0; i < len(favoritesFromList); i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(productId int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}

			products, err := s.GetProductInfo(productId)
			defer func() { <-sem }()
			if err != nil {
				fmt.Printf("bu ürün fetch edilemedi %d: %v\n", productId, err)
				cancel()
				errs <- err
				return
			}

			ch <- products
		}(favoritesFromList[i].ProductID)
	}
	go func() {
		wg.Wait()
		close(ch)
		close(errs)
	}()

	for product := range ch {
		mu.Lock()
		productInfo = append(productInfo, product)
		mu.Unlock()

	}

	if len(productInfo) != len(favoritesFromList) {
		err := fmt.Sprintf("expected %d but have %d", len(favoritesFromList), len(productInfo))

		return nil, errors.New(err)
	}

	return productInfo, err
}

func (s *FavoritesListsService) GetProductInfo(productId int) (models.ProductResponse, error) {

	productData, err := s.productClient.GetProductById(productId)
	if err != nil {
		return models.ProductResponse{}, err
	}
	data := models.ProductResponse{Id: productData.Id, ProductName: productData.ProductName, ProductCode: productData.ProductCode, ProductPrice: productData.ProductPrice}
	return data, nil
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
