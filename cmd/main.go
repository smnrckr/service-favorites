package main

import (
	"os"
	client "user-favorites-service/internals/clients"
	"user-favorites-service/internals/handlers"
	"user-favorites-service/internals/repositories"
	"user-favorites-service/internals/services"
	"user-favorites-service/pkg/postgresql"
	"user-favorites-service/utils"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-fiber/swagger"
	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.LoadEnviromentVariables()
}

func main() {

	host := os.Getenv("HOST")
	dbuser := os.Getenv("USER_NAME")
	dbname := os.Getenv("DB_NAME")
	dbpassword := os.Getenv("PASSWORD")
	port := os.Getenv("PORT")

	db := postgresql.NewDB(postgresql.DbConfig{Host: host, Dbuser: dbuser, Dbname: dbname, Dbpassword: dbpassword, Port: port})
	userClient := client.NewUserClient()

	favoritesListRepo := repositories.NewFavoritesListsRepository(db)
	favoritesRepo := repositories.NewFavoritesRepository(db)

	favoritesService := services.NewFavoritesService(favoritesRepo, favoritesListRepo, userClient)
	favoritesListService := services.NewFavoritesListsService(favoritesListRepo, favoritesRepo, userClient)

	favoritesListHandler := handlers.NewFavoritesListsHandler(favoritesListService)
	favoritesHandler := handlers.NewFavoritesHandler(favoritesService)

	app := fiber.New()

	sw := swagno.New(swagno.Config{Title: "Service Favorites", Version: "v1.0.0", Host: "localhost:8081"})
	sw.AddEndpoints(handlers.FavoritesEndpoints)
	sw.AddEndpoints(handlers.FavoritesListsEndpoints)
	swagger.SwaggerHandler(app, sw.MustToJson(), swagger.WithPrefix("/swagger"))

	favoritesHandler.FavoritesSetRoutes(app)
	favoritesListHandler.FavoritesListsSetRoutes(app)

	app.Listen(":8081")
}
