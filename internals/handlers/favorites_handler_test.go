package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"user-favorites-service/internals/handlers"
	"user-favorites-service/internals/models"
	"user-favorites-service/internals/repositories"
	"user-favorites-service/internals/services"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestFavorites(t *testing.T) {
	mockUserClient := new(handlers.MockUserClient)

	mockFavoriteListRepo := new(handlers.MockFavoriteListsRepo)
	favoriteRespo := repositories.NewFavoritesRepository(testDb)
	favoriteListService := services.NewFavoritesService(favoriteRespo, mockFavoriteListRepo, mockUserClient)
	handler := handlers.NewFavoritesHandler(favoriteListService)

	app := fiber.New()

	handler.FavoritesSetRoutes(app)

	t.Run("Add Product To Favorite List", func(t *testing.T) {
		request := models.FavoriteCreateRequest{
			UserID:    1,
			ProductID: 1,
			ListID:    1,
		}
		requestJSON, err := json.Marshal(request)
		if err != nil {
			t.Fatalf("Error marshaling request: %v", err)
		}
		req := httptest.NewRequest("POST", "/favorites", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		favorite := models.Favorite{}
		err = json.Unmarshal(jsonDataFromHttp, &favorite)
		assert.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, 1, favorite.Id)
	})

	t.Run("Add Product To Favorite List Empty Body", func(t *testing.T) {
		request := map[string]interface{}{}
		requestJSON, err := json.Marshal(request)
		if err != nil {
			t.Fatalf("Error marshaling request: %v", err)
		}
		req := httptest.NewRequest("POST", "/favorites", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		errorResponse := models.ErrorResponse{}
		err = json.Unmarshal(jsonDataFromHttp, &errorResponse)
		assert.NoError(t, err)

		assert.Equal(t, 400, resp.StatusCode)
		assert.NotEmpty(t, errorResponse)
	})

	t.Run("Add Product To Favorite List productId Invalid", func(t *testing.T) {
		request := map[string]interface{}{
			"user_id": 5,
			"list_id": 1,
		}
		requestJSON, err := json.Marshal(request)
		if err != nil {
			t.Fatalf("Error marshaling request: %v", err)
		}
		req := httptest.NewRequest("POST", "/favorites", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		errorResponse := models.ErrorResponse{}
		err = json.Unmarshal(jsonDataFromHttp, &errorResponse)
		assert.NoError(t, err)

		assert.Equal(t, 400, resp.StatusCode)
		assert.NotEmpty(t, errorResponse)
	})

	t.Run("Add Product To Favorite List listId Invalid", func(t *testing.T) {
		request := models.FavoriteCreateRequest{
			UserID:    1,
			ProductID: 1,
			ListID:    5,
		}
		requestJSON, err := json.Marshal(request)
		if err != nil {
			t.Fatalf("Error marshaling request: %v", err)
		}
		req := httptest.NewRequest("POST", "/favorites", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		errorResponse := models.ErrorResponse{}
		err = json.Unmarshal(jsonDataFromHttp, &errorResponse)
		assert.NoError(t, err)

		assert.Equal(t, 500, resp.StatusCode)
		assert.NotEmpty(t, errorResponse)
	})

	t.Run("Get All Favorites From List", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/favorites?userId=1&listId=1", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		favorites := models.FavoritesResponse{}
		err = json.Unmarshal(jsonDataFromHttp, &favorites)
		assert.NoError(t, err)

		assert.NotEmpty(t, favorites)

	})

	t.Run("Delete Product By Id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/favorites?userId=1&listId=1&productId=1", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		reqs := httptest.NewRequest("GET", "/favorites?userId=1&listId=1", nil)
		reqs.Header.Set("Content-Type", "application/json")

		respp, err := app.Test(reqs)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(respp.Body)
		assert.NoError(t, err)

		favorite := models.Favorite{}
		err = json.Unmarshal(jsonDataFromHttp, &favorite)
		assert.NoError(t, err)

		assert.Empty(t, favorite)

	})

}
