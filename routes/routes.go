package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/guneyeroglu/cine-corn-be/middleware"
)

func Routes(app *fiber.App) {
	// public api
	api := app.Group("/api/v1")
	api.Use(cors.New())

	// protected api
	api.Use(middleware.JwtMiddleware())
}
