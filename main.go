package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/cine-corn-be/middleware"
	"github.com/guneyeroglu/cine-corn-be/routes"
)

func main() {
	app := fiber.New()

	middleware.Init()

	routes.Routes(app)

	app.Listen(":8080")
}
