package handlers

import (
	"fmt"
	"user-favorites-service/internals/models"

	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/gofiber/fiber/v2"
)

type favoritesListsService interface {
	CreateFavoriteList(list *models.FavoriteList) error
	DeleteFavoriteListById(listId int, userId int) error
	GetFavoriteListsByUserId(userId int) ([]models.FavoriteList, error)
	UpdateFavoriteList(id int, updatedData models.FavoriteListUpdateRequest) error
}

type FavoritesListsHandler struct {
	favoritesListsService favoritesListsService
}

func NewFavoritesListsHandler(service favoritesListsService) *FavoritesListsHandler {
	return &FavoritesListsHandler{
		favoritesListsService: service,
	}
}
func (h *FavoritesListsHandler) handleCreateFavoriteList(c *fiber.Ctx) error {

	favoriteReq := new(models.FavoriteListCreateRequest)

	if err := c.BodyParser(favoriteReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}
	favoriteList := models.FavoriteList{Name: favoriteReq.Name, UserID: favoriteReq.UserID}
	err := h.favoritesListsService.CreateFavoriteList(&favoriteList)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(favoriteList)
}

func (h *FavoritesListsHandler) handleDeleteFavoriteListById(c *fiber.Ctx) error {
	listId, err := c.ParamsInt("listId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}
	userId := c.QueryInt("userId", -1)

	if userId == -1 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "userId required"})
	}

	err = h.favoritesListsService.DeleteFavoriteListById(listId, userId)
	if err != nil {
		if err == models.ErrorNoRowsAffected {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
		}
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{Message: "list deleted successfuly"})
}

func (h *FavoritesListsHandler) handleGetAllFavoriteLists(c *fiber.Ctx) error {
	userId := c.QueryInt("userId", -1)
	if userId == -1 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "userId required"})
	}

	favoriteLists, err := h.favoritesListsService.GetFavoriteListsByUserId(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}
	favoriteListResponse := models.FavoriteListResponse{FavoriteList: favoriteLists}

	return c.Status(fiber.StatusOK).JSON(favoriteListResponse)
}

func (h *FavoritesListsHandler) handleUpdateFavoriteList(c *fiber.Ctx) error {

	listId, err := c.ParamsInt("listId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz list ID"})
	}

	userRequest := models.FavoriteListUpdateRequest{}
	err = c.BodyParser(&userRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz JSON verisi"})
	}

	if userRequest.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": " name boş olamaz"})
	}

	err = h.favoritesListsService.UpdateFavoriteList(listId, userRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Liste ismi başarıyla güncellendi"})
}

func (h *FavoritesListsHandler) FavoritesListsSetRoutes(app *fiber.App) {
	favoriteListGroup := app.Group("/favorite-lists")

	favoriteListGroup.Delete("/:listId<int>", h.handleDeleteFavoriteListById)
	favoriteListGroup.Get("/", h.handleGetAllFavoriteLists)
	favoriteListGroup.Put("/:listId<int>", h.handleUpdateFavoriteList)
	favoriteListGroup.Post("/", h.handleCreateFavoriteList)

}

var FavoritesListsEndpoints = []*endpoint.EndPoint{
	endpoint.New(
		endpoint.GET,
		"/favorite-lists",
		endpoint.WithTags("favorite-lists"),
		endpoint.WithParams(parameter.IntParam("userId", parameter.Query, parameter.WithRequired(), parameter.WithDescription("Kullanıcı id"))),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(
			models.FavoriteListResponse{}, "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(models.ErrorResponse{}, "Bad Request", "500")}),
		endpoint.WithDescription("kullanıcıların favori listelerini döner"),
	),
	endpoint.New(
		endpoint.POST,
		"/favorite-lists",
		endpoint.WithTags("favorite-lists"),
		endpoint.WithBody(models.FavoriteListCreateRequest{}),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(
			models.FavoriteList{}, "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(models.ErrorResponse{}, "Bad Request", "500")}),
		endpoint.WithDescription("yeni favori listesi oluşturur"),
	),
	endpoint.New(
		endpoint.DELETE,
		"/favorite-lists/{listId}",
		endpoint.WithTags("favorite-lists"),
		endpoint.WithParams(parameter.IntParam("userId", parameter.Query, parameter.WithRequired(), parameter.WithDescription("Kullanıcı id")), parameter.IntParam("listId", parameter.Path, parameter.WithRequired(), parameter.WithDescription("Liste id"))),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccessResponse{}, "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(models.ErrorResponse{}, "Bad Request", "500")}),
		endpoint.WithDescription("favori listesi silinir"),
	),
	endpoint.New(
		endpoint.PUT,
		"/favorite-lists/{listId}",
		endpoint.WithTags("favorite-lists"),
		endpoint.WithParams(parameter.IntParam("listId", parameter.Path, parameter.WithRequired())),
		endpoint.WithBody(models.FavoriteListUpdateRequest{}),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(
			models.FavoriteListResponse{}, "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(models.ErrorResponse{}, "Bad Request", "500")}),
		endpoint.WithDescription("listenin ismini değiştirir"),
	),
}
