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

	if err := favoriteReq.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	newProduct := models.Favorite{UserID: favoriteReq.UserID, ListID: favoriteReq.ListID, ProductID: favoriteReq.ProductID}
	err := h.favoritesService.AddProductToFavoriteList(&newProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(newProduct)
}

func (h *FavoritesHandler) handleDeleteFavoriteById(c *fiber.Ctx) error {
	userId := c.QueryInt("userId", -1)
	if userId == -1 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "userId required"})
	}

	listId := c.QueryInt("listId", -1)
	if listId == -1 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "listId required"})
	}

	productId := c.QueryInt("productId", -1)
	if productId == -1 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "productId required"})
	}

	err := h.favoritesService.DeleteFavoritetById(userId, listId, productId)
	if err != nil {
		if err == models.ErrorNoRowsAffected {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
		}
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{Message: "Product removed from list"})
}

func (h *FavoritesHandler) FavoritesSetRoutes(app *fiber.App) {
	favoriteGroup := app.Group("/favorites")

	favoriteGroup.Delete("/", h.handleDeleteFavoriteById)
	favoriteGroup.Post("/", h.handleAddProductToFavoriteList)

}

var FavoritesEndpoints = []*endpoint.EndPoint{
	endpoint.New(
		endpoint.POST,
		"/favorites",
		endpoint.WithTags("favorites"),
		endpoint.WithBody(models.FavoriteCreateRequest{}),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(
			models.Favorite{}, "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(models.ErrorResponse{}, "Bad Request", "500")}),
		endpoint.WithDescription("favori listesine ürün ekler"),
	),
	endpoint.New(
		endpoint.DELETE,
		"/favorites",
		endpoint.WithTags("favorites"),
		endpoint.WithParams(parameter.IntParam("userId", parameter.Query, parameter.WithRequired(), parameter.WithDescription("Kullanıcı id")), parameter.IntParam("listId", parameter.Query, parameter.WithRequired(), parameter.WithDescription("Product id")), parameter.IntParam("productId", parameter.Query, parameter.WithRequired(), parameter.WithDescription("Product id"))),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccessResponse{}, "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(struct{ Error string }{}, "Bad Request", "500")}),
		endpoint.WithDescription("favori listesindeki ürün silinir"),
	),
}
