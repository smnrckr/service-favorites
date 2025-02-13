package handlers

import (
	"fmt"
	"user-favorites-service/internals/models"

	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/gofiber/fiber/v2"
)

type favoritesService interface {
	AddProductToFavoriteList(favorite *models.Favorite) error
	DeleteFavoritetById(userId int, listId int, favoriteProductId int) error
	GetAllFavoritesFromList(listId int, userId int) ([]models.Favorite, error)
}

type FavoritesHandler struct {
	favoritesService favoritesService
}

func NewFavoritesHandler(service favoritesService) *FavoritesHandler {
	return &FavoritesHandler{
		favoritesService: service,
	}
}

func (h *FavoritesHandler) handleAddProductToFavoriteList(c *fiber.Ctx) error {
	favoriteReq := models.FavoriteCreateRequest{}
	if err := c.BodyParser(&favoriteReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}
	newProduct := models.Favorite{UserID: favoriteReq.UserID, ListID: favoriteReq.ListID, ProductID: favoriteReq.ProductID}
	err := h.favoritesService.AddProductToFavoriteList(&newProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": favoriteReq})
}

func (h *FavoritesHandler) handleDeleteFavoriteById(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("userId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	listId, err := c.ParamsInt("listId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	favoriteProductId, err := c.ParamsInt("productId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	err = h.favoritesService.DeleteFavoritetById(userId, listId, favoriteProductId)
	if err != nil {
		if err == models.ErrorNoRowsAffected {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
		}
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{Message: "Product removed from list"})
}

func (h *FavoritesHandler) handleGetAllFavoritesFromList(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("userId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}
	listId, err := c.ParamsInt("listId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	favoriteLists, err := h.favoritesService.GetAllFavoritesFromList(listId, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}
	favoriteResponse := models.FavoritesResponse{Favorite: favoriteLists}

	return c.Status(fiber.StatusOK).JSON(favoriteResponse)
}

func (h *FavoritesHandler) FavoritesSetRoutess(app *fiber.App) {
	userGroup := app.Group("/user/:userId<int>")

	userGroup.Delete("/favorite-list/:listId<int>/favorites/:productId<int>", h.handleDeleteFavoriteById)
	userGroup.Get("/favorite-list/:listId<int>/favorites", h.handleGetAllFavoritesFromList)
	app.Post("/favorite-list/favorites", h.handleAddProductToFavoriteList)

}

var FavoritesEndpoints = []*endpoint.EndPoint{
	endpoint.New(
		endpoint.GET,
		"/user/{userId}/favorite-list/{listId}/favorites",
		endpoint.WithTags("get favorites from favorite-lists"),
		endpoint.WithParams(parameter.IntParam("userId", parameter.Path, parameter.WithRequired(), parameter.WithDescription("Kullanıcı id")), parameter.IntParam("listId", parameter.Path, parameter.WithRequired(), parameter.WithDescription("Liste id"))),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(
			models.FavoritesResponse{}, "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(models.ErrorResponse{}, "Bad Request", "500")}),
		endpoint.WithDescription("kullanıcının listesindeki favorileri döner"),
	),
	endpoint.New(
		endpoint.POST,
		"/favorite-list/favorites",
		endpoint.WithTags("add product to favorite-lists"),
		endpoint.WithBody(models.FavoriteCreateRequest{}),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(
			models.Favorite{}, "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(models.ErrorResponse{}, "Bad Request", "500")}),
		endpoint.WithDescription("favori listesine ürün ekler"),
	),
	endpoint.New(
		endpoint.DELETE,
		"/user/{userId}/favorite-list/{listId}/favorites/{productId}",
		endpoint.WithTags("delete product from favorite-lists"),
		endpoint.WithParams(parameter.IntParam("userId", parameter.Path, parameter.WithRequired(), parameter.WithDescription("Kullanıcı id")), parameter.IntParam("listId", parameter.Path, parameter.WithRequired(), parameter.WithDescription("Product id")), parameter.IntParam("productId", parameter.Path, parameter.WithRequired(), parameter.WithDescription("Product id"))),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccessResponse{}, "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(struct{ Error string }{}, "Bad Request", "500")}),
		endpoint.WithDescription("favori listesindeki ürün silinir"),
	),
}
