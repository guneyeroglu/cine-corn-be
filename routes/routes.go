package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/guneyeroglu/cine-corn-be/handlers"
	"github.com/guneyeroglu/cine-corn-be/middleware"
)

func Routes(app *fiber.App) {
	// public api
	api := app.Group("/api/v1")
	api.Use(cors.New())
	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)
	api.Get("/movies", handlers.GetMoviesList)
	api.Post("/movie-details", handlers.GetMovieDetailsList)

	// protected api
	protectedApi := api.Group("")
	protectedApi.Use(middleware.JwtMiddleware())
	protectedApi.Get("/auth-user", handlers.GetAuthUser)
	protectedApi.Get("/favorites", handlers.GetFavoriteMovies)
	protectedApi.Put("/favorites", handlers.ToggleFavoriteMovie)
	protectedApi.Get("/list", handlers.GetListMovies)
	protectedApi.Put("/list", handlers.ToggleListMovie)

}
