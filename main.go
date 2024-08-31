package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/cine-corn-be/config"
	"github.com/guneyeroglu/cine-corn-be/routes"
)

func main() {
	config.Init()
	app := fiber.New()
	routes.Routes(app)
	app.Listen(":8080")
}
