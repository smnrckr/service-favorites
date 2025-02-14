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

func TestFavoritesLists(t *testing.T) {
	mockUserClient := new(MockUserClient)
	favoriteListRespo := repositories.NewFavoritesListsRepository(testDb)
	favoriteListService := services.NewFavoritesListsService(favoriteListRespo, mockUserClient)
	handler := handlers.NewFavoritesListsHandler(favoriteListService)

	app := fiber.New()

	handler.FavoritesListsSetRoutes(app)

	t.Run("CreateFavoriteList", func(t *testing.T) {
		request := models.FavoriteListCreateRequest{
			Name:   "ayakkabılar",
			UserID: 1,
		}
		requestJSON, err := json.Marshal(request)
		if err != nil {
			t.Fatalf("Error marshaling request: %v", err)
		}

		req := httptest.NewRequest("POST", "/favorite-lists", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		favoriteList := models.FavoriteList{}
		err = json.Unmarshal(jsonDataFromHttp, &favoriteList)
		assert.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, 1, favoriteList.Id)

	})

	t.Run("CreateFavoriteListInvalid", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"name": "ayakkabılar",
		}

		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("Error marshaling request: %v", err)
		}

		req := httptest.NewRequest("POST", "/favorite-lists", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		errorResp := models.ErrorResponse{}
		err = json.Unmarshal(jsonDataFromHttp, &errorResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, errorResp)
		assert.Equal(t, 400, resp.StatusCode)

	})
	t.Run("CreateFavoriteListUserInvalid", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"name":    "ayakkabılar",
			"user_id": 999,
		}

		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("Error marshaling request: %v", err)
		}

		req := httptest.NewRequest("POST", "/favorite-lists", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		errorResp := models.ErrorResponse{}
		err = json.Unmarshal(jsonDataFromHttp, &errorResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, errorResp)
		assert.Equal(t, 500, resp.StatusCode)

	})
}
