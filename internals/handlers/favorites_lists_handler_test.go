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

func TestFavoritesLists(t *testing.T) {
	mockUserClient := new(handlers.MockUserClient)
	favoriteListRespo := repositories.NewFavoritesListsRepository(testDb)
	favoriteListService := services.NewFavoritesListsService(favoriteListRespo, mockUserClient)
	handler := handlers.NewFavoritesListsHandler(favoriteListService)

	app := fiber.New()

	handler.FavoritesListsSetRoutes(app)

	t.Run("Create Favorite List", func(t *testing.T) {
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
		assert.Equal(t, 2, favoriteList.Id)

	})

	t.Run("Create Favorite List Invalid", func(t *testing.T) {
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
	t.Run("Create Favorite List User Invalid", func(t *testing.T) {
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

	t.Run("Get Favorite Lists By UserId", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/favorite-lists?userId=1", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		favoriteLists := models.FavoriteListResponse{}
		err = json.Unmarshal(jsonDataFromHttp, &favoriteLists)
		assert.NoError(t, err)

		assert.NotEmpty(t, favoriteLists)

	})

	t.Run("Get Favorite Lists By UserId Invalid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/favorite-lists", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		errorResponse := models.ErrorResponse{}
		err = json.Unmarshal(jsonDataFromHttp, &errorResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, errorResponse)

	})

	t.Run("Update Favorite Lists", func(t *testing.T) {
		newData := models.FavoriteListUpdateRequest{
			Name: "çantalar",
		}
		requestJSON, err := json.Marshal(newData)
		if err != nil {
			t.Fatalf("Error marshaling request: %v", err)
		}

		req := httptest.NewRequest("PUT", "/favorite-lists/1", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		jsonDataFromHttp, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		favoriteListsUpdate := models.FavoriteList{}
		err = json.Unmarshal(jsonDataFromHttp, &favoriteListsUpdate)
		assert.NoError(t, err)

		assert.Equal(t, "çantalar", favoriteListsUpdate.Name)

	})

	t.Run("Delete Favorite Lists", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/favorite-lists/1?userId=1", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		reqs := httptest.NewRequest("GET", "/favorite-lists?userId=1", nil)
		reqs.Header.Set("Content-Type", "application/json")

		respp, err := app.Test(reqs)
		assert.NoError(t, err)

		jsonDataFromHttp, err := io.ReadAll(respp.Body)
		assert.NoError(t, err)

		favoriteLists := models.FavoriteList{}
		err = json.Unmarshal(jsonDataFromHttp, &favoriteLists)
		assert.NoError(t, err)

		assert.Empty(t, favoriteLists)

	})
}
