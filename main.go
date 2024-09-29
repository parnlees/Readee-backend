package main

import (
	"Readee-Backend/common/config"
	"Readee-Backend/common/database"
	"Readee-Backend/endpoint"

	"github.com/gofiber/fiber/v2"
	//"Readee-Backend/fiber"
)

func main() {
	config.Init()
	database.Init()
	//fiber.Init()

	// Initialize Fiber app
	app := fiber.New()

	// Register all routes
	endpoint.RegisterRoutes(app)

	// Start the server
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
